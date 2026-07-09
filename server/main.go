package main

import (
	"log"
	"stardew-panel/config"
	"stardew-panel/database"
	"stardew-panel/handler"
	"stardew-panel/router"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := database.Init(cfg.Database.Path); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 初始化默认用户（首次启动）
	authService := handler.InitAuthHandler()
	if err := authService.InitDefaultUser(); err != nil {
		log.Printf("Warning: Failed to init default user: %v", err)
	}

	// 启动会话清理任务
	authService.CleanupExpiredSessions()

	// 初始化所有 handler
	handler.InitInstallHandler(cfg)
	handler.InitServerHandler(cfg)
	handler.InitModHandler(cfg)
	handler.InitPlayerHandler()
	handler.InitSaveHandler(cfg)

	// 初始化路由
	r := router.Setup(cfg)

	// 启动服务器
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("StardewPanel server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
