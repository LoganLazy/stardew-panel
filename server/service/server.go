package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"stardew-panel/database"
	"stardew-panel/models"
	"strings"
	"sync"
	"time"
)

// sdvd compose 中的服务名（见 docker/docker-compose.yml）
const (
	gameServiceName      = "server"
	steamAuthServiceName = "steam-auth"
)

// ServerService 服务器管理服务。
// 改造后不再自己 exec 游戏进程，而是编排 sdvd/server 容器：
//   - 起停通过 docker compose 命令（compose 字段）
//   - 状态/玩家通过 sdvd REST API（gameAPI 字段）
type ServerService struct {
	compose  *ComposeRunner
	gameAPI  *GameAPIClient
	gamePort string
	vncPort  string

	operationMu sync.RWMutex
	operation   string
	lastError   string
	lifecycleMu sync.Mutex
}

// NewServerService 创建服务器服务。
func NewServerService(composeFile, composeDir, gameAPIURL, gameAPIKey, gamePort, vncPort string) *ServerService {
	return &ServerService{
		compose:  NewComposeRunner(composeFile, composeDir),
		gameAPI:  NewGameAPIClient(gameAPIURL, gameAPIKey),
		gamePort: gamePort,
		vncPort:  vncPort,
	}
}

// GetStatus 获取服务器状态。
// 优先用 compose ps 判断容器是否在跑；容器在跑时再调 sdvd REST API 拿真实游戏数据。
// 数据映射回现有 ServerStatus 结构，保证前端零改动。
func (s *ServerService) GetStatus() (*models.ServerStatus, error) {
	operation, lastError := s.operationState()

	// 1. 容器是否在运行
	running, err := s.compose.PSRunning(gameServiceName)
	if err != nil {
		status := s.baseStatus("error")
		status.Error = err.Error()
		return status, nil
	}

	if !running {
		if operation != "" {
			return s.baseStatus(operation), nil
		}
		if lastError != "" {
			status := s.baseStatus("error")
			status.Error = lastError
			return status, nil
		}
		_ = s.updateStatus("stopped", 0, nil)
		return s.baseStatus("stopped"), nil
	}

	// 2. 容器在跑，尝试拿游戏 API 真实数据
	status := s.baseStatus("running")

	gs, gerr := s.gameAPI.GetStatus()
	if gerr == nil {
		// 游戏 API 就绪
		if !gs.IsReady {
			// 容器在跑但游戏还没 ready（加载存档/过天中）
			status.Status = "starting"
		}
		status.OnlinePlayers = gs.PlayerCount
		if gs.ServerVersion != "" {
			status.Version = gs.ServerVersion
		}
		status.FarmName = gs.FarmName
		status.MaxPlayers = gs.MaxPlayers
	} else {
		// 容器在跑但 API 还没起来，视为启动中
		status.Status = "starting"
	}

	// 3. 玩家名单（API 就绪时才有）
	if players, perr := s.gameAPI.GetPlayers(); perr == nil {
		for _, p := range players {
			if p.IsOnline {
				status.Players = append(status.Players, p.Name)
			}
		}
		status.OnlinePlayers = len(status.Players)
	}
	if operation == "stopping" || operation == "restarting" {
		status.Status = operation
	}

	// 4. 计算 uptime（用数据库记录的启动时间）
	status.Uptime = s.uptimeFromDB()

	// 5. 更新数据库缓存
	if status.Status == "running" {
		now := time.Now()
		// 仅在之前不是 running 时刷新 started_at
		if prev, _ := s.statusFromDB(); prev == nil || prev.Status != "running" {
			s.updateStatus("running", 0, &now)
		}
	}

	return status, nil
}

// Start 启动服务器。
// 只做快速前置检查后立即返回，实际的 `compose up`（首次要下载几 G 游戏文件，
// 可能耗时十分钟）放到后台 goroutine 执行，避免阻塞 HTTP 请求导致浏览器超时。
// 前端置为 starting 后轮询 /status，由 GetStatus 依据 compose ps + 游戏 API
// 反映真实进度（starting → running）。
func (s *ServerService) Start() error {
	if err := s.beginOperation("starting"); err != nil {
		return err
	}

	s.lifecycleMu.Lock()
	if err := s.compose.Available(); err != nil {
		s.lifecycleMu.Unlock()
		s.finishOperation(err)
		return fmt.Errorf("docker compose 不可用，请确认宿主已安装 Docker 且面板有权限: %w", err)
	}

	// 检查是否已在运行
	running, err := s.compose.PSRunning(gameServiceName)
	if err != nil {
		s.lifecycleMu.Unlock()
		s.finishOperation(err)
		return fmt.Errorf("检查服务器状态失败: %w", err)
	}
	if running {
		s.lifecycleMu.Unlock()
		s.finishOperation(nil)
		return fmt.Errorf("服务器已在运行")
	}

	if err := s.updateStatus("starting", 0, nil); err != nil {
		s.lifecycleMu.Unlock()
		s.finishOperation(err)
		return err
	}

	// 后台拉起容器。steam-auth 是 server 的依赖，compose 会按 depends_on 处理，
	// 这里显式带上以确保拉起。
	go func() {
		defer s.lifecycleMu.Unlock()
		if _, err := s.compose.Up(steamAuthServiceName, gameServiceName); err != nil {
			log.Printf("启动游戏服务器失败: %v", err)
			_ = s.updateStatus("stopped", 0, nil)
			s.finishOperation(err)
			return
		}
		now := time.Now()
		_ = s.updateStatus("running", 0, &now)
		s.finishOperation(nil)
	}()

	return nil
}

// Stop 停止服务器：docker compose stop server。
// 只停游戏容器，steam-auth 保留（下次启动更快，且不影响已下载的游戏文件）。
func (s *ServerService) Stop() error {
	_, err := s.StopThen(nil)
	return err
}

// StopThen 停止服务器，并在仍持有生命周期锁时执行停止后的操作。
// 返回的第一个错误只表示停止后操作失败；服务器本身此时已经停止。
func (s *ServerService) StopThen(afterStop func() error) (warning error, retErr error) {
	if err := s.beginOperation("stopping"); err != nil {
		return nil, err
	}
	defer func() {
		if operation, _ := s.operationState(); operation == "stopping" {
			s.finishOperation(retErr)
		}
	}()

	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	if err := s.compose.Available(); err != nil {
		return nil, fmt.Errorf("docker compose 不可用: %w", err)
	}
	running, err := s.compose.PSRunning(gameServiceName)
	if err != nil {
		return nil, fmt.Errorf("检查服务器状态失败: %w", err)
	}
	if !running {
		return nil, fmt.Errorf("服务器未运行")
	}

	if err := s.updateStatus("stopping", 0, nil); err != nil {
		return nil, err
	}

	if out, err := s.compose.Stop(gameServiceName); err != nil {
		return nil, fmt.Errorf("停止容器失败: %w\n%s", err, out)
	}

	if err := s.updateStatus("stopped", 0, nil); err != nil {
		return nil, err
	}
	if afterStop != nil {
		warning = afterStop()
	}
	return warning, nil
}

// Restart 原子地执行停止、停止后操作和重新启动，避免其他请求插入。
func (s *ServerService) Restart(afterStop ...func() error) (retErr error) {
	if err := s.beginOperation("restarting"); err != nil {
		return err
	}
	handedOff := false
	defer func() {
		if !handedOff {
			s.lifecycleMu.Unlock()
			s.finishOperation(retErr)
		}
	}()

	s.lifecycleMu.Lock()
	if err := s.compose.Available(); err != nil {
		return fmt.Errorf("docker compose 不可用: %w", err)
	}
	running, err := s.compose.PSRunning(gameServiceName)
	if err != nil {
		return fmt.Errorf("检查服务器状态失败: %w", err)
	}
	if !running {
		return fmt.Errorf("服务器未运行")
	}
	if err := s.updateStatus("restarting", 0, nil); err != nil {
		return err
	}
	if out, err := s.compose.Stop(gameServiceName); err != nil {
		return fmt.Errorf("停止容器失败: %w\n%s", err, out)
	}
	if err := s.updateStatus("stopped", 0, nil); err != nil {
		return err
	}
	if len(afterStop) > 0 && afterStop[0] != nil {
		if err := afterStop[0](); err != nil {
			return fmt.Errorf("重启前自动备份失败: %w", err)
		}
	}
	if err := s.updateStatus("restarting", 0, nil); err != nil {
		return err
	}

	handedOff = true
	go func() {
		defer s.lifecycleMu.Unlock()
		if _, err := s.compose.Up(steamAuthServiceName, gameServiceName); err != nil {
			log.Printf("重启游戏服务器失败: %v", err)
			_ = s.updateStatus("stopped", 0, nil)
			s.finishOperation(err)
			return
		}
		now := time.Now()
		_ = s.updateStatus("running", 0, &now)
		s.finishOperation(nil)
	}()

	return nil
}

// GetLogs 返回游戏容器最近 lines 行日志（供 /logs 接口）。
// 改造后日志来自 sdvd 容器 stdout，而非本地文件。
func (s *ServerService) GetLogs(lines int) ([]string, error) {
	if lines <= 0 {
		lines = 100
	}
	out, err := s.compose.Logs(gameServiceName, lines)
	if err != nil {
		// 容器不存在/docker 不可用时不报错，返回空列表，避免前端日志页崩
		return []string{}, nil
	}
	result := []string{}
	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}

// StreamLogs 持续跟踪游戏容器日志，供 SSE 接口推送。
// ctx 取消（客户端断开）时正常结束。
func (s *ServerService) StreamLogs(ctx context.Context, tail int, onLine func(line string)) error {
	if tail <= 0 {
		tail = 200
	}
	return s.compose.StreamLogs(ctx, gameServiceName, tail, onLine)
}

// GetMetrics 返回性能监控数据：游戏容器资源占用 + 面板自身运行时指标。
// 容器未运行或 docker 不可用时 Container 为 nil，不视为错误（前端显示占位提示）。
func (s *ServerService) GetMetrics() (*models.Metrics, error) {
	metrics := &models.Metrics{
		Container: nil,
		Panel:     panelMetrics(),
	}

	container, err := s.compose.Stats(gameServiceName)
	if err != nil {
		// docker stats 失败不影响面板自身指标，保留 Container=nil 返回
		log.Printf("获取容器性能指标失败: %v", err)
		return metrics, nil
	}
	metrics.Container = container
	return metrics, nil
}

// GetInviteCode 返回联机邀请码（供前端显示）。容器未就绪时返回空。
func (s *ServerService) GetInviteCode() (string, error) {
	gs, err := s.gameAPI.GetStatus()
	if err != nil {
		return "", err
	}
	return gs.InviteCode(), nil
}

// GetSettings 透传 sdvd 的服务器设置（供前端设置页读取）。
func (s *ServerService) GetSettings() (map[string]interface{}, error) {
	return s.gameAPI.GetSettings()
}

// UpdateSetting 透传更新单个设置项到 sdvd。
func (s *ServerService) UpdateSetting(key string, value interface{}) error {
	return s.gameAPI.UpdateSetting(key, value)
}

// --- 数据库缓存辅助 ---

// statusFromDB 从数据库读取缓存的服务器状态（docker 不可用时的兜底）。
func (s *ServerService) statusFromDB() (*models.ServerStatus, error) {
	row := database.DB.QueryRow(`SELECT status, started_at FROM server_state WHERE id = 1`)
	var status string
	var startedAt sql.NullTime
	if err := row.Scan(&status, &startedAt); err != nil {
		return nil, err
	}

	uptime := int64(0)
	if status == "running" && startedAt.Valid {
		uptime = int64(time.Since(startedAt.Time).Seconds())
	}

	return &models.ServerStatus{
		Status:        status,
		OnlinePlayers: 0,
		Version:       s.cachedVersion(),
		Uptime:        uptime,
		Players:       []string{},
		GamePort:      s.gamePort,
		VNCPort:       s.vncPort,
	}, nil
}

// uptimeFromDB 从数据库 started_at 计算运行时长（秒）。
func (s *ServerService) uptimeFromDB() int64 {
	row := database.DB.QueryRow(`SELECT started_at FROM server_state WHERE id = 1`)
	var startedAt sql.NullTime
	if err := row.Scan(&startedAt); err != nil || !startedAt.Valid {
		return 0
	}
	return int64(time.Since(startedAt.Time).Seconds())
}

// cachedVersion 返回缓存的游戏版本（暂用固定值，后续可从 API 持久化）。
func (s *ServerService) cachedVersion() string {
	return "未知"
}

// updateStatus 更新数据库中的服务器状态缓存。
// pid 字段在容器模式下已无意义，恒为 0，仅保留列兼容。
func (s *ServerService) updateStatus(status string, pid int, startedAt *time.Time) error {
	if startedAt == nil {
		_, err := database.DB.Exec(`
			UPDATE server_state
			SET status = ?, pid = ?, started_at = NULL, updated_at = CURRENT_TIMESTAMP
			WHERE id = 1
		`, status, pid)
		return err
	}

	_, err := database.DB.Exec(`
		UPDATE server_state
		SET status = ?, pid = ?, started_at = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = 1
	`, status, pid, startedAt)
	return err
}

func (s *ServerService) baseStatus(status string) *models.ServerStatus {
	return &models.ServerStatus{
		Status:   status,
		Version:  s.cachedVersion(),
		Players:  []string{},
		GamePort: s.gamePort,
		VNCPort:  s.vncPort,
	}
}

func (s *ServerService) beginOperation(operation string) error {
	s.operationMu.Lock()
	defer s.operationMu.Unlock()
	if s.operation != "" {
		return fmt.Errorf("服务器正在%s，请稍候", operationLabel(s.operation))
	}
	s.operation = operation
	s.lastError = ""
	return nil
}

func (s *ServerService) finishOperation(err error) {
	s.operationMu.Lock()
	defer s.operationMu.Unlock()
	s.operation = ""
	if err != nil {
		s.lastError = err.Error()
	} else {
		s.lastError = ""
	}
}

func (s *ServerService) operationState() (string, string) {
	s.operationMu.RLock()
	defer s.operationMu.RUnlock()
	return s.operation, s.lastError
}

// AcquireStopped 将存档操作与起停流程串行化，并保证操作期间游戏容器未运行。
func (s *ServerService) AcquireStopped() (func(), error) {
	operation, _ := s.operationState()
	if operation != "" {
		return nil, fmt.Errorf("服务器正在%s，暂不能操作存档", operationLabel(operation))
	}

	s.lifecycleMu.Lock()
	operation, _ = s.operationState()
	if operation != "" {
		s.lifecycleMu.Unlock()
		return nil, fmt.Errorf("服务器正在%s，暂不能操作存档", operationLabel(operation))
	}
	running, err := s.compose.PSRunning(gameServiceName)
	if err != nil {
		s.lifecycleMu.Unlock()
		return nil, fmt.Errorf("无法确认服务器已停止: %w", err)
	}
	if running {
		s.lifecycleMu.Unlock()
		return nil, fmt.Errorf("请先停止游戏服务器，再操作存档")
	}
	return s.lifecycleMu.Unlock, nil
}

func operationLabel(operation string) string {
	switch operation {
	case "starting":
		return "启动"
	case "stopping":
		return "停止"
	case "restarting":
		return "重启"
	default:
		return "执行操作"
	}
}
