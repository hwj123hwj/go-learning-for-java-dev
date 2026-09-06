# 08 · 工程化：从 demo 到真实项目

前七个模块都是 `package main` 单文件。但真实 Go 项目（本工作区的 litellm-gateway、pi-go 就是活例子）长这样：`cmd/` + `internal/` 布局、小接口 + 依赖注入、错误是值、slog 结构化日志、测试文件占三分之一。本模块把这层"工程习惯"单独练一遍。

## 课程

| 课 | 目录 | 内容 | 练习 |
|---|---|---|---|
| 01 | [01-errors](./01-errors/) | 错误是值：哨兵错误、`%w` 包装、`errors.Is/As`、`errors.Join` | README 尾部 |
| 02 | [02-project-layout](./02-project-layout/) | **核心课**：完整 todo 服务——`cmd/`+`internal/` 布局、接口与依赖注入、stdlib `net/http`、table-driven 测试（testify）、优雅退出 | README 尾部 |
| 03 | [03-logging](./03-logging/) | `log/slog` 结构化日志：JSON、级别、子 logger、Context | README 尾部 |
| 04 | [04-generics](./04-generics/) | 泛型：类型约束、工具函数、泛型容器、什么时候不要用 | README 尾部 |

## 运行

```bash
# 01/03/04：本目录是独立 module
go run ./01-errors
go run ./03-logging
go run ./04-generics

# 02 是嵌套的独立 module（为了练真实的 import path）
cd 02-project-layout && go test ./...
cd 02-project-layout && go run ./cmd/todo
curl -X POST localhost:8090/todos -d '{"title":"学 Go 工程化"}'
curl localhost:8090/todos
```

## 为什么 02 用标准库而不是 Gin

你已经在 `02-user-api` 用过 Gin 了。先用裸 `net/http` 写一个完整服务（路由、中间件、路径参数都是 Go 1.22+ 原生能力），再回头看 Gin——它只是把这些包了一层。先会裸的，框架才不是黑盒；你网关里所有 handler 的底层就是这套。

## 答案

四课练习的答案统一在 [SOLUTIONS.md](./SOLUTIONS.md)。
