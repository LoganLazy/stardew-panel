package service

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"stardew-panel/models"
)

// panelStartedAt 面板进程启动时间，用于计算面板自身运行时长。
var panelStartedAt = time.Now()

// statsRawLine 对应 `docker stats --no-stream --format "{{json .}}"` 的单行输出，
// 所有字段都是字符串（docker CLI 模板输出的原样格式）。
type statsRawLine struct {
	CPUPerc  string `json:"CPUPerc"`
	MemPerc  string `json:"MemPerc"`
	MemUsage string `json:"MemUsage"` // "50.2MiB / 3.84GiB"
	NetIO    string `json:"NetIO"`    // "1.2kB / 3.4kB"
}

// parseStatsLine 解析一行 docker stats JSON 输出为容器指标。
func parseStatsLine(line string) (*models.ContainerMetrics, error) {
	var raw statsRawLine
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("解析 docker stats 输出失败: %w", err)
	}

	cpu, err := parsePercent(raw.CPUPerc)
	if err != nil {
		return nil, err
	}
	memPercent, err := parsePercent(raw.MemPerc)
	if err != nil {
		return nil, err
	}
	memUsed, memLimit, err := parseSizePair(raw.MemUsage)
	if err != nil {
		return nil, err
	}
	netRead, netWrite, err := parseSizePair(raw.NetIO)
	if err != nil {
		return nil, err
	}

	return &models.ContainerMetrics{
		CPUPercent:       cpu,
		MemoryBytes:      memUsed,
		MemoryLimitBytes: memLimit,
		MemoryPercent:    memPercent,
		NetReadBytes:     netRead,
		NetWriteBytes:    netWrite,
	}, nil
}

// parsePercent 解析 "0.15%" 形式的百分比字符串。
func parsePercent(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if !strings.HasSuffix(s, "%") {
		return 0, fmt.Errorf("无效的百分比格式: %q", s)
	}
	v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
	if err != nil || v < 0 {
		return 0, fmt.Errorf("无效的百分比格式: %q", s)
	}
	return v, nil
}

// dockerSizeUnits docker stats 输出使用的容量单位。
// 二进制单位（KiB/MiB/GiB/TiB）按 1024 进位，十进制单位（kB/MB/GB）按 1000 进位。
var dockerSizeUnits = map[string]float64{
	"B":   1,
	"KiB": 1024,
	"MiB": 1024 * 1024,
	"GiB": 1024 * 1024 * 1024,
	"TiB": 1024 * 1024 * 1024 * 1024,
	"kB":  1000,
	"MB":  1000 * 1000,
	"GB":  1000 * 1000 * 1000,
	"TB":  1000 * 1000 * 1000 * 1000,
}

// parseDockerSize 解析 "50.2MiB"、"512B" 形式的容量字符串为字节数。
func parseDockerSize(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("空的容量值")
	}

	// 找到数字与单位的分界（单位可能以小数点后的字母开头，如 "1.5GiB"）
	unitStart := strings.IndexAny(s, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	if unitStart <= 0 {
		return 0, fmt.Errorf("无效的容量格式: %q", s)
	}

	number, unit := s[:unitStart], s[unitStart:]
	number = strings.ReplaceAll(number, ",", "") // 容错：去掉可能出现的千分位逗号
	multiplier, ok := dockerSizeUnits[unit]
	if !ok {
		return 0, fmt.Errorf("未知的容量单位: %q", unit)
	}
	v, err := strconv.ParseFloat(number, 64)
	if err != nil || v < 0 {
		return 0, fmt.Errorf("无效的容量格式: %q", s)
	}
	return uint64(v * multiplier), nil
}

// parseSizePair 解析 "已用 / 上限" 形式的字段（MemUsage、NetIO）。
func parseSizePair(s string) (used, limit uint64, err error) {
	parts := strings.Split(strings.TrimSpace(s), "/")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("无效的用量对格式: %q", s)
	}
	if used, err = parseDockerSize(parts[0]); err != nil {
		return 0, 0, err
	}
	if limit, err = parseDockerSize(parts[1]); err != nil {
		return 0, 0, err
	}
	return used, limit, nil
}

// panelMetrics 采集面板进程自身的运行时指标。
func panelMetrics() models.PanelMetrics {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return models.PanelMetrics{
		AllocBytes:    mem.Alloc,
		SysBytes:      mem.Sys,
		Goroutines:    runtime.NumGoroutine(),
		UptimeSeconds: int64(time.Since(panelStartedAt).Seconds()),
	}
}
