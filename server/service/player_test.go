package service

import (
	"path/filepath"
	"stardew-panel/database"
	"testing"
)

func TestSyncLivePlayersMaintainsHistory(t *testing.T) {
	root := t.TempDir()
	if err := database.Init(filepath.Join(root, "panel.db")); err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	service := NewPlayerService("", "")
	if err := service.syncLivePlayers([]string{"Alice", "Bob"}); err != nil {
		t.Fatal(err)
	}
	if err := service.syncLivePlayers([]string{"Bob", "Carol"}); err != nil {
		t.Fatal(err)
	}

	players, err := service.GetOnlinePlayers()
	if err != nil {
		t.Fatal(err)
	}
	if len(players) != 2 || players[0].ConnectedAt.IsZero() || players[1].ConnectedAt.IsZero() {
		t.Fatalf("unexpected live players: %#v", players)
	}
	stats, err := service.GetPlayerStats()
	if err != nil {
		t.Fatal(err)
	}
	if stats["online_count"] != 2 || stats["total_count"] != 3 {
		t.Fatalf("unexpected stats: %#v", stats)
	}
}
