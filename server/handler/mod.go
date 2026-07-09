package handler

import (
	"net/http"
	"path/filepath"
	"stardew-panel/config"
	"stardew-panel/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var modService *service.ModService

// InitModHandler 初始化MOD处理器
func InitModHandler(cfg *config.Config) {
	modService = service.NewModService(cfg.Game.ModsPath)
}

// ListMods 获取 MOD 列表
func ListMods(c *gin.Context) {
	// 先扫描文件系统同步MOD
	modService.ScanMods()

	mods, err := modService.ListMods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mods": mods,
	})
}

// UploadMod 上传 MOD
func UploadMod(c *gin.Context) {
	// 限制上传文件大小为 500MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 500*1024*1024)

	file, err := c.FormFile("file")
	if err != nil {
		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件过大，最大支持 500MB"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败"})
		return
	}

	// 保存临时文件
	tempPath := "./data/temp/" + file.Filename
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 安装MOD
	mod, err := modService.UploadMod(tempPath, file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"mod":     mod,
	})
}

// ToggleMod 启用/禁用 MOD
func ToggleMod(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的MOD ID"})
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := modService.ToggleMod(id, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "MOD状态已更新",
	})
}

// DeleteMod 删除 MOD
func DeleteMod(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的MOD ID"})
		return
	}

	if err := modService.DeleteMod(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "MOD已删除",
	})
}
