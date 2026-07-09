package service

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"stardew-panel/database"
	"stardew-panel/models"
	"time"
)

// SaveService 存档管理服务
type SaveService struct {
	savesPath  string
	backupPath string
}

// NewSaveService 创建存档服务
func NewSaveService(savesPath string) *SaveService {
	backupPath := filepath.Join(savesPath, "backups")
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

	var backups []models.Backup
	for rows.Next() {
		var backup models.Backup
		err := rows.Scan(&backup.ID, &backup.SaveName, &backup.BackupPath, &backup.Size, &backup.CreatedAt)
		if err != nil {
			return nil, err
		}
		backups = append(backups, backup)
	}

	return backups, nil
}

// CreateBackup 创建存档备份
func (s *SaveService) CreateBackup(saveName string) (*models.Backup, error) {
	// 检查存档是否存在
	savePath := filepath.Join(s.savesPath, saveName)
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("存档不存在: %s", saveName)
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
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

	// 检查备份文件是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在")
	}

	// 目标路径
	savePath := filepath.Join(s.savesPath, saveName)

	// 备份当前存档（如果存在）
	if _, err := os.Stat(savePath); err == nil {
		timestamp := time.Now().Format("20060102_150405")
		tempBackup := filepath.Join(s.backupPath, fmt.Sprintf("%s_before_restore_%s.zip", saveName, timestamp))
		s.zipDirectory(savePath, tempBackup)
		// 删除当前存档
		os.RemoveAll(savePath)
	}

	// 解压备份
	if err := s.unzipDirectory(backupPath, savePath); err != nil {
		return fmt.Errorf("解压备份失败: %w", err)
	}

	return nil
}

// DeleteBackup 删除备份
func (s *SaveService) DeleteBackup(backupID int) error {
	// 获取备份路径
	row := database.DB.QueryRow(`SELECT backup_path FROM backups WHERE id = ?`, backupID)
	var backupPath string
	if err := row.Scan(&backupPath); err != nil {
		return err
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

	var saves []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "backups" {
			saves = append(saves, entry.Name())
		}
	}

	return saves, nil
}

// zipDirectory 压缩目录
func (s *SaveService) zipDirectory(source, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
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
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

// unzipDirectory 解压目录
func (s *SaveService) unzipDirectory(source, target string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)

		// 检查路径安全性
		if !filepath.HasPrefix(path, filepath.Clean(target)+string(os.PathSeparator)) {
			return fmt.Errorf("非法路径: %s", path)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		// 创建父目录
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		// 创建文件
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// AutoBackup 自动备份任务
func (s *SaveService) AutoBackup() error {
	// 列出所有存档
	saves, err := s.ListSaves()
	if err != nil {
		return err
	}

	// 为每个存档创建备份
	for _, save := range saves {
		_, err := s.CreateBackup(save)
		if err != nil {
			fmt.Printf("自动备份失败 [%s]: %v\n", save, err)
		}
	}

	// 清理旧备份（保留最近30个）
	return s.cleanOldBackups(30)
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
		rows.Scan(&id)
		ids = append(ids, id)
	}

	// 删除超过保留数量的备份
	if len(ids) > keep {
		for _, id := range ids[keep:] {
			s.DeleteBackup(id)
		}
	}

	return nil
}
