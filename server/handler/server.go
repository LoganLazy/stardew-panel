package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"stardew-panel/config"
	"stardew-panel/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var serverService *service.ServerService

// InitServerHandler 初始化服务器处理器
func InitServerHandler(cfg *config.Config) {
	serverService = service.NewServerService(
		cfg.Game.ComposeFile,
		cfg.Game.ComposeDir,
		cfg.Game.GameAPIURL,
		cfg.Game.GameAPIKey,
		cfg.Game.GamePort,
		cfg.Game.VNCPort,
	)
}

// GetServerStatus 获取服务器状态
func GetServerStatus(c *gin.Context) {
	status, err := serverService.GetStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// StartServer 启动服务器
func StartServer(c *gin.Context) {
	if err := serverService.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "服务器启动请求已提交",
	})
}

// StopServer 停止服务器
func StopServer(c *gin.Context) {
	warning, err := serverService.StopThen(saveService.AutoBackup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	warningMessage := ""
	if warning != nil {
		warningMessage = warning.Error()
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "服务器已停止",
		"warning": warningMessage,
	})
}

// RestartServer 重启服务器
func RestartServer(c *gin.Context) {
	if err := serverService.Restart(saveService.AutoBackup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "服务器重启请求已提交",
	})
}

// GetInviteCode 返回联机邀请码（供前端显示，玩家用它加入）。
// 容器未就绪或游戏 API 未响应时返回空码，前端据此显示"暂不可用"。
func GetInviteCode(c *gin.Context) {
	code, err := serverService.GetInviteCode()
	if err != nil {
		// 不作为错误抛出：容器没起来时拿不到码是正常的
		c.JSON(http.StatusOK, gin.H{"invite_code": "", "available": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"invite_code": code,
		"available":   code != "",
	})
}

// GetServerSettings 代理 sdvd GET /settings，返回全部服务器设置。
// 容器未就绪时返回空对象 + available:false，前端据此禁用设置表单。
func GetServerSettings(c *gin.Context) {
	settings, err := serverService.GetSettings()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"settings": gin.H{}, "available": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"settings": settings, "available": true})
}

// UpdateServerSetting 代理 sdvd PUT /settings/{key}，更新单个设置项。
// 请求体 {"value": <新值>}，key 从 URL 参数取。
func UpdateServerSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "设置项 key 不能为空"})
		return
	}

	var req struct {
		Value interface{} `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}
	if err := validateServerSetting(key, req.Value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := serverService.UpdateSetting(key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "设置已更新"})
}

// GetMetrics 性能监控数据（容器资源占用 + 面板自身运行时指标）。
// 容器未运行时 container 为 null，前端据此显示占位提示。
func GetMetrics(c *gin.Context) {
	metrics, err := serverService.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}

// StreamLogs 实时日志流（SSE）。
// 认证是 Bearer 头，原生 EventSource 带不了自定义头，前端用 fetch 流式读取。
// 事件类型：log（一行日志）、error（流异常）、end（流结束）。
func StreamLogs(c *gin.Context) {
	tail := 200
	if l := c.Query("lines"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			if parsed < 1 || parsed > 2000 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "lines 必须在 1 到 2000 之间"})
				return
			}
			tail = parsed
		}
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	// 告知 nginx 等反代关闭缓冲，否则日志会被攒着不推
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.Flush()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	lines := make(chan string, 64)
	streamErr := make(chan error, 1)
	go func() {
		defer close(lines)
		err := serverService.StreamLogs(ctx, tail, func(line string) {
			select {
			case lines <- line:
			case <-ctx.Done():
			}
		})
		if err != nil {
			streamErr <- err
		}
	}()

	writeSSE(c, "log-ready", gin.H{"tail": tail})

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			// 心跳注释行，防止空闲连接被代理掐断
			fmt.Fprintf(c.Writer, ": ping\n\n")
			c.Writer.Flush()
		case line, ok := <-lines:
			if !ok {
				select {
				case err := <-streamErr:
					writeSSE(c, "error", gin.H{"error": err.Error()})
				default:
				}
				writeSSE(c, "end", gin.H{})
				c.Writer.Flush()
				return
			}
			writeSSE(c, "log", gin.H{"line": line})
			c.Writer.Flush()
		}
	}
}

// writeSSE 写一条 SSE 事件并立即刷出。
func writeSSE(c *gin.Context, event string, data interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		return
	}
	fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, payload)
	c.Writer.Flush()
}

func validateServerSetting(key string, value interface{}) error {
	switch key {
	case "maxPlayers":
		n, ok := value.(float64)
		if !ok || math.Trunc(n) != n || n < 1 || n > 16 {
			return fmt.Errorf("最大玩家数必须是 1 到 16 的整数")
		}
	case "autoStartNewDay":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("自动开始新一天必须是布尔值")
		}
	case "newDayWaitSeconds":
		n, ok := value.(float64)
		if !ok || math.Trunc(n) != n || n < 0 || n > 3600 {
			return fmt.Errorf("过夜等待秒数必须是 0 到 3600 的整数")
		}
	default:
		return fmt.Errorf("不支持的服务器设置项")
	}
	return nil
}
