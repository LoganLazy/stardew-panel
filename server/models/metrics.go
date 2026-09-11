package models

// Metrics 性能监控数据：游戏容器资源占用 + 面板进程自身运行时指标。
type Metrics struct {
	Container *ContainerMetrics `json:"container"` // 容器未运行时为 null
	Panel     PanelMetrics      `json:"panel"`
}

// ContainerMetrics 从 docker stats 解析出的容器资源占用。
type ContainerMetrics struct {
	CPUPercent       float64 `json:"cpu_percent"`
	MemoryBytes      uint64  `json:"memory_bytes"`
	MemoryLimitBytes uint64  `json:"memory_limit_bytes"`
	MemoryPercent    float64 `json:"memory_percent"`
	NetReadBytes     uint64  `json:"net_read_bytes"`
	NetWriteBytes    uint64  `json:"net_write_bytes"`
}

// PanelMetrics 面板进程自身的 Go 运行时指标。
type PanelMetrics struct {
	AllocBytes    uint64 `json:"alloc_bytes"`     // 堆内存占用
	SysBytes      uint64 `json:"sys_bytes"`       // 从系统申请的总内存
	Goroutines    int    `json:"goroutines"`      // 协程数
	UptimeSeconds int64  `json:"uptime_seconds"`  // 面板运行时长
}
