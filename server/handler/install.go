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
		"path":      "",
		"type":      "",
		"version":   "",
	})
}

// UploadGameFiles 上传游戏文件
func UploadGameFiles(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	installSMAPI := c.PostForm("installSMAPI") == "true"

	// TODO: 实现上传和解压逻辑
	// 1. 保存上传的 zip 文件
	// 2. 解压到指定目录
	// 3. 验证游戏文件完整性
	// 4. 如果需要，安装 SMAPI
	// 5. 保存配置

	c.JSON(http.StatusOK, gin.H{
		"message":     "Upload started",
		"filename":    file.Filename,
		"size":        file.Size,
		"installSMAPI": installSMAPI,
	})
}

// VerifyGamePath 验证游戏路径
func VerifyGamePath(c *gin.Context) {
	var req struct {
		Path        string `json:"path"`
		DetectSMAPI bool   `json:"detectSMAPI"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现路径验证逻辑
	// 1. 检查路径是否存在
	// 2. 检查是否包含 StardewValley 可执行文件
	// 3. 如果需要，检测 SMAPI
	// 4. 获取游戏版本
	// 5. 保存配置

	c.JSON(http.StatusOK, gin.H{
		"valid":        true,
		"path":         req.Path,
		"hasSMAPI":     false,
		"version":      "1.6.9",
		"executable":   "StardewValley",
	})
}

// InstallViaSteamCMD 通过 SteamCMD 安装
func InstallViaSteamCMD(c *gin.Context) {
	var req struct {
		Path        string `json:"path"`
		InstallSMAPI bool   `json:"installSMAPI"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现 SteamCMD 下载逻辑
	// 1. 检查 SteamCMD 是否已安装
	// 2. 如果没有，下载并安装 SteamCMD
	// 3. 使用 SteamCMD 下载游戏 (App ID: 413150)
	// 4. 如果需要，下载并安装 SMAPI
	// 5. 保存配置

	c.JSON(http.StatusOK, gin.H{
		"message":     "SteamCMD installation started",
		"path":        req.Path,
		"installSMAPI": req.InstallSMAPI,
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
