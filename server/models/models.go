package models

import "time"

// Installation 安装配置
type Installation struct {
	ID            int       `json:"id"`
	GamePath      string    `json:"game_path"`
	InstallMethod string    `json:"install_method"` // upload, path, steamcmd
	HasSMAPI      bool      `json:"has_smapi"`
	Version       string    `json:"version"`
	InstalledAt   time.Time `json:"installed_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Mod MOD信息
type Mod struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	UniqueID    string    `json:"unique_id"`
	Version     string    `json:"version"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	FilePath    string    `json:"file_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OnlinePlayer 在线玩家
type OnlinePlayer struct {
	ID          int       `json:"id"`
	PlayerName  string    `json:"player_name"`
	ConnectedAt time.Time `json:"connected_at"`
	LastSeen    time.Time `json:"last_seen"`
	IsOnline    bool      `json:"is_online"`
}

// Backup 存档备份
type Backup struct {
	ID         int       `json:"id"`
	SaveName   string    `json:"save_name"`
	BackupPath string    `json:"backup_path"`
	Size       int64     `json:"size"`
	CreatedAt  time.Time `json:"created_at"`
}

// ServerState 服务器状态
type ServerState struct {
	ID        int        `json:"id"`
	Status    string     `json:"status"` // stopped, starting, running, stopping
	PID       int        `json:"pid"`
	StartedAt *time.Time `json:"started_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ServerStatus 服务器状态响应
type ServerStatus struct {
	Status        string   `json:"status"`
	OnlinePlayers int      `json:"online_players"`
	Version       string   `json:"version"`
	Uptime        int64    `json:"uptime"` // 秒
	Players       []string `json:"players"`
}

// InstallRequest 安装请求
type InstallRequest struct {
	Method       string `json:"method"`        // upload, path, steamcmd
	Path         string `json:"path"`          // 用于 path 和 steamcmd 方法
	InstallSMAPI bool   `json:"install_smapi"`
	DetectSMAPI  bool   `json:"detect_smapi"`
	// SteamCMD 相关
	SteamUsername string `json:"steam_username"` // Steam 用户名
	SteamPassword string `json:"steam_password"` // Steam 密码
}
