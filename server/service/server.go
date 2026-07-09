package service

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"stardew-panel/database"
	"stardew-panel/models"
	"syscall"
	"time"
)

// ServerService 服务器管理服务
type ServerService struct {
	gamePath string
	cmd      *exec.Cmd
}

// NewServerService 创建服务器服务
func NewServerService(gamePath string) *ServerService {
	return &ServerService{
		gamePath: gamePath,
	}
}

// GetStatus 获取服务器状态
func (s *ServerService) GetStatus() (*models.ServerStatus, error) {
	// 从数据库获取状态
	row := database.DB.QueryRow(`
		SELECT status, pid, started_at
		FROM server_state
		WHERE id = 1
	`)

	var status string
	var pid sql.NullInt64
	var startedAt sql.NullTime
	err := row.Scan(&status, &pid, &startedAt)
	if err != nil {
		return nil, err
	}

	// 如果状态是 running，检查进程是否真的在运行
	if status == "running" && pid.Valid {
		if !s.isProcessRunning(int(pid.Int64)) {
			// 进程已死，更新状态
			s.updateStatus("stopped", 0, nil)
			status = "stopped"
		}
	}

	// 计算运行时间
	uptime := int64(0)
	if status == "running" && startedAt.Valid {
		uptime = int64(time.Since(startedAt.Time).Seconds())
	}

	// 获取游戏版本
	version := "1.6.9"
	inst, _ := s.getInstallation()
	if inst != nil {
		version = inst.Version
	}

	// 获取在线玩家数
	onlinePlayers := 0
	playerNames := []string{}

	// 使用 PlayerService 获取在线玩家
	playerService := NewPlayerService()
	playerService.ParseLogForPlayers() // 解析最新日志

	players, err := playerService.GetOnlinePlayers()
	if err == nil {
		onlinePlayers = len(players)
		for _, p := range players {
			playerNames = append(playerNames, p.PlayerName)
		}
	}

	return &models.ServerStatus{
		Status:        status,
		OnlinePlayers: onlinePlayers,
		Version:       version,
		Uptime:        uptime,
		Players:       playerNames,
	}, nil
}

// Start 启动服务器
func (s *ServerService) Start() error {
	// 检查是否已安装
	inst, err := s.getInstallation()
	if err != nil {
		return fmt.Errorf("检查安装失败: %w", err)
	}
	if inst == nil {
		return fmt.Errorf("服务器尚未安装")
	}

	// 检查当前状态
	status, err := s.GetStatus()
	if err != nil {
		return err
	}
	if status.Status == "running" {
		return fmt.Errorf("服务器已在运行")
	}

	// 更新状态为 starting
	if err := s.updateStatus("starting", 0, nil); err != nil {
		return err
	}

	// 确定可执行文件路径
	var execPath string
	if inst.HasSMAPI {
		execPath = s.findExecutable(inst.GamePath, []string{"StardewModdingAPI", "StardewModdingAPI.exe"})
	} else {
		execPath = s.findExecutable(inst.GamePath, []string{"StardewValley", "StardewValley.exe"})
	}

	if execPath == "" {
		s.updateStatus("stopped", 0, nil)
		return fmt.Errorf("未找到游戏可执行文件")
	}

	// 启动服务器进程
	s.cmd = exec.Command(execPath, "--server")
	s.cmd.Dir = inst.GamePath

	// 创建日志文件
	logFile, err := os.OpenFile("./data/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		s.updateStatus("stopped", 0, nil)
		return fmt.Errorf("创建日志文件失败: %w", err)
	}

	s.cmd.Stdout = logFile
	s.cmd.Stderr = logFile

	// 启动进程
	if err := s.cmd.Start(); err != nil {
		logFile.Close()
		s.updateStatus("stopped", 0, nil)
		return fmt.Errorf("启动服务器失败: %w", err)
	}

	// 更新状态为 running
	now := time.Now()
	if err := s.updateStatus("running", s.cmd.Process.Pid, &now); err != nil {
		// 如果状态更新失败，终止进程
		s.cmd.Process.Kill()
		logFile.Close()
		return err
	}

	// 后台监控进程
	go func() {
		s.cmd.Wait()
		logFile.Close()
		s.updateStatus("stopped", 0, nil)
	}()

	return nil
}

// Stop 停止服务器
func (s *ServerService) Stop() error {
	// 检查当前状态
	status, err := s.GetStatus()
	if err != nil {
		return err
	}
	if status.Status != "running" {
		return fmt.Errorf("服务器未运行")
	}

	// 更新状态为 stopping
	if err := s.updateStatus("stopping", 0, nil); err != nil {
		return err
	}

	// 获取 PID
	row := database.DB.QueryRow(`SELECT pid FROM server_state WHERE id = 1`)
	var pid int
	if err := row.Scan(&pid); err != nil {
		return err
	}

	// 尝试优雅关闭
	if s.cmd != nil && s.cmd.Process != nil {
		s.cmd.Process.Signal(syscall.SIGTERM)

		// 等待最多 10 秒
		done := make(chan error, 1)
		go func() {
			done <- s.cmd.Wait()
		}()

		select {
		case <-done:
			// 进程已退出
		case <-time.After(10 * time.Second):
			// 超时，强制终止
			s.cmd.Process.Kill()
		}
	} else {
		// 如果没有保存的 cmd，直接发送信号
		process, err := os.FindProcess(pid)
		if err == nil {
			process.Signal(syscall.SIGTERM)
			time.Sleep(2 * time.Second)
			process.Kill()
		}
	}

	// 更新状态为 stopped
	if err := s.updateStatus("stopped", 0, nil); err != nil {
		return err
	}

	return nil
}

// Restart 重启服务器
func (s *ServerService) Restart() error {
	// 先停止
	status, _ := s.GetStatus()
	if status.Status == "running" {
		if err := s.Stop(); err != nil {
			return err
		}
		// 等待停止完成
		time.Sleep(2 * time.Second)
	}

	// 再启动
	return s.Start()
}

// GetLogs 获取服务器日志
func (s *ServerService) GetLogs(lines int) ([]string, error) {
	logPath := "./data/server.log"
	file, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	// 读取最后 N 行
	var result []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
		if len(result) > lines {
			result = result[1:]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// 辅助方法

func (s *ServerService) getInstallation() (*models.Installation, error) {
	row := database.DB.QueryRow(`
		SELECT id, game_path, install_method, has_smapi, version
		FROM installation
		ORDER BY id DESC LIMIT 1
	`)

	var inst models.Installation
	err := row.Scan(&inst.ID, &inst.GamePath, &inst.InstallMethod, &inst.HasSMAPI, &inst.Version)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &inst, nil
}

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

func (s *ServerService) isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// 发送信号 0 检查进程是否存在
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func (s *ServerService) findExecutable(basePath string, names []string) string {
	for _, name := range names {
		fullPath := basePath + string(os.PathSeparator) + name
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}
	return ""
}
