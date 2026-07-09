package service

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"stardew-panel/database"
	"stardew-panel/models"
	"strings"
)

// InstallService 安装服务
type InstallService struct {
	gamePath  string
	modsPath  string
	savesPath string
}

// NewInstallService 创建安装服务
func NewInstallService(gamePath, modsPath, savesPath string) *InstallService {
	return &InstallService{
		gamePath:  gamePath,
		modsPath:  modsPath,
		savesPath: savesPath,
	}
}

// CheckInstallation 检查是否已安装
func (s *InstallService) CheckInstallation() (*models.Installation, error) {
	row := database.DB.QueryRow(`
		SELECT id, game_path, install_method, has_smapi, version, installed_at, updated_at
		FROM installation
		ORDER BY id DESC LIMIT 1
	`)

	var inst models.Installation
	err := row.Scan(&inst.ID, &inst.GamePath, &inst.InstallMethod, &inst.HasSMAPI, &inst.Version, &inst.InstalledAt, &inst.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &inst, nil
}

// VerifyGamePath 验证游戏路径
func (s *InstallService) VerifyGamePath(path string, detectSMAPI bool) (*models.Installation, error) {
	// 验证路径安全性 - 防止路径注入
	if strings.Contains(path, "..") {
		return nil, fmt.Errorf("非法路径：不允许包含 '..'")
	}

	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("路径不存在: %s", path)
	}

	// 检查是否有 Stardew Valley 可执行文件
	execPath := filepath.Join(path, "StardewValley")
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		// 尝试 .exe 后缀 (Windows)
		execPath = filepath.Join(path, "StardewValley.exe")
		if _, err := os.Stat(execPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("未找到 Stardew Valley 可执行文件")
		}
	}

	// 检测 SMAPI
	hasSMAPI := false
	if detectSMAPI {
		smapiPath := filepath.Join(path, "StardewModdingAPI")
		if _, err := os.Stat(smapiPath); err == nil {
			hasSMAPI = true
		} else {
			// Windows
			smapiPath = filepath.Join(path, "StardewModdingAPI.exe")
			if _, err := os.Stat(smapiPath); err == nil {
				hasSMAPI = true
			}
		}
	}

	// 保存安装信息
	inst := &models.Installation{
		GamePath:      path,
		InstallMethod: "path",
		HasSMAPI:      hasSMAPI,
		Version:       "1.6.9", // TODO: 从游戏文件中读取版本
	}

	if err := s.saveInstallation(inst); err != nil {
		return nil, err
	}

	return inst, nil
}

// ExtractUploadedFile 解压上传的游戏文件
func (s *InstallService) ExtractUploadedFile(zipPath string, installSMAPI bool) (*models.Installation, error) {
	// 确保目标目录存在
	if err := os.MkdirAll(s.gamePath, 0755); err != nil {
		return nil, fmt.Errorf("创建游戏目录失败: %w", err)
	}

	// 解压文件
	if err := s.unzip(zipPath, s.gamePath); err != nil {
		return nil, fmt.Errorf("解压文件失败: %w", err)
	}

	// 安装 SMAPI
	hasSMAPI := false
	if installSMAPI {
		if err := s.installSMAPI(s.gamePath); err != nil {
			// SMAPI 安装失败不应阻塞整个流程
			fmt.Printf("SMAPI 安装失败: %v\n", err)
		} else {
			hasSMAPI = true
		}
	}

	// 保存安装信息
	inst := &models.Installation{
		GamePath:      s.gamePath,
		InstallMethod: "upload",
		HasSMAPI:      hasSMAPI,
		Version:       "1.6.9",
	}

	if err := s.saveInstallation(inst); err != nil {
		return nil, err
	}

	// 删除上传的压缩文件
	os.Remove(zipPath)

	return inst, nil
}

// InstallViaSteamCMD 通过 SteamCMD 安装
func (s *InstallService) InstallViaSteamCMD(installPath, username, password string, installSMAPI bool) (*models.Installation, error) {
	// 验证参数
	if username == "" {
		return nil, fmt.Errorf("Steam 用户名不能为空")
	}
	if password == "" {
		return nil, fmt.Errorf("Steam 密码不能为空")
	}

	// 验证路径安全性
	if strings.Contains(installPath, "..") {
		return nil, fmt.Errorf("非法路径：不允许包含 '..'")
	}

	// 确保安装目录存在
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return nil, fmt.Errorf("创建安装目录失败: %w", err)
	}

	// 检查 SteamCMD 是否安装
	steamcmdPath := s.findSteamCMD()
	if steamcmdPath == "" {
		return nil, fmt.Errorf("未找到 SteamCMD，请先安装 SteamCMD")
	}

	// 使用 SteamCMD 下载游戏
	// App ID 413150 是 Stardew Valley 的专用服务器
	cmd := exec.Command(steamcmdPath,
		"+login", username, password,
		"+force_install_dir", installPath,
		"+app_update", "413150", "validate",
		"+quit",
	)

	// 不显示密码在日志中
	fmt.Printf("正在使用 SteamCMD 下载游戏（用户: %s）...\n", username)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// 过滤输出中的敏感信息
		outputStr := string(output)
		filteredOutput := strings.ReplaceAll(outputStr, password, "****")

		// 检查是否是认证失败
		if strings.Contains(outputStr, "Invalid Password") {
			return nil, fmt.Errorf("Steam 密码错误")
		}
		if strings.Contains(outputStr, "Invalid Username") {
			return nil, fmt.Errorf("Steam 用户名不存在")
		}
		if strings.Contains(outputStr, "Two-factor") || strings.Contains(outputStr, "Steam Guard") {
			return nil, fmt.Errorf("需要 Steam Guard 验证码，请在 SteamCMD 中手动登录后再试")
		}
		if strings.Contains(outputStr, "No subscription") {
			return nil, fmt.Errorf("该账号未购买星露谷物语，无法下载")
		}
		return nil, fmt.Errorf("SteamCMD 执行失败: %w\n%s", err, filteredOutput)
	}

	// 安装 SMAPI
	hasSMAPI := false
	if installSMAPI {
		if err := s.installSMAPI(installPath); err != nil {
			fmt.Printf("SMAPI 安装失败: %v\n", err)
		} else {
			hasSMAPI = true
		}
	}

	// 保存安装信息
	inst := &models.Installation{
		GamePath:      installPath,
		InstallMethod: "steamcmd",
		HasSMAPI:      hasSMAPI,
		Version:       "1.6.9",
	}

	if err := s.saveInstallation(inst); err != nil {
		return nil, err
	}

	return inst, nil
}

// unzip 解压 zip 文件
func (s *InstallService) unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// 检查路径遍历漏洞
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("非法文件路径: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// 创建父目录
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// 创建文件
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// installSMAPI 安装 SMAPI
func (s *InstallService) installSMAPI(gamePath string) error {
	// TODO: 实现 SMAPI 自动安装
	// 1. 下载最新版本的 SMAPI
	// 2. 解压并安装到游戏目录
	// 3. 运行安装脚本

	// 这里返回一个提示，实际安装需要更复杂的逻辑
	return fmt.Errorf("SMAPI 自动安装功能暂未实现，请手动安装")
}

// findSteamCMD 查找 SteamCMD 可执行文件
func (s *InstallService) findSteamCMD() string {
	// 常见的 SteamCMD 路径
	paths := []string{
		"/usr/games/steamcmd",
		"/usr/local/bin/steamcmd",
		"steamcmd",
	}

	for _, path := range paths {
		if _, err := exec.LookPath(path); err == nil {
			return path
		}
	}

	return ""
}

// saveInstallation 保存安装信息
func (s *InstallService) saveInstallation(inst *models.Installation) error {
	result, err := database.DB.Exec(`
		INSERT INTO installation (game_path, install_method, has_smapi, version)
		VALUES (?, ?, ?, ?)
	`, inst.GamePath, inst.InstallMethod, inst.HasSMAPI, inst.Version)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	inst.ID = int(id)
	return nil
}
