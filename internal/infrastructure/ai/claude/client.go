package claude

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

// Client Claude客户端
type Client struct {
	apiKey  string
	baseURL string
	model   string
	timeout time.Duration
	client  *http.Client
}

// NewClient 创建Claude客户端
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

	// 调用Claude API
	claudeResp, err := c.callClaudeAPI(ctx, prompt)
	if err != nil {
		logger.Error("failed to call claude api", zap.Error(err))
		return nil, apperrors.NewAIServiceError("failed to call claude api", err)
	}

	// 构建响应
	return &types.PolishResponse{
		PolishedContent: claudeResp.Content[0].Text,
		OriginalLength:  len(req.Content),
		PolishedLength:  len(claudeResp.Content[0].Text),
		Suggestions:     c.extractSuggestions(claudeResp.Content[0].Text),
		ProviderUsed:    "claude",
		ModelUsed:       c.model,
	}, nil
}

// buildPolishPrompt 构建润色prompt
func (c *Client) buildPolishPrompt(req *types.PolishRequest) string {
	stylePrompt := ""
	switch req.Style {
	case "academic":
		stylePrompt = "Please polish the following text in an academic style, making it more formal, precise, and suitable for academic papers."
	case "formal":
		stylePrompt = "Please polish the following text in a formal style, making it more professional and appropriate for formal documents."
	case "concise":
		stylePrompt = "Please polish the following text to be more concise, removing redundancy while maintaining clarity."
	default:
		stylePrompt = "Please polish the following text to improve its clarity, coherence, and readability."
	}

	languagePrompt := ""
	if req.Language == "zh" {
		languagePrompt = "Please ensure the polished text is in Chinese."
	} else {
		languagePrompt = "Please ensure the polished text is in English."
	}

	return fmt.Sprintf(`%s %s

Original text:
%s

Please return only the polished text without any explanations or metadata.`, stylePrompt, languagePrompt, req.Content)
}

// extractSuggestions 提取改进建议（简化版实现）
func (c *Client) extractSuggestions(polishedText string) []string {
	// 这里可以根据实际需求实现更复杂的建议提取逻辑
	// 目前返回空切片，未来可扩展
	return []string{}
}

// ClaudeAPIRequest Claude API请求结构
type ClaudeAPIRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
}

// ClaudeMessage Claude消息结构
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeAPIResponse Claude API响应结构
type ClaudeAPIResponse struct {
	ID      string               `json:"id"`
	Type    string               `json:"type"`
	Role    string               `json:"role"`
	Content []ClaudeContentBlock `json:"content"`
	Model   string               `json:"model"`
}

// ClaudeContentBlock Claude内容块
type ClaudeContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ClaudeErrorResponse Claude错误响应
type ClaudeErrorResponse struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// callClaudeAPI 调用Claude API
func (c *Client) callClaudeAPI(ctx context.Context, prompt string) (*ClaudeAPIResponse, error) {
	// 构建请求体
	reqBody := ClaudeAPIRequest{
		Model:     c.model,
		MaxTokens: 4096,
		Messages: []ClaudeMessage{
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
	url := strings.TrimRight(c.baseURL, "/") + "/v1/messages"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

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
		var errResp ClaudeErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("claude api error: %s - %s", errResp.Error.Type, errResp.Error.Message)
		}
		return nil, fmt.Errorf("claude api error: status %d, body: %s", httpResp.StatusCode, string(body))
	}

	// 解析响应
	var claudeResp ClaudeAPIResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &claudeResp, nil
}
