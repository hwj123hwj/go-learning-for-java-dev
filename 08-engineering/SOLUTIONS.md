# 08 工程化 · 参考答案

## 01 错误处理

**1.1** `false`。`%v` 把错误格式化成纯字符串，错误链断了——`errors.Is` 沿着 `%w` 建立的 `Unwrap()` 链找，字符串里没有链。教训：**要保留根因就 `%w`，只在最外层打日志时用 `%v`**。

**1.2** 问题在 `fmt.Errorf("query user %d: %s", id, err.Error())`：用了 `%s` 且取 `.Error()` 字符串，等于把底层错误压扁——`errors.Is` 找不到 `ErrNotFound`，只能掉进 500。改成 `%w` 直接包 err 本身：

```go
return fmt.Errorf("query user %d: %w", id, err)
```

**1.3**

```go
type RateLimitError struct{ RetryAfter time.Duration }

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limited, retry after %v", e.RetryAfter)
}

// 调用方
var rle *RateLimitError
if errors.As(err, &rle) {
	fmt.Printf("建议 %.0f 秒后重试\n", rle.RetryAfter.Seconds())
}
```

注意 `errors.As` 第二个参数是 `**RateLimitError`（目标类型的指针的指针）。

## 02 项目布局

**2.1** Service 层：

```go
func (s *Service) Rename(ctx context.Context, id, newTitle string) (storage.Item, error) {
	newTitle = strings.TrimSpace(newTitle)
	if newTitle == "" {
		return storage.Item{}, ErrEmptyTitle
	}
	item, err := s.store.Get(ctx, id)
	if err != nil {
		return storage.Item{}, fmt.Errorf("rename todo %q: %w", id, err)
	}
	item.Title = newTitle
	item.UpdatedAt = s.now()
	if err := s.store.Save(ctx, item); err != nil {
		return storage.Item{}, fmt.Errorf("rename todo %q: %w", id, err)
	}
	return item, nil
}
```

HTTP 层（注意 PATCH 需要 `http.NewRequest` 手动发，Go 标准库 client 支持 PATCH）：

```go
mux.HandleFunc("PATCH /todos/{id}", h.rename)

func (h *Handler) rename(w http.ResponseWriter, r *http.Request) {
	var req struct{ Title string `json:"title"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	item, err := h.svc.Rename(r.Context(), r.PathValue("id"), req.Title)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
```

错误映射已经在 `writeServiceError` 里了——这就是"错误映射集中一个点"的红利：加接口不用改映射逻辑。

**2.2**

```go
type ctxKeyRequestID struct{} // 非导出结构体类型当 key，防碰撞

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := newRequestID()
		ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFrom(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyRequestID{}).(string); ok {
		return v
	}
	return ""
}
```

要点：context key 用非导出自定义类型（任何其他包都造不出这个 key）；读取封装成函数；中间件顺序 `withLogging(withRequestID(mux))` 让日志也能取到。你网关的 request_id.go 就是这个形状。

**2.3** 内存版不坏是因为 `MemoryStorage` 内部每个方法都有锁，**单次** Get/Save 是原子的。但两个请求并发 Complete 同一个 item 时：都 Get 到 done=false → 都 Save done=true → 各自返回成功，结果一致纯属运气（幂等操作掩盖了竞态）。换成 Postgres 后竞态仍在：两次 SELECT 都读到未完成、两次 UPDATE 都写入。真实做法：原子更新 `UPDATE todos SET done=true WHERE id=? AND done=false`，检查影响行数判断"是否真的由我完成"，或用版本号乐观锁（GORM 的 `Version` 字段）。

**2.4** 因为 `writeJSON` 先调 `WriteHeader(status)` 把状态行发出去了，HTTP 响应是单向流，头发出去了就收不回；之后 Encode 失败（比如客户端断连）只能记日志。缓解：把响应体先 `json.Marshal` 成字节再写——marshal 失败时头还没发，还能改 500；写响应体失败（连接断了）依然只能日志。Gin 的 `c.JSON` 内部就是先 marshal。

## 03 slog

**3.1** Text handler 输出 `level=INFO msg=... user_id=42`（人类可读，本地开发用），字段仍在但不是 JSON；`Level: slog.LevelWarn` 后 Info/Debug 直接被 handler 丢弃——**过滤发生在 handler 层，调用代码不变**，这就是级别开关的意义。

**3.2** 第二行。第一行格式化进 msg 字符串，平台上它是 `msg="user 42 logged in"` 的一部分，只能全文检索；第二行是独立的 `user_id=42` 字段，可以精确过滤和聚合（P99 延迟按接口分组这类统计全靠字段）。

**3.3** 输出等价。性能差异：`With` 返回的子 logger 把公共 attrs **预解析/预格式化**了一次，之后每条日志直接复用；每次手传则每条日志都要重新处理这些键值。高频路径（每请求都打的访问日志）用 With，低频日志无所谓。

## 04 泛型

**4.1**

```go
func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
// 使用：Keys(map[string]int{"a": 1, "b": 2})
```

（Go 1.23+ 标准库已有 `maps.Keys` 返回 iter.Seq，练习价值在于亲手写一遍约束。）

**4.2** `Sum([]bool{...})` 编译失败：`bool` 不在 `Number` 的类型集合（`~int | ~int64 | ~float64`）里。约束换成 `any` 后，`total += v` 那一行编译失败——`any` 没承诺支持 `+` 运算。结论：**约束就是"能对 T 做什么"的合同**，要算术就得声明数字集合。

**4.3** 参考结论（你可以有自己的）：出现 ≥3 处"同样形状的查找/转换"且类型不同才值得泛型；一两处直接 for 循环更直白。网关里像 `Find[T]` 这种小工具对 `[]Provider`、`[]Handler` 复用是划算的；但"为了泛型而泛型"（比如给只有一种类型的结构体加类型参数）只会增加读代码的成本。Go 社区的共识：泛型是补课，不是风格升级。
