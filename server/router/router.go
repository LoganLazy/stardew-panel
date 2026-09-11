package router

import (
	"log"
	"net/http"
	"os"
	"stardew-panel/config"
	"stardew-panel/database"
	"stardew-panel/handler"
	"stardew-panel/middleware"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS 中间件
func corsMiddleware(allowedOrigins string) gin.HandlerFunc {
	allowed := make(map[string]struct{})
	for _, origin := range strings.Split(allowedOrigins, ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			allowed[origin] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if _, ok := allowed[origin]; origin != "" && ok {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		} else if origin != "" && c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	trustedProxies := make([]string, 0)
	for _, proxy := range strings.Split(cfg.Server.TrustedProxies, ",") {
		if proxy = strings.TrimSpace(proxy); proxy != "" {
			trustedProxies = append(trustedProxies, proxy)
		}
	}
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		log.Printf("Invalid trusted_proxies configuration, proxy headers disabled: %v", err)
		_ = r.SetTrustedProxies(nil)
	}

	// 添加 CORS 中间件
	r.Use(corsMiddleware(cfg.Server.CORSOrigins))

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
			"version":   "0.4.3",
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
		// 安装检查（容器模式：引导 + 就绪状态）
		api.GET("/install/check", handler.CheckInstallation)
		api.GET("/install/setup", handler.GetSetupStatus)

		// 服务器管理
		server := api.Group("/server")
		{
			server.GET("/status", handler.GetServerStatus)
			server.POST("/start", handler.StartServer)
			server.POST("/stop", handler.StopServer)
			server.POST("/restart", handler.RestartServer)
			server.GET("/invite-code", handler.GetInviteCode)
			server.GET("/settings", handler.GetServerSettings)
			server.PUT("/settings/:key", handler.UpdateServerSetting)
			server.GET("/metrics", handler.GetMetrics)
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
		api.GET("/logs/stream", handler.StreamLogs)
	}

	// 静态文件服务（生产环境）
	// 检查是否存在 web/dist 目录
	if _, err := os.Stat("./web/dist"); err == nil {
		r.Static("/assets", "./web/dist/assets")
		r.StaticFile("/", "./web/dist/index.html")
		r.NoRoute(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{"error": "接口不存在"})
				return
			}
			c.File("./web/dist/index.html")
		})
	}

	return r
}
