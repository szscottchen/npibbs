package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

// BuildAISum 构建AI总结响应
func BuildAISum(aiSum *models.AISum) *models.AISumResponse {
	if aiSum == nil {
		return nil
	}

	resp := &models.AISumResponse{
		Id:         aiSum.Id,
		TopicId:    aiSum.TopicId,
		SumTime:    aiSum.SumTime,
		SumContent: aiSum.SumContent,
		SumValid:   aiSum.SumValid,
		CreatedAt:  aiSum.CreatedAt,
		UpdatedAt:  aiSum.UpdatedAt,
	}

	// 加载话题信息
	if aiSum.Topic != nil {
		resp.TopicTitle = aiSum.Topic.Title
	} else if aiSum.TopicId > 0 {
		topic := services.TopicService.Get(aiSum.TopicId)
		if topic != nil {
			resp.TopicTitle = topic.Title
		}
	}

	return resp
}

// BuildAISumList 构建AI总结列表响应
func BuildAISumList(aiSums []models.AISum) []models.AISumResponse {
	if len(aiSums) == 0 {
		return nil
	}

	var responses []models.AISumResponse
	for _, aiSum := range aiSums {
		resp := BuildAISum(&aiSum)
		if resp != nil {
			responses = append(responses, *resp)
		}
	}
	return responses
}