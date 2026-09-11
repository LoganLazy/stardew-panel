package service

import (
	"archive/zip"
	"os"
	"path/filepath"
	"stardew-panel/database"
	"testing"
)

func writeTestZip(t *testing.T, path string, files map[string]string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	for name, content := range files {
		entry, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestScanModsRemovesMissingDatabaseRows(t *testing.T) {
	root := t.TempDir()
	if err := database.Init(filepath.Join(root, "panel.db")); err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	mods := filepath.Join(root, "Mods")
	modDir := filepath.Join(mods, "Example")
	if err := os.MkdirAll(modDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(modDir, "manifest.json"), []byte(`{"Name":"Example","UniqueID":"test.example","Version":"1.0.0"}`), 0644); err != nil {
		t.Fatal(err)
	}

	service := NewModService(mods)
	if err := service.ScanMods(); err != nil {
		t.Fatal(err)
	}
	listed, err := service.ListMods()
	if err != nil || len(listed) != 1 {
		t.Fatalf("mod was not registered: %#v, %v", listed, err)
	}
	if err := os.RemoveAll(modDir); err != nil {
		t.Fatal(err)
	}
	if err := service.ScanMods(); err != nil {
		t.Fatal(err)
	}
	listed, err = service.ListMods()
	if err != nil || len(listed) != 0 {
		t.Fatalf("stale mod row was not removed: %#v, %v", listed, err)
	}
}

func TestUnzipModRejectsPathTraversal(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "bad.zip")
	writeTestZip(t, archive, map[string]string{"../escape.txt": "bad"})

	service := NewModService(filepath.Join(root, "mods"))
	if err := service.unzipMod(archive, filepath.Join(root, "mods", "bad")); err == nil {
		t.Fatal("expected path traversal archive to be rejected")
	}
	if _, err := os.Stat(filepath.Join(root, "escape.txt")); !os.IsNotExist(err) {
		t.Fatal("archive wrote outside destination")
	}
}

func TestUnzipModAndParseNestedManifest(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "example.zip")
	writeTestZip(t, archive, map[string]string{
		"Example/manifest.json": `{"Name":"Example","UniqueID":"test.example","Version":"1.0.0"}`,
		"Example/mod.dll":       "content",
	})

	service := NewModService(filepath.Join(root, "mods"))
	dest := filepath.Join(root, "mods", "example")
	if err := service.unzipMod(archive, dest); err != nil {
		t.Fatalf("unzip failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "manifest.json")); err != nil {
		t.Fatalf("single top-level directory was not flattened: %v", err)
	}
	mod, err := service.parseModManifest(dest)
	if err != nil {
		t.Fatalf("parse manifest failed: %v", err)
	}
	if mod.Name != "Example" || mod.UniqueID != "test.example" {
		t.Fatalf("unexpected manifest: %#v", mod)
	}
}
