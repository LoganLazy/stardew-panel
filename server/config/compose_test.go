package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDockerComposeStructure(t *testing.T) {
	path := filepath.Join("..", "..", "docker", "docker-compose.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var compose struct {
		Name     string                            `yaml:"name"`
		Services map[string]map[string]interface{} `yaml:"services"`
		Volumes  map[string]interface{}            `yaml:"volumes"`
	}
	if err := yaml.Unmarshal(data, &compose); err != nil {
		t.Fatalf("invalid compose YAML: %v", err)
	}
	if compose.Name != "stardew" {
		t.Fatalf("unexpected compose project name: %q", compose.Name)
	}
	for _, service := range []string{"stardew-panel", "server", "steam-auth"} {
		if _, ok := compose.Services[service]; !ok {
			t.Errorf("missing compose service %q", service)
		}
	}
	for _, volume := range []string{"panel-data", "game-data", "saves", "steam-session", "server-settings"} {
		if _, ok := compose.Volumes[volume]; !ok {
			t.Errorf("missing compose volume %q", volume)
		}
	}
}
