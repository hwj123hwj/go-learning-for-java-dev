# 02 · 项目布局：一个完整的 todo 服务（核心课）

这是本模块最重要的一课：一个能跑、有测试、会优雅退出的完整服务，用的全是你在 litellm-gateway / pi-go 里见到的形状。

## 目录结构与职责

```
todoapp/
├── cmd/todo/main.go          # 唯一入口：只做装配和启动
├── internal/                 # 编译器强制私有——别的项目无法 import（对比 Java 的包私有）
│   ├── storage/              # 数据层：Item 模型 + Storage 接口 + 内存实现 + 哨兵错误
│   ├── todo/                 # 业务层：校验、编排、业务规则（幂等、状态流转）
│   └── httpapi/              # 接入层：路由、中间件、错误→状态码映射
└── go.mod
```

分层与依赖方向（箭头只能向下）：

```
cmd/main ──> httpapi ──> todo ──> storage(接口)
                                └── storage/memory(实现)
```

## 必须吃透的七个点

1. **`internal/` 的语义**：放进去的包只有本 module 能 import——不是文档约定，是编译器强制的。你网关/pi-go 全是这个布局。
2. **接口定义在使用方**：`Storage` 接口长什么样由业务层（使用方）说了算，内存/Postgres 实现隐式满足。对比 Java"接口+Impl 都在实现方"。
3. **依赖注入 = 构造函数传参**：`NewService(store, logger)`——没有 @Autowired，没有全局单例，测试想换实现就换。
4. **错误链贯穿三层**：storage 哨兵错误（`ErrNotFound`）→ service `%w` 包装加上下文 → httpapi `errors.Is` 映射 404。
5. **中间件是包装函数**：`withLogging(withRecovery(mux))` 洋葱模型，对比 Java Filter/Interceptor；Gin 的中间件是同一个思想的糖。
6. **table-driven 测试**：用例是数据，断言用 testify（`require` 断前提、`assert` 断结果——真实项目的主流搭配）。
7. **优雅退出**：`signal.NotifyContext` + `srv.Shutdown`，这是常驻服务的标准收尾。

## 跑

```bash
go test ./... -v          # 13 个子测试全绿
go run ./cmd/todo         # 起服务后另开终端：
curl -X POST localhost:8090/todos -d '{"title":"学 Go 工程化"}'
curl localhost:8090/todos
curl -X POST localhost:8090/todos/任意ID/complete   # 换成上面返回的 id
curl -X POST localhost:8090/todos/不存在/complete   # 404
curl -X DELETE localhost:8090/todos/任意ID
```

## 练习（答案见 [../SOLUTIONS.md](../SOLUTIONS.md)）

**2.1〔徒手〕** 给 Service 加 `Rename(ctx, id, newTitle)`：空标题返回 `ErrEmptyTitle`、不存在返回 404 可识别的底层错误，HTTP 层加 `PATCH /todos/{id}`。

**2.2〔改造〕** 写一个 `withRequestID` 中间件：每个请求生成随机 8 字节 hex，塞进 `context`，响应头带 `X-Request-ID`，日志中间件打印它——你网关的 request_id 中间件就是这么实现的。

**2.3〔思考〕** 为什么 `Complete` 里 Get 和 Save 之间没有加锁也不会坏（单机内存版）？如果两个请求并发 Complete 同一个 item 会发生什么？换成 PostgresStorage 后这个竞态还在吗？（提示：read-modify-write；真实做法是 UPDATE ... WHERE done=false 或乐观锁版本号）

**2.4〔找 bug〕** `handler.go` 的 `writeJSON` 在 `Encode` 出错时为什么只能记日志、没法改状态码了？（提示：`WriteHeader` 已经发出去了。真实框架怎么缓解？）
