package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
	"bbs-go/internal/services/ai"
	"context"
	"fmt"
	"time"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var AISumService = newAISumService()

type aiSumService struct{}

func newAISumService() *aiSumService {
	return &aiSumService{}
}

// GenerateSummary 生成AI总结
func (s *aiSumService) GenerateSummary(ctx context.Context, topicId int64, userId int64) (*models.AISum, error) {
	// 检查话题是否存在
	topic := TopicService.Get(topicId)
	if topic == nil {
		return nil, fmt.Errorf("话题不存在")
	}

	// 获取未总结的评论
	comments, totalChars := repositories.CommentRepository.GetUnsummarizedComments(sqls.DB(), topicId)
	if len(comments) == 0 {
		return nil, fmt.Errorf("没有需要总结的评论")
	}

	// 检查是否超过字符限制
	if totalChars > 8000 {
		return nil, fmt.Errorf("评论内容过长，无法进行总结")
	}

	// 获取最新的有效总结
	var previousSummary *models.AISum
	previousSummary = repositories.AISumRepository.GetValidSummary(sqls.DB(), topicId)

	// 构建总结内容
	summaryContent, err := s.buildSummaryContent(topic, comments, previousSummary)
	if err != nil {
		return nil, fmt.Errorf("构建总结内容失败: %v", err)
	}

	// 调用AI服务生成总结
	summaryText, err := ai.AIService.GenerateSummary(ctx, summaryContent)
	if err != nil {
		return nil, fmt.Errorf("AI总结生成失败: %v", err)
	}

	// 创建新的总结记录
	now := time.Now()
	aiSum := &models.AISum{
		TopicId:    topicId,
		SumTime:    now,
		SumContent: summaryText,
		SumValid:   "Y",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 使用事务保存数据
	err = sqls.DB().Transaction(func(tx *gorm.DB) error {
		// 使之前的总结失效
		if previousSummary != nil {
			if transactionErr := repositories.AISumRepository.InvalidatePreviousSummaries(tx, topicId); transactionErr != nil {
				return transactionErr
			}
		}

		// 保存新的总结
		if transactionErr := repositories.AISumRepository.Create(tx, aiSum); transactionErr != nil {
			return transactionErr
		}

		// 标记评论为已总结
		if transactionErr := repositories.CommentRepository.MarkCommentsAsSummarized(tx, topicId); transactionErr != nil {
			return transactionErr
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("保存总结记录失败: %v", err)
	}

	return aiSum, nil
}

// GetSummary 获取话题的AI总结
func (s *aiSumService) GetSummary(topicId int64) *models.AISum {
	return repositories.AISumRepository.GetValidSummary(sqls.DB(), topicId)
}

// GetSummaryHistory 获取话题的总结历史
func (s *aiSumService) GetSummaryHistory(topicId int64) ([]models.AISum, error) {
	return repositories.AISumRepository.FindByTopicId(sqls.DB(), topicId)
}

// FindPageByParams 分页查询AI总结
func (s *aiSumService) FindPageByParams(params *params.QueryParams) ([]models.AISum, *sqls.Paging) {
	return repositories.AISumRepository.FindPageByParams(sqls.DB(), params)
}

// InvalidateSummary 使总结失效
func (s *aiSumService) InvalidateSummary(id int64) error {
	summary := repositories.AISumRepository.GetByID(sqls.DB(), id)
	if summary == nil {
		return fmt.Errorf("AI总结不存在")
	}

	summary.SumValid = "N"
	return repositories.AISumRepository.Update(sqls.DB(), summary)
}

// Delete 删除AI总结
func (s *aiSumService) Delete(id int64) error {
	return repositories.AISumRepository.Delete(sqls.DB(), id)
}

// buildSummaryContent 构建总结内容
func (s *aiSumService) buildSummaryContent(topic *models.Topic, comments []models.Comment, previousSummary *models.AISum) (string, error) {
	content := fmt.Sprintf("话题标题：%s\n", topic.Title)
	content += fmt.Sprintf("话题内容：%s\n", topic.Content)
	content += "\n需要总结的评论：\n"

	for i, comment := range comments {
		userName := "匿名用户"
		user := UserService.Get(comment.UserId)
		if user != nil {
			userName = user.Nickname
		}

		content += fmt.Sprintf("评论%d（用户：%s）：%s\n", i+1, userName, comment.Content)

		// 如果有引用评论，也包含引用内容
		if comment.QuoteId > 0 {
			quoteComment := CommentService.Get(comment.QuoteId)
			if quoteComment != nil {
				quoteUserName := "匿名用户"
				quoteUser := UserService.Get(quoteComment.UserId)
				if quoteUser != nil {
					quoteUserName = quoteUser.Nickname
				}
				content += fmt.Sprintf("  引用自（用户：%s）：%s\n", quoteUserName, quoteComment.Content)
			}
		}
		content += "\n"
	}

	if previousSummary != nil {
		content += fmt.Sprintf("\n之前的总结（仅供参考）：%s\n", previousSummary.SumContent)
	}

	content += "\n请基于以上内容，生成一个简洁、准确的总结，突出主要观点和讨论重点。"

	return content, nil
}
