package router

import (
	"stardew-panel/config"
	"stardew-panel/handler"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"version": "0.1.0",
		})
	})

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 安装检查
		api.GET("/install/check", handler.CheckInstallation)
		api.POST("/install/upload", handler.UploadGameFiles)
		api.POST("/install/verify", handler.VerifyGamePath)
		api.POST("/install/steamcmd", handler.InstallViaSteamCMD)
		api.GET("/install/status", handler.GetInstallStatus)

		// 服务器管理
		server := api.Group("/server")
		{
			server.GET("/status", handler.GetServerStatus)
			server.POST("/start", handler.StartServer)
			server.POST("/stop", handler.StopServer)
			server.POST("/restart", handler.RestartServer)
		}

		// MOD 管理
		mods := api.Group("/mods")
		{
			mods.GET("", handler.ListMods)
			mods.POST("/upload", handler.UploadMod)
			mods.PUT("/:id/toggle", handler.ToggleMod)
			mods.DELETE("/:id", handler.DeleteMod)
		}

		// 玩家管理
		players := api.Group("/players")
		{
			players.GET("", handler.ListPlayers)
			players.POST("/whitelist", handler.AddToWhitelist)
			players.DELETE("/whitelist/:id", handler.RemoveFromWhitelist)
			players.POST("/kick", handler.KickPlayer)
		}

		// 存档管理
		saves := api.Group("/saves")
		{
			saves.GET("", handler.ListSaves)
			saves.POST("/backup", handler.CreateBackup)
			saves.POST("/restore/:id", handler.RestoreSave)
			saves.DELETE("/:id", handler.DeleteSave)
		}

		// 日志
		api.GET("/logs", handler.GetLogs)
	}

	return r
}
