# 01 · 错误处理：错误是值

Go 没有 `try/catch`。错误是函数最后一个返回值，是个普通的值——这是 Go 和 Java 心智模型差异最大的一处，也是你网关源码里到处都是的模式。

## 跑

```bash
go run ./01-errors
```

## 这一课的四个概念

1. **哨兵错误**（`var ErrNotFound = errors.New(...)`）：包级导出的错误变量，调用方 `errors.Is` 匹配——等价于 Java 按异常类型 catch。
2. **`%w` 包装**：`fmt.Errorf("findUser(id=%d): %w", id, ErrNotFound)`——补充上下文同时保留错误链。`%v` 会把链压断，根因就丢了。
3. **自定义错误类型 + `errors.As`**：错误结构体带业务字段（如 `ValidationError.Field`），取出来直接读——等价于 Java 异常的 getter。
4. **`errors.Join`**：一次返回多个错误。

## 练习（答案见 [../SOLUTIONS.md](../SOLUTIONS.md)）

**1.1〔预测〕** `fmt.Errorf("load: %v", ErrNotFound)` 之后，`errors.Is(err, ErrNotFound)` 是 true 还是 false？为什么？

**1.2〔找 bug〕** 下面代码想区分"用户不存在"和"数据库连不上"，调用方却总是走进 500 分支。问题出在哪？

```go
func GetUser(id string) error {
    if _, err := query(id); err != nil {
        return fmt.Errorf("query user %d: %s", id, err.Error()) // 注意这里
    }
    return nil
}

// 调用方
if err := GetUser(42); err != nil {
    if errors.Is(err, ErrNotFound) {
        return 404 // 永远走不进来
    }
    return 500
}
```

**1.3〔徒手〕** 定义 `RateLimitError{RetryAfter time.Duration}`，实现 `Error()`；写一个会随机返回它或 `nil` 的函数，调用方用 `errors.As` 打印"建议 X 秒后重试"。
