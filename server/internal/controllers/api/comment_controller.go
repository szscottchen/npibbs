package api

import (
	"fmt"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"bbs-go/internal/spam"
)

type CommentController struct {
	Ctx iris.Context
}

func (c *CommentController) GetComments() *web.JsonResult {
	var (
		err        error
		cursor     int64
		entityType string
		entityId   int64
	)
	cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)

	if entityType, err = params.FormValueRequired(c.Ctx, "entityType"); err != nil {
		return web.JsonError(err)
	}
	if entityId, err = params.FormValueInt64(c.Ctx, "entityId"); err != nil {
		return web.JsonError(err)
	}
	currentUser := services.UserTokenService.GetCurrent(c.Ctx)
	comments, cursor, hasMore := services.CommentService.GetComments(entityType, entityId, cursor)
	return web.JsonCursorData(render.BuildComments(comments, currentUser, true, false), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *CommentController) GetReplies() *web.JsonResult {
	var (
		cursor    = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		commentId = params.FormValueInt64Default(c.Ctx, "commentId", 0)
	)
	currentUser := services.UserTokenService.GetCurrent(c.Ctx)
	comments, cursor, hasMore := services.CommentService.GetReplies(commentId, cursor, 10)
	return web.JsonCursorData(render.BuildComments(comments, currentUser, false, true), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *CommentController) PostCreate() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}
	form := models.GetCreateCommentForm(c.Ctx)
	if err := spam.CheckComment(user, form); err != nil {
		return web.JsonError(err)
	}

	comment, err := services.CommentService.Publish(user.Id, form)
	if err != nil {
		return web.JsonError(err)
	}

	return web.JsonData(render.BuildComment(comment))
}

// PostValue 价值评价
func (c *CommentController) PostValue() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonErrorCode(401, "请先登录")
	}
	
	commentId := params.FormValueInt64Default(c.Ctx, "commentId", 0)
	valueType := params.FormValueIntDefault(c.Ctx, "valueType", -1)
	
	if commentId <= 0 {
		return web.JsonErrorCode(400, "评论ID不能为空")
	}
	if valueType < 0 {
		return web.JsonErrorCode(400, "价值类型不能为空")
	}
	
	// 获取评论
	comment := services.CommentService.Get(commentId)
	if comment == nil {
		return web.JsonErrorCode(404, "评论不存在")
	}
	
	// 检查是否是主贴发布人
	topic := services.TopicService.Get(comment.EntityId)
	if topic == nil {
		return web.JsonErrorCode(404, "话题不存在")
	}
	if topic.UserId != user.Id {
		return web.JsonErrorCode(403, "只有主贴发布人才能评价价值")
	}
	
	// 不能评价自己的评论
	if comment.UserId == user.Id {
		return web.JsonErrorCode(400, "不能评价自己的评论")
	}
	
	// 获取价值类型配置
	config := services.SysConfigService.GetConfig()
	if config.ValueTypes == nil || len(config.ValueTypes) <= valueType {
		return web.JsonErrorCode(400, "价值类型不存在")
	}
	
	valueTypeConfig := config.ValueTypes[valueType]
	if valueTypeConfig.Score <= 0 {
		return web.JsonErrorCode(400, "价值积分配置错误")
	}
	
	// 检查是否已经评价过
	sourceId := strconv.FormatInt(commentId, 10)
	if services.UserScoreLogService.HasValueScore(user.Id, sourceId) {
		return web.JsonErrorCode(400, "已经对该评论进行过价值评价")
	}
	
	// 更新评论的价值标识
	err := services.CommentService.UpdateColumn(commentId, "valuable", valueTypeConfig.Label)
	if err != nil {
		return web.JsonError(err)
	}
	
	// 给用户增加积分
	err = services.UserService.IncrScoreForValue(comment.UserId, valueTypeConfig.Score, sourceId, 
		fmt.Sprintf("评论被评价为[%s]", valueTypeConfig.Label))
	if err != nil {
		return web.JsonError(err)
	}
	
	return web.JsonSuccess()
}