// Package api 提供 AI 占卜 HTTP 处理。
package api

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"fortune-teller-app/internal/fortune"
	"fortune-teller-app/internal/llm"
)

const fortuneAITimeout = 30 * time.Second

var llmClient = llm.NewClient()

// DrawFortuneAI 处理 POST /api/fortune/ai，请求体可选：{"category":"love"}，返回 AI 生成的占卜项或 503。
func DrawFortuneAI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Category string `json:"category"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	cat := strings.TrimSpace(req.Category)
	if cat == "" {
		cat = string(fortune.CategoryGeneral)
	}
	c := fortune.Category(cat)
	if c != fortune.CategoryLove && c != fortune.CategoryCareer &&
		c != fortune.CategoryHealth && c != fortune.CategoryWealth &&
		c != fortune.CategoryGeneral {
		c = fortune.CategoryGeneral
	}
	prompt := buildFortunePrompt(c)
	ctx, cancel := context.WithTimeout(r.Context(), fortuneAITimeout)
	defer cancel()
	raw, err := llmClient.Generate(ctx, prompt)
	if err != nil {
		if err == llm.ErrNotAvailable {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "AI 占卜暂不可用（未配置 OPENAI_API_KEY）"})
			return
		}
		// 401/认证类错误统一返回 503，便于前端展示友好提示
		msg := err.Error()
		if strings.Contains(msg, "401") || strings.Contains(msg, "Authentication") || strings.Contains(msg, "invalid") {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "AI 服务认证失败，请检查 OPENAI_API_KEY 是否有效"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
		return
	}
	item, err := parseFortuneResponse(raw, c)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "AI 返回格式异常，请稍后重试或使用抽签占卜"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(item)
}

func buildFortunePrompt(cat fortune.Category) string {
	label := cat.Label()
	return `你是一位简洁的占卜师。请为「` + label + `」类别生成一条占卜，只输出一行 JSON，不要其他文字。格式：{"title":"四字标题","content":"一两句运势说明","hint":"一句提示"}。`
}

// parseFortuneResponse 从 LLM 输出中提取 JSON 并解析为 fortune.Item。
func parseFortuneResponse(raw string, cat fortune.Category) (fortune.Item, error) {
	raw = strings.TrimSpace(raw)
	re := regexp.MustCompile(`\{[^{}]*\}`)
	if m := re.FindString(raw); m != "" {
		raw = m
	}
	var v struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Hint    string `json:"hint"`
	}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return fortune.Item{}, err
	}
	item := fortune.Item{
		Category: cat,
		Title:    strings.TrimSpace(v.Title),
		Content:  strings.TrimSpace(v.Content),
		Hint:     strings.TrimSpace(v.Hint),
	}
	if item.Title == "" {
		item.Title = "心诚则灵"
	}
	if item.Content == "" {
		item.Content = "顺其自然，静待时机。"
	}
	return item, nil
}
