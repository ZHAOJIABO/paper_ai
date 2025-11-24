package doubao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"paper_ai/internal/infrastructure/ai/types"
	apperrors "paper_ai/pkg/errors"
	"paper_ai/pkg/logger"
	"go.uber.org/zap"
)

// Client 豆包客户端
type Client struct {
	apiKey  string
	baseURL string
	model   string
	timeout time.Duration
	client  *http.Client
}

// NewClient 创建豆包客户端
func NewClient(apiKey, baseURL, model string, timeout time.Duration) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Polish 实现段落润色
func (c *Client) Polish(ctx context.Context, req *types.PolishRequest) (*types.PolishResponse, error) {
	// 构建prompt
	prompt := c.buildPolishPrompt(req)

	// 调用豆包API
	doubaoResp, err := c.callDoubaoAPI(ctx, prompt)
	if err != nil {
		logger.Error("failed to call doubao api", zap.Error(err))
		return nil, apperrors.NewAIServiceError("failed to call doubao api", err)
	}

	// 构建响应
	return &types.PolishResponse{
		PolishedContent: doubaoResp.Choices[0].Message.Content,
		OriginalLength:  len(req.Content),
		PolishedLength:  len(doubaoResp.Choices[0].Message.Content),
		Suggestions:     c.extractSuggestions(doubaoResp.Choices[0].Message.Content),
		ProviderUsed:    "doubao",
		ModelUsed:       c.model,
	}, nil
}

// buildPolishPrompt 构建润色prompt
func (c *Client) buildPolishPrompt(req *types.PolishRequest) string {
	stylePrompt := ""
	switch req.Style {
	case "academic":
		stylePrompt = "请以学术风格润色以下文本，使其更加正式、准确，适合用于学术论文。"
	case "formal":
		stylePrompt = "请以正式风格润色以下文本，使其更加专业，适合用于正式文档。"
	case "concise":
		stylePrompt = "请润色以下文本使其更加简洁，去除冗余内容同时保持清晰。"
	default:
		stylePrompt = "请润色以下文本以提高其清晰度、连贯性和可读性。"
	}

	languagePrompt := ""
	if req.Language == "zh" {
		languagePrompt = "请确保润色后的文本为中文。"
	} else {
		languagePrompt = "请确保润色后的文本为英文。"
	}

	return fmt.Sprintf(`%s %s

原始文本：
%s

请只返回润色后的文本，不需要任何解释或元数据。`, stylePrompt, languagePrompt, req.Content)
}

// extractSuggestions 提取改进建议（简化版实现）
func (c *Client) extractSuggestions(polishedText string) []string {
	// 这里可以根据实际需求实现更复杂的建议提取逻辑
	// 目前返回空切片，未来可扩展
	return []string{}
}

// DoubaoAPIRequest 豆包API请求结构
type DoubaoAPIRequest struct {
	Model    string           `json:"model"`
	Messages []DoubaoMessage  `json:"messages"`
}

// DoubaoMessage 豆包消息结构
type DoubaoMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DoubaoAPIResponse 豆包API响应结构
type DoubaoAPIResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []DoubaoChoice `json:"choices"`
	Usage   DoubaoUsage    `json:"usage"`
}

// DoubaoChoice 豆包选择结构
type DoubaoChoice struct {
	Index        int           `json:"index"`
	Message      DoubaoMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

// DoubaoUsage 豆包使用统计
type DoubaoUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// DoubaoErrorResponse 豆包错误响应
type DoubaoErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// callDoubaoAPI 调用豆包API
func (c *Client) callDoubaoAPI(ctx context.Context, prompt string) (*DoubaoAPIResponse, error) {
	// 构建请求体
	reqBody := DoubaoAPIRequest{
		Model: c.model,
		Messages: []DoubaoMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 构建HTTP请求
	url := strings.TrimRight(c.baseURL, "/") + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// 发送请求
	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 检查HTTP状态码
	if httpResp.StatusCode != http.StatusOK {
		var errResp DoubaoErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("doubao api error: %s - %s", errResp.Error.Type, errResp.Error.Message)
		}
		return nil, fmt.Errorf("doubao api error: status %d, body: %s", httpResp.StatusCode, string(body))
	}

	// 解析响应
	var doubaoResp DoubaoAPIResponse
	if err := json.Unmarshal(body, &doubaoResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &doubaoResp, nil
}
