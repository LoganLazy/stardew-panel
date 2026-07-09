package handler

import (
	"net/http"
	"stardew-panel/config"
	"stardew-panel/service"

	"github.com/gin-gonic/gin"
)

var serverService *service.ServerService

// InitServerHandler 初始化服务器处理器
func InitServerHandler(cfg *config.Config) {
	serverService = service.NewServerService(cfg.Game.ServerPath)
}

// GetServerStatus 获取服务器状态
func GetServerStatus(c *gin.Context) {
	status, err := serverService.GetStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// StartServer 启动服务器
func StartServer(c *gin.Context) {
	if err := serverService.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "服务器启动成功",
	})
}

// StopServer 停止服务器
func StopServer(c *gin.Context) {
	if err := serverService.Stop(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "服务器已停止",
	})
}

// RestartServer 重启服务器
func RestartServer(c *gin.Context) {
	if err := serverService.Restart(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "服务器重启成功",
	})
}
