package service

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestParsePercent(t *testing.T) {
	valid := map[string]float64{
		"0.15%": 0.15,
		"3.25%": 3.25,
		"0%":    0,
		"100%":  100,
		"12%":   12,
	}
	for input, want := range valid {
		got, err := parsePercent(input)
		if err != nil {
			t.Errorf("parsePercent(%q) failed: %v", input, err)
			continue
		}
		if got != want {
			t.Errorf("parsePercent(%q) = %v, want %v", input, got, want)
		}
	}

	invalid := []string{"3.25", "abc%", "-1%", "", "%"}
	for _, input := range invalid {
		if _, err := parsePercent(input); err == nil {
			t.Errorf("parsePercent(%q) accepted invalid input", input)
		}
	}
}

func TestParseDockerSize(t *testing.T) {
	valid := map[string]uint64{
		"0B":       0,
		"512B":     512,
		"180MiB":   180 * 1024 * 1024,
		"2GiB":     2 * 1024 * 1024 * 1024,
		"1TiB":     1024 * 1024 * 1024 * 1024,
		"656kB":    656000,
		"1.5GB":    1500000000,
		"1,024MiB": 1024 * 1024 * 1024, // 千分位逗号被容错处理
	}
	for input, want := range valid {
		got, err := parseDockerSize(input)
		if err != nil {
			t.Errorf("parseDockerSize(%q) failed: %v", input, err)
			continue
		}
		if got != want {
			t.Errorf("parseDockerSize(%q) = %d, want %d", input, got, want)
		}
	}

	invalid := []string{"", "MiB", "12", "12XB", "-5B", "1.2.3MiB"}
	for _, input := range invalid {
		if _, err := parseDockerSize(input); err == nil {
			t.Errorf("parseDockerSize(%q) accepted invalid input", input)
		}
	}
}

func TestParseStatsLine(t *testing.T) {
	// docker stats --no-stream --format "{{json .}}" 的真实单行输出
	line := `{"BlockIO":"12.3MB / 0B","CPUPerc":"3.25%","Container":"abc123","ID":"abc123def456","MemPerc":"9.5%","MemUsage":"180MiB / 2GiB","Name":"docker-server-1","NetIO":"1.1MB / 656kB","PIDs":"42"}`

	metrics, err := parseStatsLine(line)
	if err != nil {
		t.Fatalf("parseStatsLine failed: %v", err)
	}
	if metrics.CPUPercent != 3.25 {
		t.Errorf("CPUPercent = %v, want 3.25", metrics.CPUPercent)
	}
	if metrics.MemoryPercent != 9.5 {
		t.Errorf("MemoryPercent = %v, want 9.5", metrics.MemoryPercent)
	}
	if metrics.MemoryBytes != 180*1024*1024 {
		t.Errorf("MemoryBytes = %d, want %d", metrics.MemoryBytes, 180*1024*1024)
	}
	if metrics.MemoryLimitBytes != 2*1024*1024*1024 {
		t.Errorf("MemoryLimitBytes = %d, want %d", metrics.MemoryLimitBytes, 2*1024*1024*1024)
	}
	if metrics.NetReadBytes != 1100000 {
		t.Errorf("NetReadBytes = %d, want 1100000", metrics.NetReadBytes)
	}
	if metrics.NetWriteBytes != 656000 {
		t.Errorf("NetWriteBytes = %d, want 656000", metrics.NetWriteBytes)
	}

	if _, err := parseStatsLine(`{"CPUPerc":"broken"}`); err == nil {
		t.Error("parseStatsLine accepted output with missing fields")
	}
}

func TestComposeStatsReturnsNilWhenContainerStopped(t *testing.T) {
	installFakeDocker(t, `exit 0`)

	runner := NewComposeRunner("compose.yml", t.TempDir())
	metrics, err := runner.Stats("server")
	if err != nil {
		t.Fatalf("Stats failed for stopped container: %v", err)
	}
	if metrics != nil {
		t.Fatalf("expected nil metrics for stopped container, got %+v", metrics)
	}
}

func TestComposeStatsParsesRunningContainer(t *testing.T) {
	installFakeDocker(t, `
if [ "$1" = "compose" ] && [ "$2" = "ps" ]; then echo "abc123"; exit 0; fi
if [ "$1" = "stats" ]; then
  echo '{"BlockIO":"12.3MB / 0B","CPUPerc":"3.25%","Container":"abc123","ID":"abc123","MemPerc":"9.5%","MemUsage":"180MiB / 2GiB","Name":"docker-server-1","NetIO":"1.1MB / 656kB","PIDs":"42"}'
  exit 0
fi
exit 0`)

	runner := NewComposeRunner("compose.yml", t.TempDir())
	metrics, err := runner.Stats("server")
	if err != nil {
		t.Fatalf("Stats failed: %v", err)
	}
	if metrics == nil {
		t.Fatal("expected metrics for running container, got nil")
	}
	if metrics.CPUPercent != 3.25 || metrics.MemoryBytes != 180*1024*1024 {
		t.Fatalf("unexpected metrics: %+v", metrics)
	}
}

func TestStreamLogsFollowsAndStopsOnCancel(t *testing.T) {
	// exec sleep 保证取消时被杀的就是持有 stdout 管道的进程本身
	installFakeDocker(t, `
if [ "$2" = "logs" ]; then
  echo "docker-server-1  | line one"
  echo "docker-server-1  | line two"
  exec sleep 30
fi
exit 0`)

	runner := NewComposeRunner("compose.yml", t.TempDir())
	ctx, cancel := context.WithCancel(context.Background())

	var mu sync.Mutex
	var got []string
	done := make(chan error, 1)
	go func() {
		done <- runner.StreamLogs(ctx, "server", 200, func(line string) {
			mu.Lock()
			got = append(got, line)
			mu.Unlock()
		})
	}()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		count := len(got)
		mu.Unlock()
		if count >= 2 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("StreamLogs returned error after cancel: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("StreamLogs did not return after context cancel")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 || !strings.Contains(got[0], "line one") || !strings.Contains(got[1], "line two") {
		t.Fatalf("unexpected streamed lines: %v", got)
	}
}

func TestGetMetricsFallsBackToPanelOnly(t *testing.T) {
	installFakeDocker(t, `exit 0`)

	svc := NewServerService("compose.yml", t.TempDir(), "", "", "24642", "5800")
	metrics, err := svc.GetMetrics()
	if err != nil {
		t.Fatalf("GetMetrics failed: %v", err)
	}
	if metrics.Container != nil {
		t.Fatalf("expected nil container metrics, got %+v", metrics.Container)
	}
	if metrics.Panel.Goroutines <= 0 || metrics.Panel.UptimeSeconds < 0 {
		t.Fatalf("unexpected panel metrics: %+v", metrics.Panel)
	}
}
