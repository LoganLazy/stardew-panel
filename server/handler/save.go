package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListSaves 获取存档列表
func ListSaves(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"saves": []gin.H{},
	})
}

// CreateBackup 创建备份
func CreateBackup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Backup created successfully",
	})
}

// RestoreSave 恢复存档
func RestoreSave(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Save restored",
		"id":      id,
	})
}

// DeleteSave 删除存档
func DeleteSave(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Save deleted",
		"id":      id,
	})
}

// GetLogs 获取日志
func GetLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"logs": []string{},
	})
}
