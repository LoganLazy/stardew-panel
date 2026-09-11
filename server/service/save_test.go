package service

import (
	"os"
	"path/filepath"
	"stardew-panel/database"
	"strings"
	"testing"
)

func TestSafeSavePath(t *testing.T) {
	root := t.TempDir()
	service := NewSaveService(root)

	path, err := service.safeSavePath("Farm_123")
	if err != nil {
		t.Fatalf("valid save rejected: %v", err)
	}
	if path != filepath.Join(root, "Farm_123") {
		t.Fatalf("unexpected save path: %s", path)
	}

	for _, name := range []string{"../outside", `..\\outside`, "/outside", ".", ".."} {
		if _, err := service.safeSavePath(name); err == nil {
			t.Errorf("unsafe save name accepted: %q", name)
		}
	}
}

func TestSaveZipRoundTrip(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	if err := os.MkdirAll(filepath.Join(source, "nested"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "nested", "SaveGameInfo"), []byte("save-data"), 0644); err != nil {
		t.Fatal(err)
	}

	service := NewSaveService(filepath.Join(root, "saves"))
	archive := filepath.Join(root, "save.zip")
	if err := service.zipDirectory(source, archive); err != nil {
		t.Fatalf("zip failed: %v", err)
	}
	target := filepath.Join(root, "target")
	if err := service.unzipDirectory(archive, target); err != nil {
		t.Fatalf("unzip failed: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(target, "nested", "SaveGameInfo"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "save-data" {
		t.Fatalf("unexpected restored data: %q", data)
	}
}

func TestRestoreRegistersSafetyBackup(t *testing.T) {
	root := t.TempDir()
	if err := database.Init(filepath.Join(root, "panel.db")); err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	saves := filepath.Join(root, "Saves")
	backups := filepath.Join(root, "backups")
	farm := filepath.Join(saves, "Farm_1")
	if err := os.MkdirAll(farm, 0755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(farm, "SaveGameInfo")
	if err := os.WriteFile(file, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	service := NewSaveService(saves, backups)
	backup, err := service.CreateBackup("Farm_1")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := service.RestoreBackup(backup.ID); err != nil {
		t.Fatal(err)
	}

	listed, err := service.ListBackups()
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 2 {
		t.Fatalf("expected requested backup plus registered safety backup, got %d", len(listed))
	}
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "old" {
		t.Fatalf("restore did not replace save: %q", data)
	}
}

func TestAutoBackupReturnsIndividualSaveFailures(t *testing.T) {
	root := t.TempDir()
	if err := database.Init(filepath.Join(root, "panel.db")); err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	saves := filepath.Join(root, "Saves")
	farm := filepath.Join(saves, "Farm_1")
	if err := os.MkdirAll(farm, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "outside"), filepath.Join(farm, "unsafe-link")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	service := NewSaveService(saves, filepath.Join(root, "backups"))
	err := service.AutoBackup()
	if err == nil || !strings.Contains(err.Error(), "Farm_1") {
		t.Fatalf("expected per-save backup error, got %v", err)
	}
	backups, listErr := service.ListBackups()
	if listErr != nil {
		t.Fatal(listErr)
	}
	if len(backups) != 0 {
		t.Fatalf("failed archive should not be registered, got %d backups", len(backups))
	}
}
