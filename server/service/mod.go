package service

import (
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

	var mods []models.Mod
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

	return mods, nil
}

// UploadMod 上传并安装MOD
func (s *ModService) UploadMod(srcPath, filename string) (*models.Mod, error) {
	// 确保MOD目录存在
	if err := os.MkdirAll(s.modsPath, 0755); err != nil {
		return nil, fmt.Errorf("创建MOD目录失败: %w", err)
	}

	// 生成目标路径
	modDir := filepath.Join(s.modsPath, strings.TrimSuffix(filename, filepath.Ext(filename)))

	// 如果是zip文件，解压
	if strings.HasSuffix(filename, ".zip") {
		if err := s.unzipMod(srcPath, modDir); err != nil {
			return nil, fmt.Errorf("解压MOD失败: %w", err)
		}
		os.Remove(srcPath) // 删除上传的压缩文件
	} else {
		// 直接移动文件
		if err := os.MkdirAll(modDir, 0755); err != nil {
			return nil, err
		}
		destPath := filepath.Join(modDir, filename)
		if err := os.Rename(srcPath, destPath); err != nil {
			return nil, err
		}
	}

	// 读取manifest.json获取MOD信息
	mod, err := s.parseModManifest(modDir)
	if err != nil {
		// 如果无法解析manifest，使用文件名作为MOD名称
		mod = &models.Mod{
			Name:     strings.TrimSuffix(filename, filepath.Ext(filename)),
			FilePath: modDir,
			Enabled:  true,
		}
	} else {
		mod.FilePath = modDir
		mod.Enabled = true
	}

	// 保存到数据库
	result, err := database.DB.Exec(`
		INSERT INTO mods (name, unique_id, version, author, description, enabled, file_path)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, mod.Name, mod.UniqueID, mod.Version, mod.Author, mod.Description, mod.Enabled, mod.FilePath)

	if err != nil {
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

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modDir := filepath.Join(s.modsPath, entry.Name())

		// 检查是否已在数据库中
		var count int
		err := database.DB.QueryRow(`SELECT COUNT(*) FROM mods WHERE file_path = ?`, modDir).Scan(&count)
		if err != nil {
			continue
		}
		if count > 0 {
			continue // 已存在
		}

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

		// 保存到数据库
		database.DB.Exec(`
			INSERT INTO mods (name, unique_id, version, author, description, enabled, file_path)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, mod.Name, mod.UniqueID, mod.Version, mod.Author, mod.Description, mod.Enabled, mod.FilePath)
	}

	return nil
}

// parseModManifest 解析MOD的manifest.json
func (s *ModService) parseModManifest(modDir string) (*models.Mod, error) {
	manifestPath := filepath.Join(modDir, "manifest.json")
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
	// 使用install.go中的unzip逻辑
	// 这里简化实现，实际应该复用
	return fmt.Errorf("MOD解压功能需要实现")
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
