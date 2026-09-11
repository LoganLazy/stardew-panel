package service

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"stardew-panel/models"
	"strings"
	"time"
)

// ComposeRunner 封装对 docker compose 命令的调用。
// 用 exec.Command 数组传参，不经过 shell，避免命令注入。
type ComposeRunner struct {
	composeFile string // -f 指定的 compose 文件路径；空则用 workDir 下默认文件
	workDir     string // 执行 compose 命令的工作目录
}

// NewComposeRunner 创建 compose 命令执行器。
func NewComposeRunner(composeFile, workDir string) *ComposeRunner {
	return &ComposeRunner{
		composeFile: composeFile,
		workDir:     workDir,
	}
}

// baseArgs 构造 compose 子命令的公共前缀参数。
// 使用 "docker compose"（v2 插件形式），而非老的 "docker-compose"。
func (r *ComposeRunner) baseArgs() []string {
	args := []string{"compose"}
	if r.composeFile != "" {
		args = append(args, "-f", r.composeFile)
	}
	return args
}

// run 执行一条 compose 命令，返回合并后的输出。
// timeout 控制单条命令最长执行时间，防止卡死（如网络拉镜像）。
func (r *ComposeRunner) run(timeout time.Duration, subArgs ...string) (string, error) {
	args := append(r.baseArgs(), subArgs...)
	return r.runDocker(timeout, args...)
}

// runDocker 执行 docker CLI 命令。与 run 分开是为了检查 daemon，而不只是
// 检查本地是否安装了 compose 插件。
func (r *ComposeRunner) runDocker(timeout time.Duration, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", args...)
	if r.workDir != "" {
		cmd.Dir = r.workDir
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	output := out.String()
	if ctx.Err() == context.DeadlineExceeded {
		return output, fmt.Errorf("docker %s 执行超时", strings.Join(args, " "))
	}
	if err != nil {
		return output, fmt.Errorf("docker %s 失败: %w\n%s", strings.Join(args, " "), err, output)
	}
	return output, nil
}

// Up 启动指定服务（-d 后台）。首次可能拉镜像，给较长超时。
func (r *ComposeRunner) Up(services ...string) (string, error) {
	args := []string{"up", "-d"}
	args = append(args, services...)
	return r.run(10*time.Minute, args...)
}

// Stop 停止指定服务（保留容器，可再次 start）。
func (r *ComposeRunner) Stop(services ...string) (string, error) {
	args := []string{"stop"}
	args = append(args, services...)
	return r.run(2*time.Minute, args...)
}

// Restart 重启指定服务。
func (r *ComposeRunner) Restart(services ...string) (string, error) {
	args := []string{"restart"}
	args = append(args, services...)
	return r.run(3*time.Minute, args...)
}

// Down 停止并移除所有编排的容器（不删 volume）。
func (r *ComposeRunner) Down() (string, error) {
	return r.run(2*time.Minute, "down")
}

// PSRunning 返回指定服务当前是否有正在运行的容器。
// 用 `compose ps -q --status running <service>`：有输出即在运行。
func (r *ComposeRunner) PSRunning(service string) (bool, error) {
	out, err := r.run(30*time.Second, "ps", "-q", "--status", "running", service)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

// Logs 返回指定服务最近 tail 行日志。
func (r *ComposeRunner) Logs(service string, tail int) (string, error) {
	return r.run(30*time.Second, "logs", "--tail", fmt.Sprintf("%d", tail), "--no-color", service)
}

// StreamLogs 持续跟踪指定服务日志：先输出最近 tail 行，随后不断推送新行。
// 每读到一行调用一次 onLine；ctx 取消时结束并返回 nil（正常断开）。
func (r *ComposeRunner) StreamLogs(ctx context.Context, service string, tail int, onLine func(line string)) error {
	args := append(r.baseArgs(), "logs", "--follow", "--no-color", "--tail", fmt.Sprintf("%d", tail), service)

	cmd := exec.CommandContext(ctx, "docker", args...)
	if r.workDir != "" {
		cmd.Dir = r.workDir
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建日志管道失败: %w", err)
	}
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("docker logs 跟踪启动失败: %w", err)
	}

	// SMAPI 的异常堆栈可能很长，把单行上限放宽到 1MB
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		onLine(scanner.Text())
	}
	waitErr := cmd.Wait()

	// 客户端断开导致的取消不算错误
	if ctx.Err() != nil {
		return nil
	}
	if waitErr != nil {
		return fmt.Errorf("docker logs 跟踪失败: %w\n%s", waitErr, errBuf.String())
	}
	return scanner.Err()
}

// Stats 返回指定服务容器的资源占用。
// 容器未运行时返回 (nil, nil)，由调用方区分"没数据"和"出错"。
func (r *ComposeRunner) Stats(service string) (*models.ContainerMetrics, error) {
	// compose ps 拿容器 ID（有输出即在运行）
	idOut, err := r.run(30*time.Second, "ps", "-q", "--status", "running", service)
	if err != nil {
		return nil, err
	}
	containerID := strings.TrimSpace(idOut)
	if containerID == "" {
		return nil, nil
	}

	out, err := r.runDocker(30*time.Second, "stats", "--no-stream", "--format", "{{json .}}", containerID)
	if err != nil {
		return nil, err
	}
	line := strings.TrimSpace(out)
	if line == "" {
		return nil, fmt.Errorf("docker stats 无输出")
	}
	// 只取第一行（指定了单个容器，正常也只有一行）
	if idx := strings.IndexByte(line, '\n'); idx >= 0 {
		line = line[:idx]
	}
	return parseStatsLine(line)
}

// Available 检查 daemon、socket、compose 插件和当前编排文件是否都可用。
func (r *ComposeRunner) Available() error {
	if _, err := r.runDocker(15*time.Second, "info", "--format", "{{.ServerVersion}}"); err != nil {
		return err
	}
	_, err := r.run(15*time.Second, "config", "--quiet")
	return err
}
