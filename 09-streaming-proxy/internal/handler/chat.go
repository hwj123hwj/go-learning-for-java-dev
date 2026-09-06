// Package handler 是 HTTP 接入层。chat.go 处理 /v1/chat/completions——
// 这个模块的核心：非流式转发 + SSE 流式转发。
package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"streaming-proxy/internal/provider"
)

// requestTimeout 是一次请求的整体超时（含流式全程）。
// 它和上游 client 的两层超时配合：Transport 层管"等响应头"，
// 这里管"整个请求最多活多久"。见 provider/openai.go 的注释。
const requestTimeout = 60 * time.Second

type Chat struct {
	chain *provider.Chain
}

func NewChat(chain *provider.Chain) *Chat {
	return &Chat{chain: chain}
}

func (h *Chat) Register(r gin.IRouter) {
	r.POST("/v1/chat/completions", h.Handle)
}

func (h *Chat) Handle(c *gin.Context) {
	var req provider.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid JSON: " + err.Error()}})
		return
	}
	if len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "messages 不能为空"}})
		return
	}

	// ctx 组装：请求 ctx（客户端断开自动取消）+ 整体超时。
	// cancel 放在 Handle 的 defer：Handle 返回时流一定已写完，可以安全释放。
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	if req.Stream {
		h.handleStream(c, ctx, req)
		return
	}
	h.handleComplete(c, ctx, req)
}

// ─── 非流式 ─────────────────────────────────────────────────────────────────

func (h *Chat) handleComplete(c *gin.Context, ctx context.Context, req provider.ChatRequest) {
	text, err := h.chain.Complete(ctx, req)
	if err != nil {
		h.writeUpstreamError(c, err)
		return
	}
	c.JSON(http.StatusOK, chatCompletionResponse{
		ID:      newResponseID(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []choice{{
			Index:        0,
			Message:      responseMessage{Role: "assistant", Content: text},
			FinishReason: "stop",
		}},
	})
}

// ─── 流式（SSE）──────────────────────────────────────────────────────────────

func (h *Chat) handleStream(c *gin.Context, ctx context.Context, req provider.ChatRequest) {
	// SSE 响应头三件套。写头之前还有机会改状态码；
	// 头一旦发出（下面往 c.Writer 写入后），错误只能在流内表达。
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("X-Accel-Buffering", "no") // 告诉 nginx 别缓冲，逐块转发

	flusher, canFlush := c.Writer.(http.Flusher)
	id := newResponseID()
	writeChunk := func(delta, finish string) {
		chunk := chatCompletionChunk{
			ID:      id,
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []chunkChoice{{
				Index:        0,
				Delta:        responseMessage{Role: "assistant", Content: delta},
				FinishReason: finish,
			}},
		}
		raw, _ := json.Marshal(chunk)
		fmt.Fprintf(c.Writer, "data: %s\n\n", raw)
		// 关键：每个 chunk 立刻 Flush。不 Flush 的话数据攒在缓冲区，
		// 客户端到最后才一次收到全部——"流式"就名存实亡了。
		if canFlush {
			flusher.Flush()
		}
	}

	_, err := h.chain.Stream(ctx, req, func(delta string) {
		writeChunk(delta, "")
	})
	if err != nil {
		// 200 头已发出，无法再改状态码——只能在流里发一个错误事件，
		// 客户端要能处理这种"data 里带 error"的情况（OpenAI 同款做法）。
		errPayload, _ := json.Marshal(gin.H{"error": gin.H{"message": err.Error()}})
		fmt.Fprintf(c.Writer, "data: %s\n\n", errPayload)
		if canFlush {
			flusher.Flush()
		}
		return
	}
	writeChunk("", "stop")
	fmt.Fprint(c.Writer, "data: [DONE]\n\n")
	if canFlush {
		flusher.Flush()
	}
}

// writeUpstreamError 把链路错误映射成状态码——非流式时头还没发，可以正常回错误码。
func (h *Chat) writeUpstreamError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": gin.H{"message": "upstream timeout"}})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"message": err.Error()}})
	}
}

// ─── OpenAI wire 类型（只保留教学需要的字段）──────────────────────────────

type responseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type choice struct {
	Index        int             `json:"index"`
	Message      responseMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
}

type chatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []choice `json:"choices"`
}

type chunkChoice struct {
	Index        int             `json:"index"`
	Delta        responseMessage `json:"delta"`
	FinishReason string          `json:"finish_reason"`
}

type chatCompletionChunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []chunkChoice `json:"choices"`
}

func newResponseID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "chatcmpl-" + hex.EncodeToString(b)
}
