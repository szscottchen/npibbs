package services

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/iplocator"
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/repositories"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

var TopicPublishService = new(topicPublishService)

type topicPublishService struct{}

// Publish 发表
func (s *topicPublishService) Publish(userId int64, form models.CreateTopicForm) (*models.Topic, error) {
	if err := s._CheckParams(form); err != nil {
		return nil, err
	}

	// 过滤内容中的非图片文件img标签
	filteredContent := s.filterNonImageFilesInContent(form.Content)

	now := dates.NowTimestamp()
	topic := &models.Topic{
		Type:            form.Type,
		UserId:          userId,
		NodeId:          form.NodeId,
		Title:           form.Title,
		ContentType:     form.ContentType,
		Content:         filteredContent,
		HideContent:     form.HideContent,
		NeedAHand:       form.NeedAHand,
		Status:          constants.StatusOk,
		UserAgent:       form.UserAgent,
		Ip:              form.Ip,
		IpLocation:      iplocator.IpLocation(form.Ip),
		LastCommentTime: now,
		CreateTime:      now,
	}

	if len(form.ImageList) > 0 {
		imageListStr, err := jsons.ToStr(form.ImageList)
		if err == nil {
			topic.ImageList = imageListStr
		} else {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	}

	// 检查是否需要审核
	if s._IsNeedReview(form) {
		topic.Status = constants.StatusReview
	}

	var tagIds []int64
	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		var txErr error
		// 帖子
		if txErr = repositories.TopicRepository.Create(tx, topic); txErr != nil {
			return txErr
		}

		// 标签
		if tagIds, txErr = repositories.TagRepository.GetOrCreates(tx, form.Tags); txErr != nil {
			return txErr
		}
		if txErr = repositories.TopicTagRepository.AddTopicTags(tx, topic.Id, tagIds); txErr != nil {
			return txErr
		}

		// 用户计数
		if txErr = UserService.IncrTopicCount(tx, userId); txErr != nil {
			return txErr
		}

		// 积分
		UserService.IncrScoreForPostTopic(topic)

		return nil
	})
	if err != nil {
		return nil, err
	}
	// 添加索引
	search.UpdateTopicIndex(topic)
	// 发送事件
	event.Send(event.TopicCreateEvent{
		UserId:     topic.UserId,
		TopicId:    topic.Id,
		CreateTime: topic.CreateTime,
	})
	// 如果话题需要支援，刷新缓存
	if topic.NeedAHand == 1 {
		cache.TopicCache.InvalidateNeedAHand()
	}
	return topic, nil
}

// IsNeedReview 是否需要审核
func (s *topicPublishService) _IsNeedReview(form models.CreateTopicForm) bool {
	if hits := ForbiddenWordService.Check(form.Title); len(hits) > 0 {
		slog.Info("帖子标题命中违禁词", slog.String("hits", strings.Join(hits, ",")))
		return true
	}

	if hits := ForbiddenWordService.Check(form.Content); len(hits) > 0 {
		slog.Info("帖子内容命中违禁词", slog.String("hits", strings.Join(hits, ",")))
		return true
	}

	return false
}

func (s topicPublishService) _CheckParams(form models.CreateTopicForm) (err error) {
	modules := SysConfigService.GetModules()
	if form.Type == constants.TopicTypeTweet {
		if !modules.Tweet {
			return errors.New("未开启动态功能")
		}
		if strs.IsBlank(form.Content) {
			return errors.New("内容不能为空")
		}
		// if strs.IsBlank(form.Content) && len(form.ImageList) == 0 {
		// 	return errors.New("内容或图片不能为空")
		// }
	} else {
		if !modules.Topic {
			return errors.New("未开启帖子功能")
		}
		if strs.IsBlank(form.Title) {
			return errors.New("标题不能为空")
		}

		if strs.IsBlank(form.Content) {
			return errors.New("内容不能为空")
		}

		if strs.RuneLen(form.Title) > 128 {
			return errors.New("标题长度不能超过128")
		}
	}

	if form.NodeId <= 0 {
		form.NodeId = SysConfigService.GetConfig().DefaultNodeId
		if form.NodeId <= 0 {
			return errors.New("请选择节点")
		}
	}
	node := repositories.TopicNodeRepository.Get(sqls.DB(), form.NodeId)
	if node == nil || node.Status != constants.StatusOk {
		return errors.New("节点不存在")
	}

	return nil
}

// filterNonImageFilesInContent 过滤内容中的非图片文件img标签，将其转换为链接
func (s *topicPublishService) filterNonImageFilesInContent(content string) string {
	// 正则表达式匹配img标签
	imgRegex := regexp.MustCompile(`(?i)<img\s+[^>]*src\s*=\s*["']([^"']+)["'][^>]*(?:alt\s*=\s*["']([^"']*)["'])?[^>]*>`)
	
	return imgRegex.ReplaceAllStringFunc(content, func(imgTag string) string {
		// 提取src属性
		srcRegex := regexp.MustCompile(`(?i)src\s*=\s*["']([^"']+)["']`)
		srcMatch := srcRegex.FindStringSubmatch(imgTag)
		if len(srcMatch) < 2 {
			return imgTag // 如果无法提取src，保持原样
		}
		src := srcMatch[1]
		
		// 提取alt属性（作为链接文本）
		altRegex := regexp.MustCompile(`(?i)alt\s*=\s*["']([^"']*)["']`)
		altMatch := altRegex.FindStringSubmatch(imgTag)
		var altText string
		if len(altMatch) >= 2 {
			altText = altMatch[1]
		}
		
		// 检查是否为图片文件
		if s.isImageFile(src) {
			return imgTag // 图片文件保持img标签不变
		}
		
		// 非图片文件，转换为链接
		linkText := altText
		if linkText == "" {
			// 如果没有alt文本，从URL中提取文件名
			parts := strings.Split(src, "/")
			linkText = parts[len(parts)-1]
			if linkText == "" {
				linkText = "文件" // 默认文本
			}
		}
		
		return fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, src, linkText)
	})
}

// isImageFile 判断URL是否指向图片文件
func (s *topicPublishService) isImageFile(url string) bool {
	// 图片文件扩展名列表
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg", ".ico"}
	
	urlLower := strings.ToLower(url)
	for _, ext := range imageExts {
		if strings.HasSuffix(urlLower, ext) {
			return true
		}
	}
	
	return false
}
