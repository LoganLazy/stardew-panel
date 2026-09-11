package service

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"stardew-panel/database"
	"stardew-panel/models"
	"strings"
	"sync"
	"time"
)

const (
	maxBackupFiles       = 100000
	maxBackupExtractSize = uint64(20 * 1024 * 1024 * 1024)
)

// SaveService 存档管理服务
type SaveService struct {
	savesPath  string
	backupPath string
	backupMu   sync.Mutex
}

// NewSaveService 创建存档服务
func NewSaveService(savesPath string, backupPaths ...string) *SaveService {
	backupPath := filepath.Join(savesPath, "backups")
	if len(backupPaths) > 0 && backupPaths[0] != "" {
		backupPath = backupPaths[0]
	}
	os.MkdirAll(backupPath, 0755)
	return &SaveService{
		savesPath:  savesPath,
		backupPath: backupPath,
	}
}

// ListBackups 列出所有备份
func (s *SaveService) ListBackups() ([]models.Backup, error) {
	rows, err := database.DB.Query(`
		SELECT id, save_name, backup_path, size, created_at
		FROM backups
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	backups := make([]models.Backup, 0)
	for rows.Next() {
		var backup models.Backup
		err := rows.Scan(&backup.ID, &backup.SaveName, &backup.BackupPath, &backup.Size, &backup.CreatedAt)
		if err != nil {
			return nil, err
		}
		backups = append(backups, backup)
	}

	return backups, rows.Err()
}

// CreateBackup 创建存档备份
func (s *SaveService) CreateBackup(saveName string) (*models.Backup, error) {
	s.backupMu.Lock()
	defer s.backupMu.Unlock()
	return s.createBackup(saveName)
}

func (s *SaveService) createBackup(saveName string) (*models.Backup, error) {
	savePath, err := s.safeSavePath(saveName)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(savePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("存档不存在: %s", saveName)
	}
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("存档不是目录: %s", saveName)
	}
	lstat, err := os.Lstat(savePath)
	if err != nil {
		return nil, err
	}
	if lstat.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("存档目录不能是符号链接")
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405.000000000")
	backupName := fmt.Sprintf("%s_%s.zip", saveName, timestamp)
	backupPath := filepath.Join(s.backupPath, backupName)

	// 压缩存档
	if err := s.zipDirectory(savePath, backupPath); err != nil {
		return nil, fmt.Errorf("压缩存档失败: %w", err)
	}

	// 获取备份文件大小
	fileInfo, err := os.Stat(backupPath)
	if err != nil {
		return nil, err
	}

	// 保存到数据库
	result, err := database.DB.Exec(`
		INSERT INTO backups (save_name, backup_path, size)
		VALUES (?, ?, ?)
	`, saveName, backupPath, fileInfo.Size())

	if err != nil {
		os.Remove(backupPath) // 清理备份文件
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Backup{
		ID:         int(id),
		SaveName:   saveName,
		BackupPath: backupPath,
		Size:       fileInfo.Size(),
		CreatedAt:  time.Now(),
	}, nil
}

// RestoreBackup 恢复存档
func (s *SaveService) RestoreBackup(backupID int) error {
	s.backupMu.Lock()
	defer s.backupMu.Unlock()

	// 获取备份信息
	row := database.DB.QueryRow(`
		SELECT save_name, backup_path
		FROM backups
		WHERE id = ?
	`, backupID)

	var saveName, backupPath string
	if err := row.Scan(&saveName, &backupPath); err != nil {
		return err
	}
	if !isWithin(backupPath, s.backupPath) {
		return fmt.Errorf("备份路径不在允许的目录内")
	}
	savePath, err := s.safeSavePath(saveName)
	if err != nil {
		return err
	}

	// 检查备份文件是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在")
	} else if err != nil {
		return err
	}

	// 先完整解压到同一文件系统中的暂存目录，确认备份可用后再替换。
	staging, err := os.MkdirTemp(s.savesPath, ".restore-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(staging)
	if err := s.unzipDirectory(backupPath, staging); err != nil {
		return fmt.Errorf("解压备份失败: %w", err)
	}

	rollbackPath := ""
	if current, err := os.Lstat(savePath); err == nil {
		if current.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("存档目录不能是符号链接")
		}
		if _, err := s.createBackup(saveName); err != nil {
			return fmt.Errorf("恢复前备份当前存档失败: %w", err)
		}
		rollbackPath = fmt.Sprintf("%s.restore-%d", savePath, time.Now().UnixNano())
		if err := os.Rename(savePath, rollbackPath); err != nil {
			return fmt.Errorf("暂存当前存档失败: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.Rename(staging, savePath); err != nil {
		if rollbackPath != "" {
			_ = os.Rename(rollbackPath, savePath)
		}
		return fmt.Errorf("替换存档失败: %w", err)
	}
	if rollbackPath != "" {
		_ = os.RemoveAll(rollbackPath)
	}

	return nil
}

// DeleteBackup 删除备份
func (s *SaveService) DeleteBackup(backupID int) error {
	s.backupMu.Lock()
	defer s.backupMu.Unlock()
	return s.deleteBackup(backupID)
}

func (s *SaveService) deleteBackup(backupID int) error {
	// 获取备份路径
	row := database.DB.QueryRow(`SELECT backup_path FROM backups WHERE id = ?`, backupID)
	var backupPath string
	if err := row.Scan(&backupPath); err != nil {
		return err
	}
	if !isWithin(backupPath, s.backupPath) {
		return fmt.Errorf("备份路径不在允许的目录内")
	}

	// 删除文件
	if err := os.Remove(backupPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除备份文件失败: %w", err)
	}

	// 从数据库删除
	_, err := database.DB.Exec(`DELETE FROM backups WHERE id = ?`, backupID)
	return err
}

// ListSaves 列出所有存档
func (s *SaveService) ListSaves() ([]string, error) {
	// 检查存档目录是否存在
	if _, err := os.Stat(s.savesPath); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(s.savesPath)
	if err != nil {
		return nil, err
	}

	saves := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") && entry.Name() != "backups" {
			saves = append(saves, entry.Name())
		}
	}

	return saves, nil
}

// zipDirectory 压缩目录
func (s *SaveService) zipDirectory(source, target string) error {
	zipFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}

	archive := zip.NewWriter(zipFile)
	walkErr := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("存档不允许包含符号链接: %s", path)
		}

		// 创建文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 设置相对路径
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		// 写入文件头
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		// 如果不是目录，写入文件内容
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(writer, file)
			closeErr := file.Close()
			if copyErr != nil {
				return copyErr
			}
			if closeErr != nil {
				return closeErr
			}
		}

		return nil
	})

	archiveErr := archive.Close()
	fileErr := zipFile.Close()
	if walkErr != nil {
		_ = os.Remove(target)
		return walkErr
	}
	if archiveErr != nil {
		_ = os.Remove(target)
		return archiveErr
	}
	if fileErr != nil {
		_ = os.Remove(target)
		return fileErr
	}
	return nil
}

// unzipDirectory 解压目录
func (s *SaveService) unzipDirectory(source, target string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()
	if len(reader.File) == 0 || len(reader.File) > maxBackupFiles {
		return fmt.Errorf("备份文件数量无效")
	}
	var totalSize uint64
	var extractedSize uint64
	for _, file := range reader.File {
		if file.UncompressedSize64 > maxBackupExtractSize-totalSize {
			return fmt.Errorf("备份解压后大小超过限制")
		}
		totalSize += file.UncompressedSize64
	}

	for _, file := range reader.File {
		if file.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("备份不允许包含符号链接")
		}
		name := filepath.Clean(filepath.FromSlash(strings.ReplaceAll(file.Name, "\\", "/")))
		path := filepath.Join(target, name)

		// 检查路径安全性
		if !isWithin(path, target) {
			return fmt.Errorf("非法路径: %s", path)
		}

		if file.FileInfo().IsDir() {
			perm := file.Mode().Perm()
			if perm == 0 {
				perm = 0755
			}
			if err := os.MkdirAll(path, perm); err != nil {
				return err
			}
			continue
		}

		// 创建父目录
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		// 创建文件
		perm := file.Mode().Perm()
		if perm == 0 {
			perm = 0644
		}
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		remaining := maxBackupExtractSize - extractedSize
		written, copyErr := io.Copy(outFile, io.LimitReader(rc, int64(remaining)+1))
		closeOutErr := outFile.Close()
		closeInErr := rc.Close()

		if copyErr != nil {
			return copyErr
		}
		if closeOutErr != nil {
			return closeOutErr
		}
		if closeInErr != nil {
			return closeInErr
		}
		if written < 0 || uint64(written) > remaining {
			return fmt.Errorf("备份解压后大小超过限制")
		}
		extractedSize += uint64(written)
	}

	return nil
}

func (s *SaveService) safeSavePath(saveName string) (string, error) {
	saveName = strings.TrimSpace(saveName)
	if saveName == "" || saveName == "." || saveName == ".." ||
		strings.ContainsAny(saveName, `/\\`) || filepath.Base(saveName) != saveName {
		return "", fmt.Errorf("无效的存档名称")
	}
	path := filepath.Join(s.savesPath, saveName)
	if !isWithin(path, s.savesPath) {
		return "", fmt.Errorf("存档路径不在允许的目录内")
	}
	return path, nil
}

// AutoBackup 自动备份任务
func (s *SaveService) AutoBackup() error {
	s.backupMu.Lock()
	defer s.backupMu.Unlock()

	// 列出所有存档
	saves, err := s.ListSaves()
	if err != nil {
		return err
	}

	// 为每个存档创建备份，并汇总失败，避免调用方误判为全部成功。
	var backupErrors []error
	for _, save := range saves {
		_, err := s.createBackup(save)
		if err != nil {
			backupErrors = append(backupErrors, fmt.Errorf("%s: %w", save, err))
		}
	}

	// 清理旧备份（保留最近30个）
	if err := s.cleanOldBackups(30); err != nil {
		backupErrors = append(backupErrors, fmt.Errorf("清理旧备份: %w", err))
	}
	return errors.Join(backupErrors...)
}

// StartAutoBackup 每隔 interval 执行一次由调用方协调过的安全备份。
func (s *SaveService) StartAutoBackup(interval time.Duration, run func() error) {
	if interval <= 0 || run == nil {
		return
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if err := run(); err != nil {
				log.Printf("自动备份跳过或失败: %v", err)
			}
		}
	}()
}

// cleanOldBackups 清理旧备份
func (s *SaveService) cleanOldBackups(keep int) error {
	// 获取所有备份，按时间倒序
	rows, err := database.DB.Query(`
		SELECT id
		FROM backups
		ORDER BY created_at DESC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// 删除超过保留数量的备份
	if len(ids) > keep {
		for _, id := range ids[keep:] {
			if err := s.deleteBackup(id); err != nil {
				return err
			}
		}
	}

	return nil
}
