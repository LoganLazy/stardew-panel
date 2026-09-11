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
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	CORSOrigins    string `yaml:"cors_origins"`
	TrustedProxies string `yaml:"trusted_proxies"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type GameConfig struct {
	ServerPath  string `yaml:"server_path"`
	ModsPath    string `yaml:"mods_path"`
	SavesPath   string `yaml:"saves_path"`
	BackupsPath string `yaml:"backups_path"`

	// 编排相关（改造后新增）
	ComposeFile string `yaml:"compose_file"` // docker compose 文件路径
	ComposeDir  string `yaml:"compose_dir"`  // 执行 compose 命令的工作目录
	GameAPIURL  string `yaml:"game_api_url"` // sdvd 游戏容器 REST API 地址，如 http://server:8080
	GameAPIKey  string `yaml:"game_api_key"` // sdvd API_KEY，转发时用
	GamePort    string `yaml:"game_port"`    // 宿主对外联机端口
	VNCPort     string `yaml:"vnc_port"`     // 宿主对外 VNC 网页端口
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

	applyEnvOverrides(&cfg)
	return &cfg, nil
}

// applyEnvOverrides 用环境变量覆盖配置。
// compose 部署时敏感信息（API_KEY）和地址通过 env 注入，优先级高于 yaml。
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("PANEL_PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("CORS_ORIGINS"); v != "" {
		cfg.Server.CORSOrigins = v
	}
	if v := os.Getenv("TRUSTED_PROXIES"); v != "" {
		cfg.Server.TrustedProxies = v
	}
	if v := os.Getenv("GAME_API_URL"); v != "" {
		cfg.Game.GameAPIURL = v
	}
	if v := os.Getenv("API_KEY"); v != "" {
		cfg.Game.GameAPIKey = v
	}
	if v := os.Getenv("GAME_PORT"); v != "" {
		cfg.Game.GamePort = v
	}
	if v := os.Getenv("VNC_PORT"); v != "" {
		cfg.Game.VNCPort = v
	}
	if v := os.Getenv("COMPOSE_FILE"); v != "" {
		cfg.Game.ComposeFile = v
	}
	if v := os.Getenv("COMPOSE_DIR"); v != "" {
		cfg.Game.ComposeDir = v
	}
}

func getConfigPath() string {
	if path := os.Getenv("STARDEW_PANEL_CONFIG"); path != "" {
		return path
	}
	return "./config.yaml"
}

func defaultConfig() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: "9090", // 改：避开 sdvd 游戏容器的 8080
		},
		Database: DatabaseConfig{
			Path: "./data/stardew-panel.db",
		},
		Game: GameConfig{
			ServerPath:  "./game/server",
			ModsPath:    "./game/server/Mods",
			SavesPath:   "./game/saves/Saves",
			BackupsPath: "./game/saves/backups",
			// 编排默认值：compose 文件在项目 docker/ 下，游戏 API 走容器网络内的 server:8080
			ComposeFile: "docker/docker-compose.yml",
			ComposeDir:  ".",
			GameAPIURL:  "http://server:8080",
			GamePort:    "24642",
			VNCPort:     "5800",
		},
	}
	applyEnvOverrides(cfg)
	return cfg
}
