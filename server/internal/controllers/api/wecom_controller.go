package api

import (
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/services"
	"bbs-go/internal/pkg/errs"
	"fmt"
	"net/url"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

type WeComController struct {
	Ctx iris.Context
}

// 企业微信应用入口
func (c *WeComController) GetEntry() *web.JsonResult {
	// 获取授权URL
	config, err := services.WeComService.GetConfig()
	if err != nil {
		return web.JsonError(err)
	}
	
	// 生成state参数，防止CSRF
	state := uuid.New().String()
	
	// 构建企业微信OAuth2授权URL
	redirectUri := fmt.Sprintf("%s/api/wecom/callback", config.RedirectUri)
	authUrl := fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect",
		config.CorpId,
		url.QueryEscape(redirectUri),
		state,
	)
	
	return web.NewEmptyRspBuilder().
		Put("authUrl", authUrl).
		Put("state", state).
		JsonResult()
}

// 企业微信授权回调
func (c *WeComController) GetCallback() *web.JsonResult {
	code := c.Ctx.URLParam("code")
	
	if code == "" {
		return web.JsonErrorCode(400, "授权码不能为空")
	}
	
	// 1. 获取企业微信用户信息
	wecomUserInfo, err := services.WeComService.GetUserInfoByCode(code)
	if err != nil {
		return web.JsonError(err)
	}
	
	// 2. 只通过工号查找对应用户
	user := services.WeComBindService.FindUserByWeCom(wecomUserInfo)
	
	if user == nil {
		// 未找到对应用户，返回绑定页面所需信息
		return web.NewEmptyRspBuilder().
			Put("needBind", true).
			Put("wecomUserInfo", wecomUserInfo).
			Put("bindUrl", "/mobile/bind/password").
			JsonResult()
	}

	if user.WeComUserId == "" {
		// 找到用户但未绑定，也需要绑定
		return web.NewEmptyRspBuilder().
			Put("needBind", true).
			Put("wecomUserInfo", wecomUserInfo).
			Put("suggestedUsername", user.Username.String).
			Put("bindUrl", "/mobile/bind/password").
			JsonResult()
	}
	
	// 3. 已绑定，直接登录（不再更新登录时间）
	user, err = services.WeComBindService.LoginByWeCom(wecomUserInfo)
	if err != nil {
		return web.JsonError(err)
	}
	
	return render.BuildLoginSuccess(c.Ctx, user, "/mobile")
}

// 绑定企业微信账号
func (c *WeComController) PostBind() *web.JsonResult {
	// 从JSON请求体中读取参数
	var req struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		WecomUserId string `json:"wecomUserId"`
	}

	if err := c.Ctx.ReadJSON(&req); err != nil {
		c.Ctx.Application().Logger().Errorf("读取请求参数失败: %v", err)
		return web.JsonErrorCode(400, "参数格式错误")
	}

	// 添加调试日志
	c.Ctx.Application().Logger().Infof("企业微信绑定请求: username=%s, wecomUserId=%s", req.Username, req.WecomUserId)

	// 参数验证
	if req.Username == "" {
		return web.JsonErrorCode(400, "账号不能为空")
	}
	if req.Password == "" {
		return web.JsonErrorCode(400, "密码不能为空")
	}
	if req.WecomUserId == "" {
		return web.JsonErrorCode(400, "企业微信用户ID不能为空")
	}

	// 获取企业微信用户信息
	c.Ctx.Application().Logger().Infof("正在获取企业微信用户信息: wecomUserId=%s", req.WecomUserId)
	wecomUserInfo, err := services.WeComService.GetUserInfoByUserId(req.WecomUserId)
	if err != nil {
		c.Ctx.Application().Logger().Errorf("获取企业微信用户信息失败: %v", err)
		return web.JsonError(err)
	}
	c.Ctx.Application().Logger().Infof("企业微信用户信息: %+v", wecomUserInfo)

	// 执行绑定（不再记录绑定时间）
	c.Ctx.Application().Logger().Infof("正在绑定用户: username=%s", req.Username)
	user, err := services.WeComBindService.BindWeComByUsername(req.Username, req.Password, wecomUserInfo)
	if err != nil {
		c.Ctx.Application().Logger().Errorf("绑定失败: %v", err)
		return web.JsonError(err)
	}

	c.Ctx.Application().Logger().Infof("绑定成功: userId=%d", user.Id)

	// 绑定成功，直接登录
	return render.BuildLoginSuccess(c.Ctx, user, "/mobile")
}

// 解绑企业微信账号
func (c *WeComController) PostUnbind() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin())
	}

	// 检查是否已绑定企业微信
	if user.WeComUserId == "" {
		return web.JsonErrorCode(400, "您尚未绑定企业微信账号")
	}

	err := services.WeComBindService.UnbindWeCom(user.Id)
	if err != nil {
		return web.JsonError(err)
	}

	return web.JsonSuccess()
}