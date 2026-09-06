// Package todo 是业务层：校验、编排、业务规则。它依赖 storage.Storage
// 接口而不是具体实现——依赖注入让测试可以注入内存实现（见 service_test.go）。
package todo

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"todoapp/internal/storage"
)

// ErrEmptyTitle 是业务层自己的哨兵错误：空标题在校验阶段就该拦下，
// 不应该等存储层报错（分层校验：业务错误在 service，存储错误在 storage）。
var ErrEmptyTitle = errors.New("title must not be empty")

// Service 是业务服务。
// 依赖全部通过构造函数注入（对比 Java 的 @Autowired 字段注入）：
// 显式、可测试、没有隐藏的全局状态。你网关的 handlers 也是这个形状。
type Service struct {
	store storage.Storage
	log   *slog.Logger
	now   func() time.Time // 时间也是依赖：测试里可替换成固定时钟
}

func NewService(store storage.Storage, log *slog.Logger) *Service {
	return &Service{
		store: store,
		log:   log,
		now:   time.Now,
	}
}

func (s *Service) Add(ctx context.Context, title string) (storage.Item, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return storage.Item{}, ErrEmptyTitle
	}
	item := storage.Item{
		ID:        newID(),
		Title:     title,
		UpdatedAt: s.now(),
	}
	if err := s.store.Create(ctx, item); err != nil {
		// 处理不了 → 包装上抛，附带"我当时在干什么"
		return storage.Item{}, fmt.Errorf("add todo: %w", err)
	}
	s.log.InfoContext(ctx, "todo created", "id", item.ID, "title", item.Title)
	return item, nil
}

// Complete 演示业务规则 + 错误链：Get 和 Save 的错误都会带着
// "complete todo" 上下文上抛，根因（ErrNotFound）仍在链上，
// HTTP 层据此映射 404。
func (s *Service) Complete(ctx context.Context, id string) (storage.Item, error) {
	item, err := s.store.Get(ctx, id)
	if err != nil {
		return storage.Item{}, fmt.Errorf("complete todo %q: %w", id, err)
	}
	if item.Done {
		return item, nil // 业务规则：重复完成是幂等的，不报错
	}
	item.Done = true
	item.UpdatedAt = s.now()
	if err := s.store.Save(ctx, item); err != nil {
		return storage.Item{}, fmt.Errorf("complete todo %q: %w", id, err)
	}
	s.log.InfoContext(ctx, "todo completed", "id", id)
	return item, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.store.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete todo %q: %w", id, err)
	}
	s.log.InfoContext(ctx, "todo deleted", "id", id)
	return nil
}

func (s *Service) List(ctx context.Context) ([]storage.Item, error) {
	return s.store.List(ctx)
}

// newID 生成随机 ID。crypto/rand 而不是 math/rand：
// ID 不是密钥但保证全局唯一即可，8 字节够用且无序（不泄露创建顺序）。
func newID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand 失败意味着系统熵源故障，直接 panic 是合理选择
		panic(fmt.Sprintf("crypto/rand failed: %v", err))
	}
	return hex.EncodeToString(b)
}
