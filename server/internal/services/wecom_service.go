package services

import (
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/mlogclub/simple/sqls"
)

type weComService struct{}

var WeComService = &weComService{}

type WeComTokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserId      string `json:"UserId"`
}

type WeComUserInfo struct {
	UserId     string `json:"userid"`
	Name       string `json:"name"`
	Department []int  `json:"department"`
	Position   string `json:"position"`
	Mobile     string `json:"mobile"`
	Gender     string `json:"gender"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	Status     int    `json:"status"`
	Enable     int    `json:"enable"`
}

// 企业微信配置结构
type WeComConfig struct {
	Enabled    bool   `json:"enabled"`
	CorpId     string `json:"corpId"`
	Secret     string `json:"secret"`
	AgentId    string `json:"agentId"`
	RedirectUri string `json:"redirectUri"`
}

// 获取企业微信配置
func (s *weComService) GetConfig() (*WeComConfig, error) {
	// 优先从配置文件读取（如果配置文件中有CorpId等关键配置）
	wecomConfig := &WeComConfig{}
	if config.Instance.WeCom.CorpId != "" {
		wecomConfig.Enabled = config.Instance.WeCom.Enabled
		wecomConfig.CorpId = config.Instance.WeCom.CorpId
		wecomConfig.Secret = config.Instance.WeCom.Secret
		wecomConfig.AgentId = config.Instance.WeCom.AgentId
		wecomConfig.RedirectUri = config.Instance.WeCom.RedirectUri

		// 检查必要配置
		if wecomConfig.Enabled {
			if wecomConfig.CorpId == "" || wecomConfig.Secret == "" || wecomConfig.AgentId == "" {
				return nil, fmt.Errorf("企业微信配置不完整")
			}
			return wecomConfig, nil
		}
	}

	// 如果配置文件中没有启用，从t_sys_config表获取
	// 获取启用状态
	if enabledConfig := repositories.SysConfigRepository.GetByKey(sqls.DB(), "wecom.enabled"); enabledConfig != nil {
		wecomConfig.Enabled = enabledConfig.Value == "true"
	} else {
		wecomConfig.Enabled = false
	}

	// 获取其他配置项
	if corpIdConfig := repositories.SysConfigRepository.GetByKey(sqls.DB(), "wecom.corp_id"); corpIdConfig != nil {
		wecomConfig.CorpId = corpIdConfig.Value
	}

	if secretConfig := repositories.SysConfigRepository.GetByKey(sqls.DB(), "wecom.secret"); secretConfig != nil {
		wecomConfig.Secret = secretConfig.Value
	}

	if agentIdConfig := repositories.SysConfigRepository.GetByKey(sqls.DB(), "wecom.agent_id"); agentIdConfig != nil {
		wecomConfig.AgentId = agentIdConfig.Value
	}

	if redirectUriConfig := repositories.SysConfigRepository.GetByKey(sqls.DB(), "wecom.redirect_uri"); redirectUriConfig != nil {
		wecomConfig.RedirectUri = redirectUriConfig.Value
	}

	// 检查必要配置
	if !wecomConfig.Enabled {
		return nil, fmt.Errorf("企业微信登录未启用")
	}

	if wecomConfig.CorpId == "" || wecomConfig.Secret == "" || wecomConfig.AgentId == "" {
		return nil, fmt.Errorf("企业微信配置不完整")
	}

	return wecomConfig, nil
}

// InitConfig 初始化企业微信配置，将配置文件中的配置同步到数据库
func (s *weComService) InitConfig() error {
	// 如果配置文件中没有配置，不执行同步
	if config.Instance.WeCom.CorpId == "" {
		return nil
	}

	// 同步配置到数据库
	_ = SysConfigService.Set("wecom.enabled", fmt.Sprintf("%v", config.Instance.WeCom.Enabled))
	_ = SysConfigService.Set("wecom.corp_id", config.Instance.WeCom.CorpId)
	_ = SysConfigService.Set("wecom.secret", config.Instance.WeCom.Secret)
	_ = SysConfigService.Set("wecom.agent_id", config.Instance.WeCom.AgentId)
	_ = SysConfigService.Set("wecom.redirect_uri", config.Instance.WeCom.RedirectUri)

	return nil
}

// 获取企业微信Access Token
func (s *weComService) GetCorpToken() (*WeComTokenResponse, error) {
	config, err := s.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("获取企业微信配置失败: %v", err)
	}
	
	url := fmt.Sprintf(
		"https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s",
		config.CorpId,
		config.Secret,
	)
	
	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求企业微信API失败: %v", err)
	}
	
	var tokenResp WeComTokenResponse
	if err := json.Unmarshal(resp.Body(), &tokenResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	
	if tokenResp.ErrCode != 0 {
		return nil, fmt.Errorf("企业微信API错误: %d - %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}
	
	return &tokenResp, nil
}

// 通过授权码获取用户信息
func (s *weComService) GetUserInfoByCode(code string) (*WeComUserInfo, error) {
	// 获取access_token
	tokenResp, err := s.GetCorpToken()
	if err != nil {
		return nil, err
	}
	
	// 通过code获取用户ID
	userUrl := fmt.Sprintf(
		"https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s",
		tokenResp.AccessToken,
		code,
	)
	
	client := resty.New()
	resp, err := client.R().Get(userUrl)
	if err != nil {
		return nil, fmt.Errorf("获取用户ID失败: %v", err)
	}
	
	var userResp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		UserId  string `json:"userid"`
	}
	
	err = json.Unmarshal(resp.Body(), &userResp)
	if err != nil {
		return nil, fmt.Errorf("解析用户ID响应失败: %v", err)
	}
	
	if userResp.ErrCode != 0 {
		return nil, fmt.Errorf("获取用户ID失败: %d - %s", userResp.ErrCode, userResp.ErrMsg)
	}
	
	// 获取详细用户信息
	detailUrl := fmt.Sprintf(
		"https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s",
		tokenResp.AccessToken,
		userResp.UserId,
	)
	
	detailResp, err := client.R().Get(detailUrl)
	if err != nil {
		return nil, fmt.Errorf("获取用户详细信息失败: %v", err)
	}
	
	var userInfo WeComUserInfo
	err = json.Unmarshal(detailResp.Body(), &userInfo)
	if err != nil {
		return nil, fmt.Errorf("解析用户详细信息失败: %v", err)
	}
	
	return &userInfo, nil
}

// 根据用户ID获取用户信息
func (s *weComService) GetUserInfoByUserId(userId string) (*WeComUserInfo, error) {
	tokenResp, err := s.GetCorpToken()
	if err != nil {
		return nil, err
	}

	detailUrl := fmt.Sprintf(
		"https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s",
		tokenResp.AccessToken,
		userId,
	)

	client := resty.New()
	resp, err := client.R().Get(detailUrl)
	if err != nil {
		return nil, fmt.Errorf("获取用户详细信息失败: %v", err)
	}

	// 解析响应并检查错误码
	var respBody struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		WeComUserInfo
	}

	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return nil, fmt.Errorf("解析用户详细信息失败: %v", err)
	}

	if respBody.ErrCode != 0 {
		return nil, fmt.Errorf("获取用户详细信息失败: %d - %s", respBody.ErrCode, respBody.ErrMsg)
	}

	// 复制用户信息
	userInfo := respBody.WeComUserInfo
	return &userInfo, nil
}

// SendAppMessage 发送企业微信应用消息（触发红点提醒）
// touser: 企业微信用户ID
// title: 消息标题
// description: 消息描述
// url: 点击后跳转的链接
func (s *weComService) SendAppMessage(touser, title, description, url string) error {
	if touser == "" {
		return fmt.Errorf("企业微信用户ID为空")
	}

	config, err := s.GetConfig()
	if err != nil {
		return fmt.Errorf("获取企业微信配置失败: %v", err)
	}

	tokenResp, err := s.GetCorpToken()
	if err != nil {
		return err
	}

	// 调用企业微信发送应用消息接口
	apiUrl := fmt.Sprintf(
		"https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s",
		tokenResp.AccessToken,
	)

	// 构建文本卡片消息
	messageData := map[string]interface{}{
		"touser":  touser,
		"msgtype": "textcard",
		"agentid": config.AgentId,
		"textcard": map[string]string{
			"title":       title,
			"description": description,
			"url":         url,
			"btntxt":      "查看详情",
		},
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(messageData).
		Post(apiUrl)

	if err != nil {
		return fmt.Errorf("发送企业微信消息失败: %v", err)
	}

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("企业微信API错误: %d - %s", result.ErrCode, result.ErrMsg)
	}

	return nil
}