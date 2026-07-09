package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"stardew-panel/config"
	"stardew-panel/models"
	"stardew-panel/service"

	"github.com/gin-gonic/gin"
)

var installService *service.InstallService

// InitInstallHandler 初始化安装处理器
func InitInstallHandler(cfg *config.Config) {
	installService = service.NewInstallService(cfg.Game.ServerPath, cfg.Game.ModsPath, cfg.Game.SavesPath)
}

// CheckInstallation 检查安装状态
func CheckInstallation(c *gin.Context) {
	inst, err := installService.CheckInstallation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if inst == nil {
		c.JSON(http.StatusOK, gin.H{"installed": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"installed":    true,
		"installation": inst,
	})
}

// VerifyGamePath 验证游戏路径
func VerifyGamePath(c *gin.Context) {
	var req models.InstallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if req.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "路径不能为空"})
		return
	}

	inst, err := installService.VerifyGamePath(req.Path, req.DetectSMAPI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"installation": inst,
	})
}

// UploadGameFiles 上传游戏文件
func UploadGameFiles(c *gin.Context) {
	// 限制上传文件大小为 5GB（游戏文件较大）
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5*1024*1024*1024)

	file, err := c.FormFile("file")
	if err != nil {
		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件过大，最大支持 5GB"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败"})
		return
	}

	if filepath.Ext(file.Filename) != ".zip" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持 .zip 格式"})
		return
	}

	tempDir := "./data/temp"
	os.MkdirAll(tempDir, 0755)
	tempPath := filepath.Join(tempDir, file.Filename)

	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	installSMAPI := c.PostForm("installSMAPI") == "true"

	inst, err := installService.ExtractUploadedFile(tempPath, installSMAPI)
	if err != nil {
		os.Remove(tempPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"installation": inst,
	})
}

// InstallViaSteamCMD 通过 SteamCMD 安装
func InstallViaSteamCMD(c *gin.Context) {
	var req models.InstallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if req.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "安装路径不能为空"})
		return
	}

	if req.SteamUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Steam 用户名不能为空"})
		return
	}

	if req.SteamPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Steam 密码不能为空"})
		return
	}

	inst, err := installService.InstallViaSteamCMD(req.Path, req.SteamUsername, req.SteamPassword, req.InstallSMAPI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"installation": inst,
	})
}

// GetInstallStatus 获取安装进度
func GetInstallStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":   "idle",
		"progress": 100,
		"message":  "安装完成",
	})
}
