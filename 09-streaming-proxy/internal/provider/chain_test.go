package provider_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"streaming-proxy/internal/provider"
)

func quietLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// fakeUpstream 起一个行为可控的 OpenAI 兼容上游。
// status != 0 时返回该状态码；否则返回 content。
type fakeUpstream struct {
	status  int
	content string
	hits    atomic.Int64
}

func (f *fakeUpstream) server(t *testing.T) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.hits.Add(1)
		if f.status != 0 {
			http.Error(w, "planned failure", f.status)
			return
		}
		fmt.Fprintf(w, `{"choices":[{"message":{"content":%q}}]}`, f.content)
	}))
	t.Cleanup(ts.Close)
	return ts
}

func newProvider(ts *httptest.Server, name string) provider.Provider {
	return provider.NewOpenAICompat(name, ts.URL+"/v1", "")
}

func request() provider.ChatRequest {
	return provider.ChatRequest{
		Model:    "test",
		Messages: []provider.Message{{Role: "user", Content: "hi"}},
	}
}

// TestCompleteFallback 是本模块的核心行为：第一个上游挂了自动切第二个。
// table-driven：四种场景共用一份断言逻辑。
func TestCompleteFallback(t *testing.T) {
	tests := []struct {
		name        string
		firstStatus int // 0 = 正常
		lastStatus  int
		wantText    string
		wantErr     bool
		wantFirst   int64 // 第一个上游被请求次数
		wantLast    int64
	}{
		{
			name: "第一个挂(500)，fallback到第二个",
			firstStatus: 500, lastStatus: 0,
			wantText: "from-second", wantFirst: 2, wantLast: 1, // attempts=2：同上游先重试一次再切换
		},
		{
			name: "第一个429限流，fallback到第二个",
			firstStatus: 429, lastStatus: 0,
			wantText: "from-second", wantFirst: 2, wantLast: 1,
		},
		{
			name: "都挂，报全部失败",
			firstStatus: 500, lastStatus: 500,
			wantErr: true, wantFirst: 2, wantLast: 2,
		},
		{
			name: "第一个返回400(不可重试)，立即返回不碰第二上游",
			firstStatus: 400, lastStatus: 0,
			wantErr: true, wantFirst: 1, wantLast: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			first := &fakeUpstream{status: tt.firstStatus, content: "from-first"}
			last := &fakeUpstream{status: tt.lastStatus, content: "from-second"}
			firstTS := first.server(t)
			lastTS := last.server(t)

			chain := provider.NewChain(quietLogger(), 2, time.Millisecond,
				newProvider(firstTS, "first"), newProvider(lastTS, "last"))

			text, err := chain.Complete(context.Background(), request())

			if tt.wantErr {
				require.Error(t, err)
				if tt.firstStatus == 500 && tt.lastStatus == 500 {
					assert.Contains(t, err.Error(), "all 2 providers failed")
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantText, text)
			}
			assert.Equal(t, tt.wantFirst, first.hits.Load(), "第一个上游命中次数")
			assert.Equal(t, tt.wantLast, last.hits.Load(), "第二个上游命中次数")
		})
	}
}

// TestNetworkErrorFallback：连接拒绝也是可重试错误（status=0 的 RetryableError）。
func TestNetworkErrorFallback(t *testing.T) {
	// 起一个马上关掉的服务器，留下一个必然连不上的端口
	dead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	deadURL := dead.URL
	dead.Close()

	alive := &fakeUpstream{content: "alive"}
	aliveTS := alive.server(t)

	chain := provider.NewChain(quietLogger(), 1, time.Millisecond,
		provider.NewOpenAICompat("dead", deadURL+"/v1", ""),
		newProvider(aliveTS, "alive"))

	text, err := chain.Complete(context.Background(), request())

	require.NoError(t, err)
	assert.Equal(t, "alive", text)
	assert.Equal(t, int64(1), alive.hits.Load())
}

// TestStreamFallback：流式同样支持 fallback——上游在吐首字节前失败即可切换。
func TestStreamFallback(t *testing.T) {
	broken := &fakeUpstream{status: 500}
	brokenTS := broken.server(t)

	streaming := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"流\"}}]}\n\n")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"式\"}}]}\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	t.Cleanup(streaming.Close)

	chain := provider.NewChain(quietLogger(), 2, time.Millisecond,
		newProvider(brokenTS, "broken"), newProvider(streaming, "streaming"))

	var deltas []string
	full, err := chain.Stream(context.Background(), request(), func(d string) {
		deltas = append(deltas, d)
	})

	require.NoError(t, err)
	assert.Equal(t, "流式", full)
	assert.Equal(t, []string{"流", "式"}, deltas)
}
