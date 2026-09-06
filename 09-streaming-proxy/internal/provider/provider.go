// Package provider 定义"上游"这个概念：类型、接口、可重试错误。
// 你网关的 internal/ 里对应的是各 provider 的适配层。
package provider

import (
	"context"
	"fmt"
)

// Message 是 OpenAI 风格的消息。
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 是 /v1/chat/completions 的请求体（只保留教学需要的字段）。
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// Provider 是一个上游（OpenAI 兼容端点）。
// 小接口 + 具体类型返回值：Chain 依赖它，mock 上游/真实上游都实现它。
type Provider interface {
	Name() string

	// Complete 非流式：一次性返回完整回复。
	Complete(ctx context.Context, req ChatRequest) (string, error)

	// Stream 流式：每收到一段增量文本回调一次 onDelta，返回完整拼接文本。
	Stream(ctx context.Context, req ChatRequest, onDelta func(delta string)) (string, error)
}

// RetryableError 标记"值得重试或切换上游"的错误（网络错误、5xx、429）。
// 对比：400 参数错误换哪个上游都没用，不该标记可重试。
// 用类型 + errors.As 判断，而不是错误字符串。
type RetryableError struct {
	Provider string
	Status   int // 0 表示网络层错误，没有 HTTP 状态码
	Err      error
}

func (e *RetryableError) Error() string {
	if e.Status > 0 {
		return fmt.Sprintf("provider %s: HTTP %d", e.Provider, e.Status)
	}
	return fmt.Sprintf("provider %s: %v", e.Provider, e.Err)
}

// Unwrap 让 errors.Is/As 能穿透到底层错误（如 context.DeadlineExceeded）。
func (e *RetryableError) Unwrap() error { return e.Err }
