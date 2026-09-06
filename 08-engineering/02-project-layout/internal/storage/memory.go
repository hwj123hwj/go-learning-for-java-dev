package storage

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// MemoryStorage 是 Storage 的内存实现：测试与本地开发用。
// 换 Postgres 时新写一个 PostgresStorage 放同包，上层代码零改动——
// 这就是依赖接口的意义（对比 Java：换实现只改 @Configuration）。
type MemoryStorage struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{items: make(map[string]Item)}
}

// ctx 参数在内存实现里用不上（没有 IO），保留它是为了契约一致：
// 真实实现里它承担超时取消（context.WithTimeout）和链路追踪。
func (m *MemoryStorage) Create(_ context.Context, item Item) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.items[item.ID]; exists {
		return fmt.Errorf("create item %q: %w", item.ID, ErrConflict)
	}
	m.items[item.ID] = item
	return nil
}

func (m *MemoryStorage) Get(_ context.Context, id string) (Item, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	item, ok := m.items[id]
	if !ok {
		return Item{}, fmt.Errorf("get item %q: %w", id, ErrNotFound)
	}
	return item, nil
}

func (m *MemoryStorage) List(_ context.Context) ([]Item, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Item, 0, len(m.items))
	for _, it := range m.items {
		out = append(out, it)
	}
	// 排序保证输出稳定，测试才能断言（map 遍历顺序是随机的）
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt.Before(out[j].UpdatedAt) })
	return out, nil
}

func (m *MemoryStorage) Save(_ context.Context, item Item) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.items[item.ID]; !ok {
		return fmt.Errorf("save item %q: %w", item.ID, ErrNotFound)
	}
	m.items[item.ID] = item
	return nil
}

func (m *MemoryStorage) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.items[id]; !ok {
		return fmt.Errorf("delete item %q: %w", id, ErrNotFound)
	}
	delete(m.items, id)
	return nil
}
