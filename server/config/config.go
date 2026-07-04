package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Game     GameConfig     `yaml:"game"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type GameConfig struct {
	ServerPath string `yaml:"server_path"`
	ModsPath   string `yaml:"mods_path"`
	SavesPath  string `yaml:"saves_path"`
}

// Load 加载配置文件
func Load() (*Config, error) {
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		// 如果配置文件不存在，返回默认配置
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func getConfigPath() string {
	if path := os.Getenv("STARDEW_PANEL_CONFIG"); path != "" {
		return path
	}
	return "./config.yaml"
}

func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: "8080",
		},
		Database: DatabaseConfig{
			Path: "./data/stardew-panel.db",
		},
		Game: GameConfig{
			ServerPath: "./game/server",
			ModsPath:   "./game/mods",
			SavesPath:  "./game/saves",
		},
	}
}
