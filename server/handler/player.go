package handler

import (
	"errors"
	"net/http"
	"stardew-panel/config"
	"stardew-panel/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var playerService *service.PlayerService

// InitPlayerHandler 初始化玩家处理器
func InitPlayerHandler(cfg *config.Config) {
	playerService = service.NewPlayerService(cfg.Game.GameAPIURL, cfg.Game.GameAPIKey)
}

// ListPlayers 获取在线玩家列表。
// 优先从 sdvd 游戏容器 API 拿实时数据；API 不可用时回退到日志解析 + 数据库。
func ListPlayers(c *gin.Context) {
	// 优先：游戏容器实时玩家
	if live, err := playerService.GetLivePlayers(); err == nil {
		stats, statsErr := playerService.GetPlayerStats()
		if statsErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": statsErr.Error()})
			return
		}
		stats["source"] = "game_api"
		c.JSON(http.StatusOK, gin.H{
			"players": live,
			"stats":   stats,
		})
		return
	}

	// 回退：日志解析 + 数据库缓存
	playerService.ParseLogForPlayers()

	players, err := playerService.GetOnlinePlayers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stats, err := playerService.GetPlayerStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range players {
		players[i].LastSeen = time.Now()
	}

	c.JSON(http.StatusOK, gin.H{
		"players": players,
		"stats":   stats,
	})
}

// RefreshPlayers 手动刷新玩家列表（解析日志）
func RefreshPlayers(c *gin.Context) {
	if _, err := playerService.GetLivePlayers(); err == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "玩家列表已刷新"})
		return
	}
	if err := playerService.ParseLogForPlayers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "玩家列表已刷新",
	})
}

// KickPlayer 踢出玩家
func KickPlayer(c *gin.Context) {
	var req struct {
		PlayerName string `json:"player_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if req.PlayerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "玩家名称不能为空"})
		return
	}

	if err := playerService.KickPlayer(req.PlayerName); err != nil {
		if errors.Is(err, service.ErrKickPlayerUnsupported) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "玩家已被踢出",
	})
}

// CleanupPlayers 清理离线玩家记录
func CleanupPlayers(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			if parsed < 0 || parsed > 3650 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "days 必须在 0 到 3650 之间"})
				return
			}
			days = parsed
		}
	}

	if err := playerService.CleanupOfflinePlayers(days); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "清理完成",
	})
}

// 保留旧的白名单接口以兼容性，但标记为废弃
// AddToWhitelist - 已废弃
func AddToWhitelist(c *gin.Context) {
	c.JSON(http.StatusGone, gin.H{"error": "白名单功能已移除"})
}

// RemoveFromWhitelist - 已废弃
func RemoveFromWhitelist(c *gin.Context) {
	c.JSON(http.StatusGone, gin.H{"error": "白名单功能已移除"})
}
