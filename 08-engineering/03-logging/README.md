# 03 · slog 结构化日志

`log/slog` 是 Go 1.21+ 的官方结构化日志库，你两个真实项目都在用它。核心心智模型：**日志是一行行的键值对，不是格式化字符串**。

## 跑

```bash
go run ./03-logging
```

## 练习（答案见 [../SOLUTIONS.md](../SOLUTIONS.md)）

**3.1〔改造〕** 把 todoapp（02 课）的 JSON logger 换成 Text handler 但保留字段，观察输出差异；再把 `Level` 改成 `slog.LevelWarn`，Info 日志去哪了？

**3.2〔找 bug〕** 下面两行日志，哪一行在日志平台上检索 `user_id=42` 能命中？为什么？

```go
logger.Info(fmt.Sprintf("user %d logged in", userID))
logger.Info("user logged in", "user_id", userID)
```

**3.3〔思考〕** `logger.With("component", "api")` 和每次传 `"component", "api"` 效果一样吗？性能差在哪？（提示：With 返回的子 logger 预序列化了公共 attrs）
