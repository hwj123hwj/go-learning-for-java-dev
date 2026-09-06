package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAICompat 是 OpenAI 兼容端点的客户端（OpenAI/智谱/DeepSeek/你自己的 LiteLLM 都是这套协议）。
type OpenAICompat struct {
	name    string
	baseURL string // 形如 https://api.example.com/v1
	apiKey  string
	client  *http.Client
}

func NewOpenAICompat(name, baseURL, apiKey string) *OpenAICompat {
	return &OpenAICompat{
		name:    name,
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		// 关键决策：不设 client.Timeout！
		// Timeout 覆盖"整个请求周期包括读 body"，长回复的流式响应会被拦腰砍断。
		// 流式上游的超时拆成两层：
		//   ResponseHeaderTimeout —— 等响应头最多 10s（防上游假死）
		//   整体时长由调用方的 context 控制（见 handler 层）
		client: &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: 10 * time.Second,
				MaxIdleConnsPerHost:   8, // 连接复用：网关的标配，避免每请求新建 TCP
			},
		},
	}
}

func (o *OpenAICompat) Name() string { return o.name }

// newRequest 构造发往上游的 POST。
// ctx 从 handler 一路传下来：客户端断开 → 整条链路自动取消。
func (o *OpenAICompat) newRequest(ctx context.Context, req ChatRequest) (*http.Request, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		o.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if o.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)
	}
	return httpReq, nil
}

// drainBody 读掉并丢弃响应体（保持连接可复用），上限防内存被恶意响应撑爆。
func drainBody(resp *http.Response) {
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 64*1024))
}

// upstreamError 把非 2xx 统一翻译成可重试/不可重试的错误。
func (o *OpenAICompat) upstreamError(resp *http.Response) error {
	drainBody(resp)
	status := resp.StatusCode
	switch {
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		// 401/403：换上游可能有用（另一个 key），但重试自己没用
		return &RetryableError{Provider: o.name, Status: status,
			Err: fmt.Errorf("auth failed")}
	case status >= 500 || status == http.StatusTooManyRequests:
		// 5xx/429：典型的"值得再试"，标记可重试
		return &RetryableError{Provider: o.name, Status: status,
			Err: fmt.Errorf("server error or rate limited")}
	default:
		// 4xx（400 参数错、404 路径错）：请求本身有问题，重试无意义
		return fmt.Errorf("provider %s: HTTP %d: non-retryable", o.name, status)
	}
}

// Complete 非流式调用。
func (o *OpenAICompat) Complete(ctx context.Context, req ChatRequest) (string, error) {
	req.Stream = false
	httpReq, err := o.newRequest(ctx, req)
	if err != nil {
		return "", err
	}

	resp, err := o.client.Do(httpReq)
	if err != nil {
		// 网络层错误（连接拒绝/DNS/超时）：值得换一个上游重试
		return "", &RetryableError{Provider: o.name, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", o.upstreamError(resp)
	}

	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", &RetryableError{Provider: o.name, Err: fmt.Errorf("decode response: %w", err)}
	}
	if len(out.Choices) == 0 {
		return "", &RetryableError{Provider: o.name, Err: fmt.Errorf("empty choices")}
	}
	return out.Choices[0].Message.Content, nil
}

// Stream 流式调用：逐行解析 SSE，把增量文本回调出去。
//
// SSE 线格式（就这么朴素）：
//   data: {"choices":[{"delta":{"content":"你"}}]}
//
//   data: {"choices":[{"delta":{"content":"好"}}]}
//
//   data: [DONE]
func (o *OpenAICompat) Stream(ctx context.Context, req ChatRequest, onDelta func(string)) (string, error) {
	req.Stream = true
	httpReq, err := o.newRequest(ctx, req)
	if err != nil {
		return "", err
	}

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return "", &RetryableError{Provider: o.name, Err: err}
	}
	defer resp.Body.Close()

	// 响应头都拿到了：到这里为止还没给客户端吐任何字节，仍可以安全 fallback。
	if resp.StatusCode >= 400 {
		return "", o.upstreamError(resp)
	}

	// 从这里开始把增量交给回调——一旦吐出第一个字节，就再也回不了头。
	var full strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue // 跳过空行/注释行/event 行
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			break
		}

		var ev struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &ev); err != nil || len(ev.Choices) == 0 {
			continue // 容忍畸形事件：流式协议里什么都可能混进来
		}
		if d := ev.Choices[0].Delta.Content; d != "" {
			full.WriteString(d)
			onDelta(d)
		}
	}
	if err := scanner.Err(); err != nil {
		// 流中途断掉：客户端已经看到一半内容，不能重放、不能 fallback，
		// 只能把错误交给 handler 在流内告知。这就是"首字节规则"。
		return full.String(), fmt.Errorf("provider %s: stream broken after %d bytes: %w",
			o.name, full.Len(), err)
	}
	return full.String(), nil
}
