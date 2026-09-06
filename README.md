# Java 转 Go · 后端学习仓库

带路线图、练习题和实战项目的 Go 后端学习仓库。目标不是"看懂"，而是**吃透并能徒手写出来**——每个模块都按「跑 → 练 → 对答案 → 重写」的循环使用。

## 这个仓库怎么用

1. **跑**：进入模块目录 `go run ./xx_demo`，对照输出理解每一行在干什么
2. **练**：打开 `EXERCISES.md`，遮住源码做题（预测输出 / 改代码 / 徒手实现 / 找 bug）
3. **对答案**：做完再看 `SOLUTIONS.md`，答错的题回到源码找原因
4. **重写**：最高标准——不看源码，把 demo 的核心逻辑自己写一遍

学习路线的"为什么"、技术选型观点和阶段复盘记录在 [JOURNEY.md](./JOURNEY.md)（按日期追加的成长轨迹）。

## 学习路径

### 主线：Go 后端

| 阶段 | 目录 | 内容 | 练习 |
|---|---|---|---|
| 01 | [01-quickstart](./01-quickstart/) | 基础语法 7 讲（变量/函数/结构体/接口/并发/集合） | — |
| 02 | [02-user-api](./02-user-api/) | Gin + GORM 入门：CRUD、分层架构、PostgreSQL | — |
| 03 | [03-user-api-advanced](./03-user-api-advanced/) | 进阶：JWT 认证、中间件、viper 配置、Service 层单测、Docker 部署 | — |
| 04 | [04-grpc-demo](./04-grpc-demo/) | gRPC 基础：proto、四种模式、拦截器、错误处理 | — |
| 05 | [05-grpc-microservice](./05-grpc-microservice/) | gRPC 微服务：多服务互通、服务间调用 | — |
| 06 | [06-concurrency](./06-concurrency/) | 并发十讲：goroutine/泄漏/channel/select/锁/atomic/sync/Context | [EXERCISES](./06-concurrency/EXERCISES.md) |
| 07 | [07-backend-patterns](./07-backend-patterns/) | 后端三大件：分布式锁、限流、缓存穿透/击穿/雪崩 | [EXERCISES](./07-backend-patterns/EXERCISES.md) |
| 08 | [08-engineering](./08-engineering/) | 工程化：错误处理、cmd/internal 布局、slog、泛型、table-driven 测试 | 各课 README 尾部 |
| 09 | [09-streaming-proxy](./09-streaming-proxy/) | 结业实战：迷你 LLM 网关——SSE 流式转发、超时/重试/fallback | 各课 README 尾部 |

### 支线：面试准备

| 目录 | 内容 |
|---|---|
| [interview-prep/data-structures](./interview-prep/data-structures/) | 链表/栈/队列/树/图/排序/查找 |
| [interview-prep/operating-system](./interview-prep/operating-system/) | 进程调度/内存管理/页面置换/死锁/同步 |
| [interview-prep/computer-network](./interview-prep/computer-network/) | TCP/UDP/HTTP/DNS/协议栈 |

## 08、09 两个新模块为什么存在

这两个模块对应真实工作里最常用、但前七个模块没覆盖的能力：

- **08-engineering**：真实 Go 项目（如本工作区的 litellm-gateway、pi-go）都是 `cmd/` + `internal/` 布局、接口拆分、依赖注入、测试文件占三分之一。这一模块把这些工程习惯单独练一遍，其中 `02-project-layout` 是一个完整的 stdlib `net/http` todo 服务。
- **09-streaming-proxy**：LLM 网关的日常 = SSE 流式转发 + `http.Client` 超时/重试/fallback。这个模块做一个能跑的迷你网关（带 mock 上游，离线可玩），做完回头看 litellm-gateway 的源码会顺畅很多。

## 进度清单

- [ ] 01 基础语法：7 个脚本全部跑过，能说出和 Java 的核心差异
- [ ] 02 user-api：CRUD 跑通，理解 Controller/Service/Repository 分层
- [ ] 03 user-api-advanced：JWT 流程能徒手画出来；`go test ./...` 通过；Docker 跑通
- [ ] 04 gRPC：四种通信模式各跑一遍，能说清楚和 REST 的取舍
- [ ] 05 微服务：两个服务互相调用跑通
- [ ] 06 并发：**EXERCISES.md 全部做完**（含 3 道徒手实现题）
- [ ] 07 后端模式：**EXERCISES.md 全部做完**；缓存三连问题 + 方案能脱稿讲
- [ ] 08 工程化：4 课练完；todoapp 的测试全部看懂，`Complete` 的错误链能徒手写
- [ ] 09 结业：不看源码实现"fallback + SSE 转发"的核心逻辑；mock 双上游演示跑通
- [ ] 结业自查：随机挑一个 demo，遮住源码重写，编译运行通过

## 常用命令

```bash
# 单模块运行（每个目录都是独立 Go module）
cd 06-concurrency && go run ./01_goroutine
cd 09-streaming-proxy && go run ./cmd/proxy

# 单模块测试
cd 08-engineering/02-project-layout && go test ./... -v

# 结业实战三终端演示（SSE 流式 + fallback）
go run ./cmd/mock-upstream --addr :9090 --label 上游A     # 终端1
go run ./cmd/mock-upstream --addr :9091 --label 上游B     # 终端2
go run ./cmd/proxy                                         # 终端3
curl -N -X POST localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"demo","messages":[{"role":"user","content":"你好"}],"stream":true}'
# 杀掉 9090 再请求 → 自动 fallback 到上游B
```

## 历史档案

`.trae/specs/` 是 2026-06 学习期间的规划存档，只作历史参考，不再维护；当前路线以本 README 为准。
