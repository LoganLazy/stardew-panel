package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListMods 获取 MOD 列表
func ListMods(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"mods": []gin.H{},
	})
}

// UploadMod 上传 MOD
func UploadMod(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "MOD uploaded successfully",
	})
}

// ToggleMod 启用/禁用 MOD
func ToggleMod(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "MOD toggled",
		"id":      id,
	})
}

// DeleteMod 删除 MOD
func DeleteMod(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "MOD deleted",
		"id":      id,
	})
}
