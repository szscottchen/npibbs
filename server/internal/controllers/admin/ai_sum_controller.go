package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/services"
)

type AISumController struct {
	Ctx iris.Context
}

func (c *AISumController) GetBy(id int64) *web.JsonResult {
	summary := services.AISumService.GetSummary(id)
	if summary == nil {
		return web.JsonErrorMsg("AI总结不存在，id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(summary)
}

func (c *AISumController) AnyList() *web.JsonResult {
	list, paging := services.AISumService.FindPageByParams(params.NewQueryParams(c.Ctx).
		EqByReq("topic_id").EqByReq("sum_valid").PageByReq().Desc("sum_time"))

	var results []map[string]interface{}
	for _, summary := range list {
		item := map[string]interface{}{
			"id":         summary.Id,
			"topicId":    summary.TopicId,
			"sumTime":    summary.SumTime,
			"sumContent": summary.SumContent,
			"sumValid":   summary.SumValid,
			"createdAt":  summary.CreatedAt,
			"updatedAt":  summary.UpdatedAt,
		}

		// 获取话题信息
		topic := services.TopicService.Get(summary.TopicId)
		if topic != nil {
			item["topicTitle"] = topic.Title
		}

		results = append(results, item)
	}

	return web.JsonData(&web.PageResult{Results: results, Page: paging})
}

// GenerateSummary 管理员生成AI总结
func (c *AISumController) PostGenerateSummary() *web.JsonResult {
	topicId, err := params.FormValueInt64(c.Ctx, "topicId")
	if err != nil {
		return web.JsonError(err)
	}

	// 检查话题是否存在
	topic := services.TopicService.Get(topicId)
	if topic == nil {
		return web.JsonErrorMsg("话题不存在")
	}

	// 生成AI总结
	aiSum, err := services.AISumService.GenerateSummary(c.Ctx.Request().Context(), topicId, 0)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	return web.JsonData(map[string]interface{}{
		"summaryId": aiSum.Id,
		"content":   aiSum.SumContent,
		"sumTime":   aiSum.SumTime,
	})
}

// InvalidateSummary 使总结失效
func (c *AISumController) PostInvalidateSummary() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}

	summary := services.AISumService.GetSummary(id)
	if summary == nil {
		return web.JsonErrorMsg("AI总结不存在")
	}

	err = services.AISumService.InvalidateSummary(id)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	return web.JsonSuccess()
}

// Delete 删除AI总结
func (c *AISumController) PostDelete() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}

	err = services.AISumService.Delete(id)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	return web.JsonSuccess()
}