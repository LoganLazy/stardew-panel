package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Init 初始化数据库
func Init(dbPath string) error {
	// 确保数据目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// 打开数据库连接
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	DB = db

	// 创建表
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// createTables 创建数据库表
func createTables() error {
	schemas := []string{
		// 用户表（认证）
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// 会话表（持久化）
		`CREATE TABLE IF NOT EXISTS sessions (
			token TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL,
			ip_address TEXT,
			user_agent TEXT
		)`,

		// 安装配置表
		`CREATE TABLE IF NOT EXISTS installation (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_path TEXT NOT NULL,
			install_method TEXT NOT NULL,
			has_smapi BOOLEAN NOT NULL DEFAULT 0,
			version TEXT,
			installed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// MOD 表
		`CREATE TABLE IF NOT EXISTS mods (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			unique_id TEXT UNIQUE,
			version TEXT,
			author TEXT,
			description TEXT,
			enabled BOOLEAN NOT NULL DEFAULT 1,
			file_path TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// 在线玩家记录表
		`CREATE TABLE IF NOT EXISTS online_players (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			player_name TEXT NOT NULL,
			connected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_online BOOLEAN NOT NULL DEFAULT 1
		)`,

		// 存档备份表
		`CREATE TABLE IF NOT EXISTS backups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			save_name TEXT NOT NULL,
			backup_path TEXT NOT NULL,
			size INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// 服务器状态表
		`CREATE TABLE IF NOT EXISTS server_state (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			status TEXT NOT NULL DEFAULT 'stopped',
			pid INTEGER,
			started_at DATETIME,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, schema := range schemas {
		if _, err := DB.Exec(schema); err != nil {
			return err
		}
	}

	// 初始化服务器状态
	_, err := DB.Exec(`INSERT OR IGNORE INTO server_state (id, status) VALUES (1, 'stopped')`)
	return err
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
