// Package storage 定义数据怎么存取。它只有两个角色：
// 数据形状（Item）+ 存储契约（Storage 接口）+ 哨兵错误。
package storage

import (
	"context"
	"errors"
	"time"
)

// Item 是待办事项的数据模型。
// json tag 决定 API 出入参的字段名（snake_case 是常见约定）。
type Item struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 哨兵错误：调用方用 errors.Is 判断，永远不要匹配错误字符串。
var (
	ErrNotFound = errors.New("item not found")
	ErrConflict = errors.New("item already exists")
)

// Storage 是存储契约。
//
// Go 惯用法：接口定义在"使用方"这一侧（这里是被 service 使用的包），
// 小接口（3~5 个方法），实现方隐式满足——不需要 implements 关键字。
// 对比 Java：接口通常由实现方定义并导出（UserRepository 接口 + Impl），
// Go 鼓励"接受接口，返回结构体"，接口按消费者需要裁剪。
type Storage interface {
	Create(ctx context.Context, item Item) error
	Get(ctx context.Context, id string) (Item, error)
	List(ctx context.Context) ([]Item, error)
	Save(ctx context.Context, item Item) error
	Delete(ctx context.Context, id string) error
}
