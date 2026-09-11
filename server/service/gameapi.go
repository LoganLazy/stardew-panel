package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GameAPIClient 封装对 sdvd/server 游戏容器 REST API 的调用。
// 参考: https://github.com/stardew-valley-dedicated-server/server (docs/features/rest-api.md)
type GameAPIClient struct {
	baseURL string // 如 http://server:8080
	apiKey  string // 对应 sdvd 的 API_KEY，非空时走 Bearer 认证
	client  *http.Client
}

// NewGameAPIClient 创建游戏 API 客户端。baseURL 为空表示未配置（未编排）。
func NewGameAPIClient(baseURL, apiKey string) *GameAPIClient {
	return &GameAPIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// Configured 是否已配置游戏 API 地址。
func (g *GameAPIClient) Configured() bool {
	return g.baseURL != ""
}

// GameStatus 对应 sdvd /status 端点返回。只取 panel 需要的字段。
type GameStatus struct {
	PlayerCount     int    `json:"playerCount"`
	MaxPlayers      int    `json:"maxPlayers"`
	SteamInviteCode string `json:"steamInviteCode"`
	GogInviteCode   string `json:"gogInviteCode"`
	ServerVersion   string `json:"serverVersion"`
	IsOnline        bool   `json:"isOnline"`
	IsReady         bool   `json:"isReady"`
	FarmName        string `json:"farmName"`
	Day             int    `json:"day"`
	Season          string `json:"season"`
	Year            int    `json:"year"`
	TimeOfDay       int    `json:"timeOfDay"`
	IsPaused        bool   `json:"isPaused"`
}

// InviteCode 返回可用的邀请码（优先 Steam，其次 GOG）。
func (s *GameStatus) InviteCode() string {
	if s.SteamInviteCode != "" {
		return s.SteamInviteCode
	}
	return s.GogInviteCode
}

// GamePlayer 对应 sdvd /players 里的单个玩家。
type GamePlayer struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsOnline bool   `json:"isOnline"`
}

type playersResponse struct {
	Players []GamePlayer `json:"players"`
	Version int64        `json:"version"`
}

// get 发起一次 GET 请求并把 JSON 解到 out。
func (g *GameAPIClient) get(path string, out interface{}) error {
	if !g.Configured() {
		return fmt.Errorf("游戏 API 未配置")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.baseURL+path, nil)
	if err != nil {
		return err
	}
	if g.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+g.apiKey)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("请求游戏 API %s 失败: %w", path, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("游戏 API %s 返回 %d: %s", path, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("解析游戏 API %s 响应失败: %w", path, err)
	}
	return nil
}

// GetStatus 调用 sdvd /status。
func (g *GameAPIClient) GetStatus() (*GameStatus, error) {
	var st GameStatus
	if err := g.get("/status", &st); err != nil {
		return nil, err
	}
	return &st, nil
}

// GetPlayers 调用 sdvd /players，返回在线玩家列表。
func (g *GameAPIClient) GetPlayers() ([]GamePlayer, error) {
	var pr playersResponse
	if err := g.get("/players", &pr); err != nil {
		return nil, err
	}
	return pr.Players, nil
}

// Healthy 调用 sdvd /health（无需认证），判断游戏 API 是否已就绪。
func (g *GameAPIClient) Healthy() bool {
	if !g.Configured() {
		return false
	}
	var discard map[string]interface{}
	// /health 是公开端点，失败即视为未就绪
	return g.get("/health", &discard) == nil
}

// GetSettings 调用 sdvd GET /settings，返回全部服务器设置（原样透传的 map）。
func (g *GameAPIClient) GetSettings() (map[string]interface{}, error) {
	var settings map[string]interface{}
	if err := g.get("/settings", &settings); err != nil {
		return nil, err
	}
	return settings, nil
}

// UpdateSetting 调用 sdvd PUT /settings/{key}，更新单个设置项。
// sdvd 约定请求体为 {"value": <新值>}。
func (g *GameAPIClient) UpdateSetting(key string, value interface{}) error {
	body := map[string]interface{}{"value": value}
	return g.put("/settings/"+key, body)
}

// put 发起一次 PUT 请求，body 序列化为 JSON 发送。不解析响应体，只看状态码。
func (g *GameAPIClient) put(path string, body interface{}) error {
	if !g.Configured() {
		return fmt.Errorf("游戏 API 未配置")
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, g.baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if g.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+g.apiKey)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("请求游戏 API %s 失败: %w", path, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("游戏 API %s 返回 %d: %s", path, resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}
