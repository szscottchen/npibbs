package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/services"
	"context"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/mlogclub/simple/web"
)

type AISumController struct {
	Ctx iris.Context
}

// BeforeActivation 在控制器激活前注册路由
func (c *AISumController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/topic/{topicId}/generate-summary", "GenerateSummary")
	b.Handle("GET", "/topic/{topicId}/summary", "GetSummary")
	b.Handle("GET", "/topic/{topicId}/summary-history", "GetSummaryHistory")
}

// GenerateSummary 生成AI总结
func (c *AISumController) GenerateSummary(topicIdStr string) *web.JsonResult {
	topicId, err := strconv.ParseInt(topicIdStr, 10, 64)
	if err != nil || topicId <= 0 {
		return web.JsonErrorMsg(locales.Get("common.not_found"))
	}

	// 检查用户权限
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin())
	}

	// 检查话题是否存在
	topic := services.TopicService.Get(topicId)
	if topic == nil {
		return web.JsonErrorMsg(locales.Get("common.not_found"))
	}

	// 检查用户是否有权限生成总结
	if !c.hasPermission(user.Id, topicId) {
		return web.JsonErrorMsg(locales.Get("topic.no_permission"))
	}

	// 生成AI总结
	ctx := context.Background()
	aiSum, err := services.AISumService.GenerateSummary(ctx, topicId, user.Id)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	return web.JsonData(map[string]interface{}{
		"summaryId": aiSum.Id,
		"content":   aiSum.SumContent,
		"sumTime":   aiSum.SumTime,
	})
}

// GetSummary 获取当前有效的AI总结
func (c *AISumController) GetSummary(topicIdStr string) *web.JsonResult {
	topicId, err := strconv.ParseInt(topicIdStr, 10, 64)
	if err != nil || topicId <= 0 {
		return web.JsonErrorMsg(locales.Get("common.not_found"))
	}

	// 检查话题是否存在
	topic := services.TopicService.Get(topicId)
	if topic == nil {
		return web.JsonErrorMsg(locales.Get("common.not_found"))
	}

	// 获取当前有效的总结
	aiSum := services.AISumService.GetSummary(topicId)
	if aiSum == nil {
		return web.JsonErrorMsg("该话题暂无AI总结")
	}

	return web.JsonData(map[string]interface{}{
		"summaryId": aiSum.Id,
		"content":   aiSum.SumContent,
		"sumTime":   aiSum.SumTime,
		"topicId":   aiSum.TopicId,
	})
}

// GetSummaryHistory 获取AI总结历史
func (c *AISumController) GetSummaryHistory(topicIdStr string) *web.JsonResult {
	topicId, err := strconv.ParseInt(topicIdStr, 10, 64)
	if err != nil || topicId <= 0 {
		return web.JsonErrorMsg(locales.Get("common.not_found"))
	}

	// 检查话题是否存在
	topic := services.TopicService.Get(topicId)
	if topic == nil {
		return web.JsonErrorMsg(locales.Get("common.not_found"))
	}

	// 获取总结历史
	summaries, err := services.AISumService.GetSummaryHistory(topicId)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	var result []map[string]interface{}
	for _, summary := range summaries {
		result = append(result, map[string]interface{}{
			"summaryId": summary.Id,
			"content":   summary.SumContent,
			"sumTime":   summary.SumTime,
			"valid":     summary.SumValid,
		})
	}

	return web.JsonData(result)
}

// hasPermission 检查用户是否有权限生成AI总结
func (c *AISumController) hasPermission(userId, topicId int64) bool {
	// 管理员和话题作者有权限
	topic := services.TopicService.Get(topicId)
	if topic == nil {
		return false
	}

	user := services.UserService.Get(userId)
	if user == nil {
		return false
	}

	// 检查是否是管理员
	if user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return true
	}

	// 检查是否是话题作者
	if topic.UserId == userId {
		return true
	}

	return false
}
