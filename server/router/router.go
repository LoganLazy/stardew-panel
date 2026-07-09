package router

import (
	"net/http"
	"os"
	"stardew-panel/config"
	"stardew-panel/database"
	"stardew-panel/handler"
	"stardew-panel/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(corsMiddleware())

	// 添加登录限流中间件
	r.Use(middleware.LoginRateLimit())

	// 健康检查（无需认证）
	r.GET("/health", func(c *gin.Context) {
		// 检查数据库连接
		dbOk := database.DB.Ping() == nil

		status := "ok"
		if !dbOk {
			status = "degraded"
		}

		c.JSON(200, gin.H{
			"status":    status,
			"version":   "0.3.0",
			"database":  dbOk,
			"timestamp": time.Now().Unix(),
		})
	})

	// API 路由组
	api := r.Group("/api/v1")

	// 认证路由（无需认证）
	auth := api.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/logout", handler.Logout)
		auth.GET("/check", handler.AuthMiddleware(), handler.CheckAuth)
		auth.POST("/change-password", handler.AuthMiddleware(), handler.ChangePassword)
	}

	// 其他所有路由需要认证
	api.Use(handler.AuthMiddleware())
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
			players.POST("/refresh", handler.RefreshPlayers)
			players.POST("/kick", handler.KickPlayer)
			players.DELETE("/cleanup", handler.CleanupPlayers)
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

	// 静态文件服务（生产环境）
	// 检查是否存在 web/dist 目录
	if _, err := os.Stat("./web/dist"); err == nil {
		r.Static("/assets", "./web/dist/assets")
		r.StaticFile("/", "./web/dist/index.html")
		r.NoRoute(func(c *gin.Context) {
			c.File("./web/dist/index.html")
		})
	}

	return r
}
