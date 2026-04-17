package ai

import (
	"bbs-go/internal/pkg/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var AIService = newAIService()

type aiService struct{}

func newAIService() *aiService {
	return &aiService{}
}

// GenerateSummary 生成AI总结
func (s *aiService) GenerateSummary(ctx context.Context, content string) (string, error) {
	aiConfig := config.Instance.AI

	if !aiConfig.Enabled {
		return "", fmt.Errorf("AI服务未启用")
	}

	switch aiConfig.Provider {
	case "openai":
		return s.callOpenAI(ctx, content, aiConfig)
	case "azure":
		return s.callAzureOpenAI(ctx, content, aiConfig)
	case "deepseek":
		return s.callDeepSeek(ctx, content, aiConfig)
	default:
		return "", fmt.Errorf("不支持的AI服务提供商: %s", aiConfig.Provider)
	}
}

// callOpenAI 调用OpenAI API
func (s *aiService) callOpenAI(ctx context.Context, content string, aiConfig config.AIConfig) (string, error) {
	requestBody := map[string]interface{}{
		"model":       aiConfig.Model,
		"messages":    []map[string]string{{"role": "user", "content": content}},
		"max_tokens":  1000,
		"temperature": 0.7,
	}

	return s.makeAIRequest(ctx, "https://api.openai.com/v1/chat/completions", aiConfig.APIKey, requestBody)
}

// callAzureOpenAI 调用Azure OpenAI API
func (s *aiService) callAzureOpenAI(ctx context.Context, content string, aiConfig config.AIConfig) (string, error) {
	requestBody := map[string]interface{}{
		"messages":    []map[string]string{{"role": "user", "content": content}},
		"max_tokens":  1000,
		"temperature": 0.7,
	}

	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s",
		aiConfig.Endpoint, aiConfig.Model, aiConfig.APIVersion)

	return s.makeAIRequest(ctx, url, aiConfig.APIKey, requestBody)
}

// callDeepSeek 调用DeepSeek API
func (s *aiService) callDeepSeek(ctx context.Context, content string, aiConfig config.AIConfig) (string, error) {
	requestBody := map[string]interface{}{
		"model":       aiConfig.Model,
		"messages":    []map[string]string{{"role": "user", "content": content}},
		"max_tokens":  1000,
		"temperature": 0.7,
	}

	return s.makeAIRequest(ctx, "https://api.deepseek.com/v1/chat/completions", aiConfig.APIKey, requestBody)
}

// makeAIRequest 通用的AI API请求
func (s *aiService) makeAIRequest(ctx context.Context, url, apiKey string, requestBody map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API返回错误: %s, 响应: %s", resp.Status, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("无效的API响应格式")
	}

	choice := choices[0].(map[string]interface{})
	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("无效的消息格式")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("无效的内容格式")
	}

	return strings.TrimSpace(content), nil
}
