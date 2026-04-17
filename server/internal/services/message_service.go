package services

import (
	"fmt"
	"log/slog"

	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/email"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/msg"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"github.com/tidwall/gjson"
)

var MessageService = newMessageService()

func newMessageService() *messageService {
	return &messageService{}
}

type messageService struct {
}

func (s *messageService) Get(id int64) *models.Message {
	return repositories.MessageRepository.Get(sqls.DB(), id)
}

func (s *messageService) Take(where ...interface{}) *models.Message {
	return repositories.MessageRepository.Take(sqls.DB(), where...)
}

func (s *messageService) Find(cnd *sqls.Cnd) []models.Message {
	return repositories.MessageRepository.Find(sqls.DB(), cnd)
}

func (s *messageService) FindOne(cnd *sqls.Cnd) *models.Message {
	return repositories.MessageRepository.FindOne(sqls.DB(), cnd)
}

func (s *messageService) FindPageByParams(params *params.QueryParams) (list []models.Message, paging *sqls.Paging) {
	return repositories.MessageRepository.FindPageByParams(sqls.DB(), params)
}

func (s *messageService) FindPageByCnd(cnd *sqls.Cnd) (list []models.Message, paging *sqls.Paging) {
	return repositories.MessageRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *messageService) Create(t *models.Message) error {
	return repositories.MessageRepository.Create(sqls.DB(), t)
}

func (s *messageService) Update(t *models.Message) error {
	return repositories.MessageRepository.Update(sqls.DB(), t)
}

func (s *messageService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.MessageRepository.Updates(sqls.DB(), id, columns)
}

func (s *messageService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.MessageRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *messageService) Delete(id int64) {
	repositories.MessageRepository.Delete(sqls.DB(), id)
}

// GetUnReadCount 获取未读消息数量
func (s *messageService) GetUnReadCount(userId int64) (count int64) {
	sqls.DB().Where("user_id = ? and status = ?", userId, msg.StatusUnread).Model(&models.Message{}).Count(&count)
	return
}

// MarkRead 将所有消息标记为已读
func (s *messageService) MarkRead(userId int64) {
	sqls.DB().Exec("update t_message set status = ? where user_id = ? and status = ?", msg.StatusHaveRead,
		userId, msg.StatusUnread)
}

// SendMsg 发送消息
func (s *messageService) SendMsg(from, to int64, msgType msg.Type,
	title, content, quoteContent string, extraData interface{}) {

	t := &models.Message{
		FromId:       from,
		UserId:       to,
		Title:        title,
		Content:      content,
		QuoteContent: quoteContent,
		Type:         int(msgType),
		ExtraData:    jsons.ToJsonStr(extraData),
		Status:       msg.StatusUnread,
		CreateTime:   dates.NowTimestamp(),
	}
	if err := s.Create(t); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	} else {
		s.SendEmailNotice(t)
		s.SendWeComNotice(t)
	}
}

// SendEmailNotice 发送邮件通知
func (s *messageService) SendEmailNotice(t *models.Message) {
	msgType := msg.Type(t.Type)

	// 话题被删除不发送邮件提醒
	if msgType == msg.TypeTopicDelete {
		return
	}
	user := cache.UserCache.Get(t.UserId)
	if user == nil || len(user.Email.String) == 0 {
		return
	}
	var (
		siteTitle  = cache.SysConfigCache.GetStr(constants.SysConfigSiteTitle)
		emailTitle = siteTitle + " - " + locales.Get("email.new_message")
	)

	if msgType == msg.TypeTopicComment {
		emailTitle = siteTitle + " - " + locales.Get("email.topic_comment")
	} else if msgType == msg.TypeCommentReply {
		emailTitle = siteTitle + " - " + locales.Get("email.comment_reply")
	} else if msgType == msg.TypeTopicLike {
		emailTitle = siteTitle + " - " + locales.Get("email.topic_like")
	} else if msgType == msg.TypeTopicFavorite {
		emailTitle = siteTitle + " - " + locales.Get("email.topic_favorite")
	} else if msgType == msg.TypeTopicRecommend {
		emailTitle = siteTitle + " - " + locales.Get("email.topic_recommend")
	} else if msgType == msg.TypeTopicDelete {
		emailTitle = siteTitle + " - " + locales.Get("email.topic_delete")
	} else if msgType == msg.TypeArticleComment {
		emailTitle = siteTitle + " - " + locales.Get("email.article_comment")
	}

	var from *models.User
	if t.FromId > 0 {
		from = cache.UserCache.Get(t.FromId)
	}
	err := email.SendTemplateEmail(from, user.Email.String, emailTitle, emailTitle, t.Content,
		t.QuoteContent, &dto.ActionLink{
			Title: locales.Get("email.view_details"),
			Url:   bbsurls.AbsUrl("/user/messages"),
		})
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
}

// SendWeComNotice 发送企业微信通知（触发红点提醒）
func (s *messageService) SendWeComNotice(t *models.Message) {
	msgType := msg.Type(t.Type)

	// 话题被删除不发送企业微信提醒
	if msgType == msg.TypeTopicDelete {
		return
	}

	// 获取接收用户
	user := cache.UserCache.Get(t.UserId)
	if user == nil || user.WeComUserId == "" {
		return
	}

	// 构建消息标题和内容
	var (
		siteTitle = cache.SysConfigCache.GetStr(constants.SysConfigSiteTitle)
		title     string
		desc      string
		detailUrl string
	)

	switch msgType {
	case msg.TypeTopicComment:
		title = "收到新评论"
		desc = fmt.Sprintf("%s：%s", t.Title, t.Content)
		detailUrl = s.getMessageUrl(t)
	case msg.TypeArticleComment:
		title = "文章收到新评论"
		desc = fmt.Sprintf("%s：%s", t.Title, t.Content)
		detailUrl = s.getMessageUrl(t)
	case msg.TypeCommentReply:
		title = "收到新回复"
		desc = fmt.Sprintf("%s：%s", t.Title, t.Content)
		detailUrl = s.getMessageUrl(t)
	case msg.TypeTopicLike:
		title = "收到新点赞"
		desc = fmt.Sprintf("有人%s", t.Title)
		detailUrl = s.getMessageUrl(t)
	case msg.TypeTopicFavorite:
		title = "话题被收藏"
		desc = fmt.Sprintf("有人%s", t.Title)
		detailUrl = s.getMessageUrl(t)
	case msg.TypeTopicRecommend:
		title = "话题被推荐"
		desc = fmt.Sprintf("你的%s", t.Title)
		detailUrl = s.getMessageUrl(t)
	default:
		return
	}

	// 发送企业微信应用消息
	err := WeComService.SendAppMessage(user.WeComUserId, siteTitle+" - "+title, desc, detailUrl)
	if err != nil {
		slog.Error("发送企业微信通知失败", slog.Any("err", err), slog.Int64("userId", t.UserId))
	}
}

// getMessageUrl 获取消息对应的跳转链接
func (s *messageService) getMessageUrl(t *models.Message) string {
	msgType := msg.Type(t.Type)
	switch msgType {
	case msg.TypeTopicComment, msg.TypeArticleComment:
		entityType := gjson.Get(t.ExtraData, "entityType")
		entityId := gjson.Get(t.ExtraData, "entityId")
		if entityType.String() == "article" {
			return bbsurls.AbsUrl(fmt.Sprintf("/article/%d", entityId.Int()))
		} else {
			return bbsurls.AbsUrl(fmt.Sprintf("/topic/%d", entityId.Int()))
		}
	case msg.TypeCommentReply:
		entityType := gjson.Get(t.ExtraData, "rootEntityType")
		entityId := gjson.Get(t.ExtraData, "rootEntityId")
		if entityType.String() == "article" {
			return bbsurls.AbsUrl(fmt.Sprintf("/article/%d", entityId.Int()))
		} else {
			return bbsurls.AbsUrl(fmt.Sprintf("/topic/%d", entityId.Int()))
		}
	case msg.TypeTopicLike, msg.TypeTopicFavorite, msg.TypeTopicRecommend:
		topicId := gjson.Get(t.ExtraData, "topicId")
		if topicId.Exists() && topicId.Int() > 0 {
			return bbsurls.AbsUrl(fmt.Sprintf("/topic/%d", topicId.Int()))
		}
	}
	return bbsurls.AbsUrl("/user/messages")
}
