package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
	"fmt"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/passwd"
	"github.com/mlogclub/simple/sqls"
)

type weComBindService struct{}

var WeComBindService = &weComBindService{}

// 通过企业微信信息查找对应用户（只通过工号）
func (s *weComBindService) FindUserByWeCom(wecomUserInfo *WeComUserInfo) *models.User {
	// 1. 首先检查是否已绑定
	if user := s.GetUserByWeComUserId(wecomUserInfo.UserId); user != nil {
		return user
	}
	
	// 2. 只通过工号查找（工号 = username）
	if user := s.GetUserByUsername(wecomUserInfo.UserId); user != nil {
		return user
	}
	
	return nil
}

// 检查用户是否已绑定企业微信
func (s *weComBindService) GetUserByWeComUserId(wecomUserId string) *models.User {
	return repositories.UserRepository.GetByWeComUserId(sqls.DB(), wecomUserId)
}

// 通过用户名查找用户
func (s *weComBindService) GetUserByUsername(username string) *models.User {
	return repositories.UserRepository.GetByUsername(sqls.DB(), username)
}

// 绑定企业微信账号（简化版，不记录时间）
func (s *weComBindService) BindWeCom(userId int64, wecomUserInfo *WeComUserInfo, password string) (*models.User, error) {
	// 1. 获取用户信息
	user := repositories.UserRepository.Get(sqls.DB(), userId)
	if user == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	
	// 2. 验证密码
	if !passwd.ValidatePassword(user.Password, password) {
		return nil, fmt.Errorf("密码错误")
	}
	
	// 3. 检查企业微信账号是否已被其他用户绑定
	existUser := s.GetUserByWeComUserId(wecomUserInfo.UserId)
	if existUser != nil && existUser.Id != userId {
		return nil, fmt.Errorf("该企业微信账号已被其他用户绑定")
	}
	
	// 4. 绑定企业微信（只更新we_com_user_id字段）
	updateColumns := map[string]interface{}{
		"we_com_user_id": wecomUserInfo.UserId,  // 改为 we_com_user_id
		"update_time": dates.NowTimestamp(),
	}
	
	err := repositories.UserRepository.UpdateColumns(sqls.DB(), userId, updateColumns)
	if err != nil {
		return nil, fmt.Errorf("绑定失败: %v", err)
	}
	
	// 5. 返回更新后的用户信息
	return repositories.UserRepository.Get(sqls.DB(), userId), nil
}

// 企业微信登录（简化版，不记录登录时间）
func (s *weComBindService) LoginByWeCom(wecomUserInfo *WeComUserInfo) (*models.User, error) {
	user := s.GetUserByWeComUserId(wecomUserInfo.UserId)
	if user == nil {
		return nil, fmt.Errorf("企业微信账号未绑定")
	}
	
	return user, nil
}

// 通过用户名和企业微信信息绑定
func (s *weComBindService) BindWeComByUsername(username, password string, wecomUserInfo *WeComUserInfo) (*models.User, error) {
	// 1. 根据用户名查找用户
	user := s.GetUserByUsername(username)
	if user == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	
	// 2. 验证密码
	if !passwd.ValidatePassword(user.Password, password) {
		return nil, fmt.Errorf("密码错误")
	}
	
	// 3. 检查企业微信账号是否已被其他用户绑定
	existUser := s.GetUserByWeComUserId(wecomUserInfo.UserId)
	if existUser != nil && existUser.Id != user.Id {
		return nil, fmt.Errorf("该企业微信账号已被其他用户绑定")
	}
	
	// 4. 绑定企业微信
	updateColumns := map[string]interface{}{
		"we_com_user_id": wecomUserInfo.UserId,  // 改为 we_com_user_id
		"update_time":  dates.NowTimestamp(),
	}
	
	err := repositories.UserRepository.UpdateColumns(sqls.DB(), user.Id, updateColumns)
	if err != nil {
		return nil, fmt.Errorf("绑定失败: %v", err)
	}
	
	// 5. 返回更新后的用户信息
	return repositories.UserRepository.Get(sqls.DB(), user.Id), nil
}

// 解绑企业微信账号
func (s *weComBindService) UnbindWeCom(userId int64) error {
	updateColumns := map[string]interface{}{
		"we_com_user_id": "",  // 改为 we_com_user_id
		"update_time":  dates.NowTimestamp(),
	}
	
	return repositories.UserRepository.UpdateColumns(sqls.DB(), userId, updateColumns)
}