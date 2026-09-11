package service

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"stardew-panel/database"
	"stardew-panel/models"
	"strings"
)

var ErrKickPlayerUnsupported = errors.New("当前游戏服务器 API 不支持踢出玩家")

// PlayerService 玩家管理服务。
// 改造后优先从 sdvd 游戏容器 REST API 拿实时玩家（gameAPI），
// 拿不到时回退到日志解析 + 数据库缓存（老逻辑保留作兜底）。
type PlayerService struct {
	gameAPI *GameAPIClient
}

// NewPlayerService 创建玩家服务。
// gameAPIURL/gameAPIKey 为空时退化为纯日志/数据库模式。
func NewPlayerService(gameAPIURL, gameAPIKey string) *PlayerService {
	return &PlayerService{
		gameAPI: NewGameAPIClient(gameAPIURL, gameAPIKey),
	}
}

// GetLivePlayers 从 sdvd 游戏容器 API 拿实时在线玩家。
// 返回的是玩家名列表；API 不可用时返回 error，由调用方回退到数据库。
func (s *PlayerService) GetLivePlayers() ([]models.OnlinePlayer, error) {
	apiPlayers, err := s.gameAPI.GetPlayers()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(apiPlayers))
	for _, p := range apiPlayers {
		if !p.IsOnline {
			continue
		}
		names = append(names, p.Name)
	}
	if err := s.syncLivePlayers(names); err != nil {
		return nil, err
	}
	return s.GetOnlinePlayers()
}

// syncLivePlayers 将游戏 API 快照写入历史表，同时保留仍在线玩家的加入时间。
func (s *PlayerService) syncLivePlayers(names []string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT player_name FROM online_players WHERE is_online = 1`)
	if err != nil {
		return err
	}
	current := make(map[string]struct{})
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			rows.Close()
			return err
		}
		current[name] = struct{}{}
	}
	if err := rows.Close(); err != nil {
		return err
	}

	live := make(map[string]struct{}, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		live[name] = struct{}{}
		if _, ok := current[name]; ok {
			if _, err := tx.Exec(`UPDATE online_players SET last_seen = CURRENT_TIMESTAMP WHERE player_name = ? AND is_online = 1`, name); err != nil {
				return err
			}
			continue
		}
		if _, err := tx.Exec(`
			INSERT INTO online_players (player_name, connected_at, last_seen, is_online)
			VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 1)
		`, name); err != nil {
			return err
		}
	}

	for name := range current {
		if _, ok := live[name]; ok {
			continue
		}
		if _, err := tx.Exec(`
			UPDATE online_players SET is_online = 0, last_seen = CURRENT_TIMESTAMP
			WHERE player_name = ? AND is_online = 1
		`, name); err != nil {
			return err
		}
	}
	return tx.Commit()
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

	players := make([]models.OnlinePlayer, 0)
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
	return ErrKickPlayerUnsupported
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

	return map[string]interface{}{
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
