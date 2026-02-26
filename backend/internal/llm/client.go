// Package llm 封装 OpenAI 兼容 LLM 调用，仅从环境变量读取 OPENAI_API_KEY、OPENAI_BASE_URL。
package llm

import (
	"context"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// Client 封装 LLM 调用，无 key 时 client 为 nil，调用方需处理不可用。
type Client struct {
	client llms.LLM
}

// ErrNotAvailable 表示未配置 API Key 或 LLM 不可用。
var ErrNotAvailable = &NotAvailableError{}

type NotAvailableError struct{}

func (e *NotAvailableError) Error() string { return "LLM not available: OPENAI_API_KEY not set or init failed" }

// NewClient 创建 LLM 客户端。API Key 仅从环境变量 OPENAI_API_KEY 读取；OPENAI_BASE_URL 可选，默认 https://api.deepseek.com。
func NewClient() *Client {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return &Client{client: nil}
	}
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	os.Setenv("OPENAI_API_KEY", apiKey)
	os.Setenv("OPENAI_BASE_URL", baseURL)
	c, err := openai.New(
		openai.WithModel("deepseek-chat"),
		openai.WithBaseURL(baseURL),
	)
	if err != nil {
		return &Client{client: nil}
	}
	return &Client{client: c}
}

// Generate 根据 prompt 生成文本。若 client 未初始化返回 ErrNotAvailable。
func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	if c.client == nil {
		return "", ErrNotAvailable
	}
	return llms.GenerateFromSinglePrompt(ctx, c.client, prompt)
}
