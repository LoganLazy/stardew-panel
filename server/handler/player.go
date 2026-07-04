package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListPlayers 获取玩家列表
func ListPlayers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"online":    []gin.H{},
		"whitelist": []gin.H{},
	})
}

// AddToWhitelist 添加到白名单
func AddToWhitelist(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Player added to whitelist",
	})
}

// RemoveFromWhitelist 从白名单移除
func RemoveFromWhitelist(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Player removed from whitelist",
		"id":      id,
	})
}

// KickPlayer 踢出玩家
func KickPlayer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Player kicked",
	})
}
