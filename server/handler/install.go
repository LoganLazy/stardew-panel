package handler

import (
	"net/http"
	"stardew-panel/config"
	"stardew-panel/service"

	"github.com/gin-gonic/gin"
)

var installService *service.InstallService

// InitInstallHandler 初始化安装处理器
func InitInstallHandler(cfg *config.Config) {
	installService = service.NewInstallService(cfg.Game.ComposeFile, cfg.Game.ComposeDir)
}

// GetSetupStatus 返回安装向导所需的就绪状态：
// Steam 账号是否已配置、docker 是否可用、游戏容器是否已就绪。
// 容器模式下不再有"上传/SteamCMD 下载"这类由 panel 执行的安装动作，
// 游戏文件由 steam-auth 容器负责，首次需在宿主手动过 Steam Guard。
func GetSetupStatus(c *gin.Context) {
	c.JSON(http.StatusOK, installService.GetSetupStatus())
}

// CheckInstallation 兼容旧前端的安装检查端点。
// 映射为：Steam 已配置 + docker 可用即视为"已安装"。
func CheckInstallation(c *gin.Context) {
	st := installService.GetSetupStatus()
	c.JSON(http.StatusOK, gin.H{
		"installed": st.SteamConfigured && st.DockerAvailable,
		"setup":     st,
	})
}
