package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckInstallation 检查服务端是否已安装
func CheckInstallation(c *gin.Context) {
	// TODO: 实现检查逻辑
	c.JSON(http.StatusOK, gin.H{
		"installed": false,
	})
}

// InstallServer 安装服务端
func InstallServer(c *gin.Context) {
	var req struct {
		Mode string `json:"mode"` // official 或 smapi
		Path string `json:"path"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现安装逻辑
	// 1. 下载 SteamCMD
	// 2. 下载星露谷服务端
	// 3. 如果是 SMAPI 模式，安装 SMAPI
	// 4. 配置服务器

	c.JSON(http.StatusOK, gin.H{
		"message": "Installation started",
		"mode":    req.Mode,
		"path":    req.Path,
	})
}

// GetInstallStatus 获取安装状态
func GetInstallStatus(c *gin.Context) {
	// TODO: 实现获取安装进度
	c.JSON(http.StatusOK, gin.H{
		"status":   "installing",
		"progress": 50,
		"message":  "Downloading server files...",
	})
}
