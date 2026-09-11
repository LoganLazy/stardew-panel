package service

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"stardew-panel/database"
	"strings"
	"testing"
	"time"
)

func installFakeDocker(t *testing.T, body string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("fake docker shell script requires a Unix test environment")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "docker")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body+"\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func TestComposeAvailableChecksDockerDaemon(t *testing.T) {
	installFakeDocker(t, `
if [ "$1" = "info" ]; then
  echo "daemon unavailable" >&2
  exit 1
fi
exit 0`)

	runner := NewComposeRunner("compose.yml", t.TempDir())
	if err := runner.Available(); err == nil || !strings.Contains(err.Error(), "daemon unavailable") {
		t.Fatalf("expected daemon error, got %v", err)
	}
}

func TestStartReportsBackgroundComposeFailure(t *testing.T) {
	installFakeDocker(t, `
if [ "$1" = "info" ]; then exit 0; fi
case " $* " in
  *" up -d "*) echo "pull failed" >&2; exit 1 ;;
esac
exit 0`)
	if err := database.Init(filepath.Join(t.TempDir(), "panel.db")); err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	service := NewServerService("compose.yml", t.TempDir(), "", "", "24642", "5800")
	if err := service.Start(); err != nil {
		t.Fatalf("start request failed before background operation: %v", err)
	}

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		status, err := service.GetStatus()
		if err != nil {
			t.Fatal(err)
		}
		if status.Status == "error" {
			if !strings.Contains(status.Error, "pull failed") {
				t.Fatalf("unexpected operation error: %q", status.Error)
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("background compose failure was not exposed")
}

func TestStopThenKeepsLifecycleOperationActiveDuringCallback(t *testing.T) {
	installFakeDocker(t, `
if [ "$1" = "info" ]; then exit 0; fi
case " $* " in
  *" ps -q --status running server "*) echo "container-id" ;;
esac
exit 0`)
	if err := database.Init(filepath.Join(t.TempDir(), "panel.db")); err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	server := NewServerService("compose.yml", t.TempDir(), "", "", "24642", "5800")
	var competingErr error
	warning, err := server.StopThen(func() error {
		competingErr = server.Start()
		return nil
	})
	if err != nil {
		t.Fatalf("stop failed: %v", err)
	}
	if warning != nil {
		t.Fatalf("unexpected callback warning: %v", warning)
	}
	if competingErr == nil || !strings.Contains(competingErr.Error(), "正在停止") {
		t.Fatalf("competing start was not rejected during callback: %v", competingErr)
	}
	if operation, _ := server.operationState(); operation != "" {
		t.Fatalf("operation was not released after stop: %q", operation)
	}
}

func TestRestartAbortsWhenBackupFails(t *testing.T) {
	installFakeDocker(t, `
if [ "$1" = "info" ]; then exit 0; fi
case " $* " in
  *" ps -q --status running server "*) echo "container-id" ;;
esac
exit 0`)
	if err := database.Init(filepath.Join(t.TempDir(), "panel.db")); err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	server := NewServerService("compose.yml", t.TempDir(), "", "", "24642", "5800")
	err := server.Restart(func() error { return fmt.Errorf("backup failed") })
	if err == nil || !strings.Contains(err.Error(), "backup failed") {
		t.Fatalf("expected synchronous backup failure, got %v", err)
	}
	operation, lastError := server.operationState()
	if operation != "" || !strings.Contains(lastError, "backup failed") {
		t.Fatalf("unexpected operation state after failure: operation=%q error=%q", operation, lastError)
	}
}
