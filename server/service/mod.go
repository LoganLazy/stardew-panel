package service

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"stardew-panel/database"
	"stardew-panel/models"
	"strings"
)

const (
	maxModFiles       = 10000
	maxModExtractSize = uint64(2 * 1024 * 1024 * 1024)
)

// ModService MOD管理服务
type ModService struct {
	modsPath string
}

// NewModService 创建MOD服务
func NewModService(modsPath string) *ModService {
	return &ModService{
		modsPath: modsPath,
	}
}

// ListMods 列出所有MOD
func (s *ModService) ListMods() ([]models.Mod, error) {
	rows, err := database.DB.Query(`
		SELECT id, name, unique_id, version, author, description, enabled, file_path, created_at, updated_at
		FROM mods
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mods := make([]models.Mod, 0)
	for rows.Next() {
		var mod models.Mod
		var uniqueID, version, author, description sql.NullString

		err := rows.Scan(&mod.ID, &mod.Name, &uniqueID, &version, &author, &description,
			&mod.Enabled, &mod.FilePath, &mod.CreatedAt, &mod.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if uniqueID.Valid {
			mod.UniqueID = uniqueID.String
		}
		if version.Valid {
			mod.Version = version.String
		}
		if author.Valid {
			mod.Author = author.String
		}
		if description.Valid {
			mod.Description = description.String
		}

		mods = append(mods, mod)
	}

	return mods, rows.Err()
}

// UploadMod 上传并安装MOD
func (s *ModService) UploadMod(srcPath, filename string) (*models.Mod, error) {
	// 确保MOD目录存在
	if err := os.MkdirAll(s.modsPath, 0755); err != nil {
		return nil, fmt.Errorf("创建MOD目录失败: %w", err)
	}

	// Multipart 文件名可能包含路径分隔符，必须只使用最后一段。
	filename = filepath.Base(strings.ReplaceAll(filename, "\\", "/"))
	if filename == "." || filename == ".." || filename == "" {
		return nil, fmt.Errorf("无效的 MOD 文件名")
	}
	modName := strings.TrimSuffix(filename, filepath.Ext(filename))
	if modName == "" || modName == "." || modName == ".." {
		return nil, fmt.Errorf("无效的 MOD 名称")
	}
	modDir := filepath.Join(s.modsPath, modName)
	if existing, err := os.Stat(modDir); err == nil && existing.IsDir() {
		return nil, fmt.Errorf("MOD 已存在: %s", modName)
	}

	if !strings.EqualFold(filepath.Ext(filename), ".zip") {
		return nil, fmt.Errorf("仅支持 .zip 格式的 MOD")
	}
	if err := s.unzipMod(srcPath, modDir); err != nil {
		return nil, fmt.Errorf("解压MOD失败: %w", err)
	}
	_ = os.Remove(srcPath)

	// 读取manifest.json获取MOD信息
	mod, err := s.parseModManifest(modDir)
	if err != nil {
		_ = os.RemoveAll(modDir)
		return nil, fmt.Errorf("MOD 缺少有效的 manifest.json: %w", err)
	}
	mod.FilePath = modDir
	mod.Enabled = true

	// 保存到数据库
	result, err := database.DB.Exec(`
		INSERT INTO mods (name, unique_id, version, author, description, enabled, file_path)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, mod.Name, nullableString(mod.UniqueID), mod.Version, mod.Author, mod.Description, mod.Enabled, mod.FilePath)

	if err != nil {
		_ = os.RemoveAll(modDir)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	mod.ID = int(id)

	return mod, nil
}

// ToggleMod 启用/禁用MOD
func (s *ModService) ToggleMod(id int, enabled bool) error {
	// 获取MOD信息
	row := database.DB.QueryRow(`SELECT file_path, enabled FROM mods WHERE id = ?`, id)
	var filePath string
	var currentEnabled bool
	if err := row.Scan(&filePath, &currentEnabled); err != nil {
		return err
	}
	if !isWithin(filePath, s.modsPath) {
		return fmt.Errorf("MOD 路径不在允许的目录内")
	}

	// 如果状态已经是目标状态，直接返回
	if currentEnabled == enabled {
		return nil
	}

	// 禁用MOD：重命名文件夹，添加 .disabled 后缀
	if !enabled {
		disabledPath := filePath + ".disabled"
		if err := os.Rename(filePath, disabledPath); err != nil {
			return fmt.Errorf("禁用MOD失败: %w", err)
		}
		// 更新数据库中的路径
		_, err := database.DB.Exec(`UPDATE mods SET enabled = ?, file_path = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
			enabled, disabledPath, id)
		if err != nil {
			_ = os.Rename(disabledPath, filePath)
		}
		return err
	}

	// 启用MOD：移除 .disabled 后缀
	if strings.HasSuffix(filePath, ".disabled") {
		enabledPath := strings.TrimSuffix(filePath, ".disabled")
		if err := os.Rename(filePath, enabledPath); err != nil {
			return fmt.Errorf("启用MOD失败: %w", err)
		}
		// 更新数据库中的路径
		_, err := database.DB.Exec(`UPDATE mods SET enabled = ?, file_path = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
			enabled, enabledPath, id)
		if err != nil {
			_ = os.Rename(enabledPath, filePath)
		}
		return err
	}

	// 如果文件路径没有.disabled后缀，只更新数据库
	_, err := database.DB.Exec(`UPDATE mods SET enabled = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, enabled, id)
	return err
}

// DeleteMod 删除MOD
func (s *ModService) DeleteMod(id int) error {
	// 获取MOD文件路径
	row := database.DB.QueryRow(`SELECT file_path FROM mods WHERE id = ?`, id)
	var filePath string
	if err := row.Scan(&filePath); err != nil {
		return err
	}
	if !isWithin(filePath, s.modsPath) {
		return fmt.Errorf("MOD 路径不在允许的目录内")
	}

	// 删除文件
	if err := os.RemoveAll(filePath); err != nil {
		return fmt.Errorf("删除MOD文件失败: %w", err)
	}

	// 从数据库删除
	_, err := database.DB.Exec(`DELETE FROM mods WHERE id = ?`, id)
	return err
}

// ScanMods 扫描MOD目录并同步到数据库
func (s *ModService) ScanMods() error {
	// 检查MOD目录是否存在
	if _, err := os.Stat(s.modsPath); os.IsNotExist(err) {
		return nil // MOD目录不存在，跳过扫描
	}

	// 读取目录
	entries, err := os.ReadDir(s.modsPath)
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	seen := make(map[string]struct{})
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modDir := filepath.Join(s.modsPath, entry.Name())
		seen[filepath.Clean(modDir)] = struct{}{}

		// 解析manifest
		mod, err := s.parseModManifest(modDir)
		if err != nil {
			// 无法解析，使用目录名
			mod = &models.Mod{
				Name:     entry.Name(),
				FilePath: modDir,
				Enabled:  !strings.HasSuffix(entry.Name(), ".disabled"),
			}
		} else {
			mod.FilePath = modDir
			mod.Enabled = !strings.HasSuffix(entry.Name(), ".disabled")
		}

		result, err := tx.Exec(`
			UPDATE mods
			SET name = ?, unique_id = ?, version = ?, author = ?, description = ?,
				enabled = ?, updated_at = CURRENT_TIMESTAMP
			WHERE file_path = ?
		`, mod.Name, nullableString(mod.UniqueID), mod.Version, mod.Author, mod.Description, mod.Enabled, mod.FilePath)
		if err != nil {
			return err
		}
		updated, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if updated == 0 {
			if _, err := tx.Exec(`
			INSERT INTO mods (name, unique_id, version, author, description, enabled, file_path)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			`, mod.Name, nullableString(mod.UniqueID), mod.Version, mod.Author, mod.Description, mod.Enabled, mod.FilePath); err != nil {
				return err
			}
		}
	}

	rows, err := tx.Query(`SELECT id, file_path FROM mods`)
	if err != nil {
		return err
	}
	stale := make([]int, 0)
	for rows.Next() {
		var id int
		var path string
		if err := rows.Scan(&id, &path); err != nil {
			rows.Close()
			return err
		}
		if !isWithin(path, s.modsPath) {
			continue
		}
		if _, ok := seen[filepath.Clean(path)]; !ok {
			stale = append(stale, id)
		}
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, id := range stale {
		if _, err := tx.Exec(`DELETE FROM mods WHERE id = ?`, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func nullableString(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

// parseModManifest 解析MOD的manifest.json
func (s *ModService) parseModManifest(modDir string) (*models.Mod, error) {
	var manifestPath string
	err := filepath.Walk(modDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !info.IsDir() && strings.EqualFold(info.Name(), "manifest.json") {
			manifestPath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if manifestPath == "" {
		return nil, fmt.Errorf("manifest.json 不存在")
	}
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var manifest struct {
		Name        string `json:"Name"`
		UniqueID    string `json:"UniqueID"`
		Version     string `json:"Version"`
		Author      string `json:"Author"`
		Description string `json:"Description"`
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &models.Mod{
		Name:        manifest.Name,
		UniqueID:    manifest.UniqueID,
		Version:     manifest.Version,
		Author:      manifest.Author,
		Description: manifest.Description,
	}, nil
}

// unzipMod 解压MOD文件
func (s *ModService) unzipMod(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	if len(reader.File) == 0 || len(reader.File) > maxModFiles {
		return fmt.Errorf("MOD 压缩包文件数量无效")
	}
	var totalSize uint64
	var extractedSize uint64
	for _, file := range reader.File {
		name := strings.ReplaceAll(file.Name, "\\", "/")
		if name == "" || strings.HasPrefix(name, "/") {
			return fmt.Errorf("MOD 包含非法路径: %s", file.Name)
		}
		clean := filepath.Clean(filepath.FromSlash(name))
		if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
			return fmt.Errorf("MOD 包含非法路径: %s", file.Name)
		}
		if file.UncompressedSize64 > maxModExtractSize-totalSize {
			return fmt.Errorf("MOD 解压后大小超过限制")
		}
		totalSize += file.UncompressedSize64
		if file.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("MOD 不允许包含符号链接")
		}
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	staging, err := os.MkdirTemp(filepath.Dir(dest), ".mod-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(staging)

	for _, file := range reader.File {
		name := filepath.Clean(filepath.FromSlash(strings.ReplaceAll(file.Name, "\\", "/")))
		path := filepath.Join(staging, name)
		if !isWithin(path, staging) {
			return fmt.Errorf("MOD 包含非法路径: %s", file.Name)
		}
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			return err
		}
		in, err := file.Open()
		if err == nil {
			remaining := maxModExtractSize - extractedSize
			var written int64
			written, err = io.Copy(out, io.LimitReader(in, int64(remaining)+1))
			if err == nil && (written < 0 || uint64(written) > remaining) {
				err = fmt.Errorf("MOD 解压后大小超过限制")
			}
			if err == nil {
				extractedSize += uint64(written)
			}
			_ = in.Close()
		}
		closeErr := out.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
	}

	// 大多数 MOD 压缩包自带一个顶层目录，去掉这一层避免 Mods/A/A/manifest.json。
	entries, err := os.ReadDir(staging)
	if err != nil {
		return err
	}
	installRoot := staging
	if len(entries) == 1 && entries[0].IsDir() {
		installRoot = filepath.Join(staging, entries[0].Name())
	}
	return os.Rename(installRoot, dest)
}

func isWithin(path, root string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) && !filepath.IsAbs(rel)
}

// CopyFile 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
