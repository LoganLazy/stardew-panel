package handler

import (
	"net/http"
	"stardew-panel/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var playerService *service.PlayerService

// InitPlayerHandler 初始化玩家处理器
func InitPlayerHandler() {
	playerService = service.NewPlayerService()
}

// ListPlayers 获取在线玩家列表
func ListPlayers(c *gin.Context) {
	// 先解析日志更新玩家状态
	playerService.ParseLogForPlayers()

	// 获取在线玩家
	players, err := playerService.GetOnlinePlayers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取统计信息
	stats, err := playerService.GetPlayerStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 计算在线时长
	for i := range players {
		duration := time.Since(players[i].ConnectedAt)
		players[i].LastSeen = time.Now()
	}

	c.JSON(http.StatusOK, gin.H{
		"players": players,
		"stats":   stats,
	})
}

// RefreshPlayers 手动刷新玩家列表（解析日志）
func RefreshPlayers(c *gin.Context) {
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

	// TODO: 实现真正的踢人逻辑（需要游戏服务器API支持）
	if err := playerService.KickPlayer(req.PlayerName); err != nil {
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
