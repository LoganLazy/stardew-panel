package service

import (
	"bufio"
	"database/sql"
	"os"
	"regexp"
	"stardew-panel/database"
	"stardew-panel/models"
	"strings"
	"time"
)

// PlayerService 玩家管理服务
type PlayerService struct{}

// NewPlayerService 创建玩家服务
func NewPlayerService() *PlayerService {
	return &PlayerService{}
}

// GetOnlinePlayers 获取在线玩家列表
func (s *PlayerService) GetOnlinePlayers() ([]models.OnlinePlayer, error) {
	// 从数据库获取在线玩家
	rows, err := database.DB.Query(`
		SELECT id, player_name, connected_at, last_seen, is_online
		FROM online_players
		WHERE is_online = 1
		ORDER BY connected_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []models.OnlinePlayer
	for rows.Next() {
		var player models.OnlinePlayer
		err := rows.Scan(&player.ID, &player.PlayerName, &player.ConnectedAt, &player.LastSeen, &player.IsOnline)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

// ParseLogForPlayers 从日志文件解析玩家连接信息
func (s *PlayerService) ParseLogForPlayers() error {
	logPath := "./data/server.log"

	// 检查日志文件是否存在
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return nil // 日志文件不存在，跳过
	}

	file, err := os.Open(logPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取文件大小
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// 如果文件很大（>10MB），只读取最后 1MB
	const maxReadSize = 1024 * 1024 // 1MB
	fileSize := stat.Size()

	if fileSize > 10*1024*1024 {
		// 大文件：只读取最后 1MB
		offset := fileSize - maxReadSize
		if offset < 0 {
			offset = 0
		}
		_, err = file.Seek(offset, 0)
		if err != nil {
			return err
		}
	}

	// 正则表达式匹配玩家加入/离开
	// 修改为支持中文和特殊字符
	joinPattern := regexp.MustCompile(`(.+?) joined the game`)
	leftPattern := regexp.MustCompile(`(.+?) left the game`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 检查玩家加入
		if matches := joinPattern.FindStringSubmatch(line); matches != nil {
			playerName := strings.TrimSpace(matches[1])
			if playerName != "" {
				s.playerJoined(playerName)
			}
		}

		// 检查玩家离开
		if matches := leftPattern.FindStringSubmatch(line); matches != nil {
			playerName := strings.TrimSpace(matches[1])
			if playerName != "" {
				s.playerLeft(playerName)
			}
		}
	}

	return scanner.Err()
}

// playerJoined 玩家加入游戏
func (s *PlayerService) playerJoined(playerName string) error {
	// 检查玩家是否已存在
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM online_players
		WHERE player_name = ? AND is_online = 1
	`, playerName).Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 {
		// 已在线，更新最后看到时间
		_, err = database.DB.Exec(`
			UPDATE online_players
			SET last_seen = CURRENT_TIMESTAMP
			WHERE player_name = ? AND is_online = 1
		`, playerName)
		return err
	}

	// 新加入，插入记录
	_, err = database.DB.Exec(`
		INSERT INTO online_players (player_name, connected_at, last_seen, is_online)
		VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 1)
	`, playerName)

	return err
}

// playerLeft 玩家离开游戏
func (s *PlayerService) playerLeft(playerName string) error {
	_, err := database.DB.Exec(`
		UPDATE online_players
		SET is_online = 0, last_seen = CURRENT_TIMESTAMP
		WHERE player_name = ? AND is_online = 1
	`, playerName)
	return err
}

// KickPlayer 踢出玩家
func (s *PlayerService) KickPlayer(playerName string) error {
	// TODO: 需要游戏服务器API支持
	// 这里只是标记为离线
	return s.playerLeft(playerName)
}

// GetPlayerStats 获取玩家统计
func (s *PlayerService) GetPlayerStats() (map[string]interface{}, error) {
	// 在线人数
	var onlineCount int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM online_players WHERE is_online = 1
	`).Scan(&onlineCount)
	if err != nil {
		return nil, err
	}

	// 历史总人数
	var totalCount int
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT player_name) FROM online_players
	`).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	return map[string]interface{
		"online_count": onlineCount,
		"total_count":  totalCount,
	}, nil
}

// CleanupOfflinePlayers 清理离线时间过长的记录
func (s *PlayerService) CleanupOfflinePlayers(days int) error {
	_, err := database.DB.Exec(`
		DELETE FROM online_players
		WHERE is_online = 0
		AND last_seen < datetime('now', '-' || ? || ' days')
	`, days)
	return err
}
