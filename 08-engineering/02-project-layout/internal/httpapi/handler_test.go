package httpapi_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"todoapp/internal/httpapi"
	"todoapp/internal/storage"
	"todoapp/internal/todo"
)

// newTestServer 起一个真实的 HTTP 测试服务器（httptest），
// 测试走完整网络栈：路由 → 中间件 → handler → service → 内存存储。
// 对比 Java：MockMvc / @SpringBootTest + TestRestTemplate。
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	svc := todo.NewService(storage.NewMemoryStorage(),
		slog.New(slog.NewTextHandler(io.Discard, nil)))
	h := httpapi.NewHandler(svc,
		slog.New(slog.NewTextHandler(io.Discard, nil)))
	return httptest.NewServer(h.Router())
}

// post helper
func post(t *testing.T, ts *httptest.Server, path string, body any) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	resp, err := http.Post(ts.URL+path, "application/json", bytes.NewReader(b))
	require.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func TestTodoLifecycle(t *testing.T) {
	ts := newTestServer(t)
	t.Cleanup(ts.Close)

	// 1. 创建
	resp := post(t, ts, "/todos", map[string]string{"title": "学 Go 工程化"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created storage.Item
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	require.NotEmpty(t, created.ID)
	assert.Equal(t, "学 Go 工程化", created.Title)

	// 2. 列表里能查到
	listResp, err := http.Get(ts.URL + "/todos")
	require.NoError(t, err)
	t.Cleanup(func() { listResp.Body.Close() })

	var items []storage.Item
	require.NoError(t, json.NewDecoder(listResp.Body).Decode(&items))
	require.Len(t, items, 1)

	// 3. 完成
	doneResp := post(t, ts, "/todos/"+created.ID+"/complete", nil)
	require.Equal(t, http.StatusOK, doneResp.StatusCode)

	var completed storage.Item
	require.NoError(t, json.NewDecoder(doneResp.Body).Decode(&completed))
	assert.True(t, completed.Done)
}

// TestErrorMapping 端到端验证"错误 → 状态码"映射：
// 业务校验错误 400、不存在 404。这是 handler 最该测的部分。
func TestErrorMapping(t *testing.T) {
	ts := newTestServer(t)
	t.Cleanup(ts.Close)

	t.Run("空标题返回 400", func(t *testing.T) {
		resp := post(t, ts, "/todos", map[string]string{"title": ""})
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("畸形 JSON 返回 400", func(t *testing.T) {
		resp, err := http.Post(ts.URL+"/todos", "application/json",
			bytes.NewReader([]byte(`{"title": 123}`)))
		require.NoError(t, err)
		t.Cleanup(func() { resp.Body.Close() })
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("完成不存在的事项返回 404", func(t *testing.T) {
		resp := post(t, ts, "/todos/no-such-id/complete", nil)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("删除不存在的事项返回 404", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, ts.URL+"/todos/no-such-id", nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { resp.Body.Close() })
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
