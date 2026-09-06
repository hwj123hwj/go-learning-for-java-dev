package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"streaming-proxy/internal/config"
	"streaming-proxy/internal/provider"
	"streaming-proxy/internal/server"
)

// newProxy 起一个完整代理（真实 gin 引擎 + httptest 假上游），
// 请求走：client → gin → chain → httptest 上游，全链路。
func newProxy(t *testing.T, upstream *httptest.Server) http.Handler {
	t.Helper()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	chain := provider.NewChain(log, 1, time.Millisecond,
		provider.NewOpenAICompat("test", upstream.URL+"/v1", ""))
	cfg := &config.Config{Addr: ":0", Upstreams: []config.Upstream{
		{Name: "test", BaseURL: upstream.URL + "/v1"},
	}}
	return server.New(cfg, chain, log)
}

// openAIUpstream 模拟 OpenAI 兼容上游：非流式返回 JSON，流式返回 SSE。
func openAIUpstream(t *testing.T) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Stream bool `json:"stream"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		if !req.Stream {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"choices": []any{map[string]any{
					"message": map[string]string{"role": "assistant", "content": "完整回复"},
				}},
			})
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		for _, word := range []string{"你", "好", "世", "界"} {
			fmt.Fprintf(w, "data: {\"choices\":[{\"delta\":{\"content\":%q}}]}\n\n", word)
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	t.Cleanup(ts.Close)
	return ts
}

func TestChatNonStream(t *testing.T) {
	proxy := newProxy(t, openAIUpstream(t))

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions",
		bytes.NewReader([]byte(`{"model":"gpt-demo","messages":[{"role":"user","content":"hi"}]}`)))
	rec := httptest.NewRecorder()
	proxy.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Len(t, resp.Choices, 1)
	assert.Equal(t, "完整回复", resp.Choices[0].Message.Content)
}

// TestChatStream 端到端验证 SSE：多 chunk + finish + [DONE]，且 Content-Type 正确。
func TestChatStream(t *testing.T) {
	proxy := newProxy(t, openAIUpstream(t))

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions",
		bytes.NewReader([]byte(`{"model":"gpt-demo","stream":true,"messages":[{"role":"user","content":"hi"}]}`)))
	rec := httptest.NewRecorder()
	proxy.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "text/event-stream")

	body := rec.Body.String()
	// 4 个增量 chunk + 1 个 finish chunk + DONE
	assert.Equal(t, 5, strings.Count(body, "data: {"), "应有 5 个 JSON 事件")
	assert.Contains(t, body, "[DONE]")
	assert.Contains(t, body, `"content":"你"`)
	assert.Contains(t, body, `"finish_reason":"stop"`)
}

func TestChatValidation(t *testing.T) {
	proxy := newProxy(t, openAIUpstream(t))

	t.Run("空 messages 返回 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions",
			bytes.NewReader([]byte(`{"model":"m"}`)))
		rec := httptest.NewRecorder()
		proxy.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("畸形 JSON 返回 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions",
			bytes.NewReader([]byte(`{bad json`)))
		rec := httptest.NewRecorder()
		proxy.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// TestAllUpstreamsDown：所有上游挂 → 非流式 502。
func TestAllUpstreamsDown(t *testing.T) {
	dead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusInternalServerError)
	}))
	t.Cleanup(dead.Close)
	proxy := newProxy(t, dead)

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions",
		bytes.NewReader([]byte(`{"model":"m","messages":[{"role":"user","content":"hi"}]}`)))
	rec := httptest.NewRecorder()
	proxy.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadGateway, rec.Code)
}
