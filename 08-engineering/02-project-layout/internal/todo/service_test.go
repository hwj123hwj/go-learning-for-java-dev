package todo_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"todoapp/internal/storage"
	"todoapp/internal/todo"
)

// newTestService：每个用例一套全新的 Service + 内存存储，互不污染。
// logger 输出丢弃（io.Discard），测试日志安静。
func newTestService(t *testing.T) *todo.Service {
	t.Helper() // 标记这是测试辅助函数，报错时定位到调用行
	return todo.NewService(storage.NewMemoryStorage(),
		slog.New(slog.NewTextHandler(io.Discard, nil)))
}

// TestAdd 是 table-driven 测试：用例是数据，测试逻辑只写一遍。
// 对比 Java：@ParameterizedTest + @CsvSource。Go 用切片 + t.Run 子测试。
func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error // nil 表示期望成功
	}{
		{name: "正常创建", title: "学 Go 工程化", wantErr: nil},
		{name: "空标题被拦截", title: "", wantErr: todo.ErrEmptyTitle},
		{name: "纯空格被拦截", title: "   ", wantErr: todo.ErrEmptyTitle},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t)

			item, err := svc.Add(context.Background(), tt.title)

			if tt.wantErr != nil {
				// require：失败立即终止本用例（后面代码依赖前提）
				require.Error(t, err)
				// assert：判断错误身份用 errors.Is，绝不匹配字符串
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, item.ID, "应该生成 ID")
			assert.Equal(t, tt.title, item.Title)
			assert.False(t, item.Done, "新事项应该是未完成")
		})
	}
}

func TestComplete(t *testing.T) {
	t.Run("完成未完成的事项", func(t *testing.T) {
		svc := newTestService(t)
		item, err := svc.Add(context.Background(), "买牛奶")
		require.NoError(t, err)

		got, err := svc.Complete(context.Background(), item.ID)

		require.NoError(t, err)
		assert.True(t, got.Done)
	})

	t.Run("重复完成是幂等的", func(t *testing.T) {
		svc := newTestService(t)
		item, _ := svc.Add(context.Background(), "买牛奶")
		_, err := svc.Complete(context.Background(), item.ID)
		require.NoError(t, err)

		got, err := svc.Complete(context.Background(), item.ID)

		require.NoError(t, err, "第二次完成不应报错")
		assert.True(t, got.Done)
	})

	t.Run("不存在的事项返回可识别的底层错误", func(t *testing.T) {
		svc := newTestService(t)

		_, err := svc.Complete(context.Background(), "no-such-id")

		require.Error(t, err)
		// 错误链：service 包装过，但根因仍可被 errors.Is 识别——
		// 这就是 HTTP 层能映射 404 的原因
		assert.ErrorIs(t, err, storage.ErrNotFound)
	})
}

// TestListOrder 验证 List 返回按时间排序（依赖存储实现的稳定顺序）。
func TestListOrder(t *testing.T) {
	svc := newTestService(t)
	for _, title := range []string{"第一件", "第二件", "第三件"} {
		_, err := svc.Add(context.Background(), title)
		require.NoError(t, err)
		time.Sleep(1 * time.Millisecond) // 保证 UpdatedAt 严格递增
	}

	items, err := svc.List(context.Background())

	require.NoError(t, err)
	require.Len(t, items, 3)
	assert.Equal(t, "第一件", items[0].Title)
	assert.Equal(t, "第三件", items[2].Title)
}
