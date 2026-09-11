package service

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// InstallService 安装/就绪检查服务。
//
// 改造后不再自己解压 / 跑 SteamCMD。游戏文件由 sdvd 的 steam-auth 容器用
// 用户的正版 Steam 账号下载，首次需在宿主手动跑一次交互式命令过 Steam Guard：
//
//	docker compose run --rm -it steam-auth setup
//
// 因此 panel 的"安装"页职责简化为：引导 + 就绪状态检查——
//   - 读挂载进来的 .env，判断 Steam 账号是否已配置（只回显用户名，绝不回显密码）
//   - 用 compose ps 判断 steam-auth 容器是否已就绪（游戏文件下载由它负责）
type InstallService struct {
	compose *ComposeRunner
	envPath string // COMPOSE_DIR/.env，只读挂载
}

// NewInstallService 创建安装服务。
func NewInstallService(composeFile, composeDir string) *InstallService {
	return &InstallService{
		compose: NewComposeRunner(composeFile, composeDir),
		envPath: filepath.Join(composeDir, ".env"),
	}
}

// SetupStatus 安装向导状态，返回给前端。
type SetupStatus struct {
	SteamConfigured bool   `json:"steam_configured"` // .env 里是否已填 STEAM_USERNAME
	SteamUsername   string `json:"steam_username"`   // 已配置的 Steam 用户名（脱敏用，不含密码）
	GameReady       bool   `json:"game_ready"`       // steam-auth 容器是否已就绪
	DockerAvailable bool   `json:"docker_available"` // 宿主 docker compose 是否可用
}

// GetSetupStatus 返回当前安装/就绪状态。
func (s *InstallService) GetSetupStatus() *SetupStatus {
	st := &SetupStatus{}

	// Steam 账号是否已在 .env 配置
	username := s.readEnvKey("STEAM_USERNAME")
	st.SteamUsername = username
	st.SteamConfigured = username != ""

	// docker 是否可用
	if err := s.compose.Available(); err == nil {
		st.DockerAvailable = true
		// steam-auth 容器在跑即视为游戏文件已就绪（它负责下载游戏）
		if running, err := s.compose.PSRunning(steamAuthServiceName); err == nil {
			st.GameReady = running
		}
	}

	return st
}

// readEnvKey 从挂载的 .env 读取指定键的值。文件不存在或无此键返回空字符串。
// 只做最简单的 KEY=VALUE 解析，去掉两侧引号和空白。
func (s *InstallService) readEnvKey(key string) string {
	f, err := os.Open(s.envPath)
	if err != nil {
		return ""
	}
	defer f.Close()

	prefix := key + "="
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || !strings.HasPrefix(line, prefix) {
			continue
		}
		val := strings.TrimPrefix(line, prefix)
		val = strings.TrimSpace(val)
		val = strings.Trim(val, `"'`)
		return val
	}
	return ""
}
