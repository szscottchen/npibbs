package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/pkg/config"
	"bbs-go/internal/services"
)

type ConfigController struct {
	Ctx iris.Context
}

type SysConfigController struct {
	Ctx iris.Context
}

func (c *ConfigController) GetConfigs() *web.JsonResult {
	cfg := config.Instance

	var b *web.RspBuilder
	if cfg.Installed {
		sysConfig := services.SysConfigService.GetConfig()
		sysConfig.UploadConfig = nil
		b = web.NewRspBuilder(sysConfig)
	} else {
		b = web.NewEmptyRspBuilder()
	}
	b.Put("installed", cfg.Installed)
	b.Put("language", cfg.Language)
	return b.JsonResult()
}

// GetConfigs 获取系统配置（包含价值类型）
func (c *SysConfigController) GetConfigs() *web.JsonResult {
	config := services.SysConfigService.GetConfig()
	return web.JsonData(config)
}
