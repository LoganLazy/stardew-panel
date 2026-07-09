package handler

import (
	"net/http"
	"stardew-panel/config"
	"stardew-panel/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var saveService *service.SaveService

// InitSaveHandler 初始化存档处理器
func InitSaveHandler(cfg *config.Config) {
	saveService = service.NewSaveService(cfg.Game.SavesPath)
}

// ListSaves 获取存档列表
func ListSaves(c *gin.Context) {
	// 获取备份列表
	backups, err := saveService.ListBackups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取存档列表
	saves, err := saveService.ListSaves()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"saves":   saves,
		"backups": backups,
	})
}

// CreateBackup 创建备份
func CreateBackup(c *gin.Context) {
	var req struct {
		SaveName string `json:"save_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if req.SaveName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "存档名称不能为空"})
		return
	}

	backup, err := saveService.CreateBackup(req.SaveName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"backup":  backup,
	})
}

// RestoreSave 恢复存档
func RestoreSave(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的备份ID"})
		return
	}

	if err := saveService.RestoreBackup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "存档恢复成功",
	})
}

// DeleteSave 删除存档
func DeleteSave(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的备份ID"})
		return
	}

	if err := saveService.DeleteBackup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "备份已删除",
	})
}

// GetLogs 获取日志
func GetLogs(c *gin.Context) {
	lines := 100
	if l := c.Query("lines"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			lines = parsed
		}
	}

	// 使用 serverService 获取日志
	logs, err := serverService.GetLogs(lines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}
