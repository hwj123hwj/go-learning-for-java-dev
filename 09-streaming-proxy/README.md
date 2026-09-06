# 09 · 结业实战：迷你 LLM 网关（SSE 流式 + fallback）

做一个能跑的 LLM API 网关缩影：接收 OpenAI 格式请求，转发给上游（可配多个），支持 **SSE 流式透传**、**超时分层**、**重试 + 多上游 fallback**。做完这个模块，回头看 litellm-gateway 的 `go-gateway/internal/handlers/chat_completions.go` 和 `streaming.go`，你会发现每一块都似曾相识。

## 架构

```
client ──POST /v1/chat/completions──▶ proxy
                                       │ gin 中间件: recovery / requestID / accessLog / (bearerAuth)
                                       │ handler: 参数校验 + ctx 超时
                                       │ chain: 重试 + fallback（首字节规则）
                                       ▼
                              OpenAI 兼容上游 ×N（mock-upstream 可离线模拟）
```

```
streaming-proxy/
├── cmd/proxy/main.go            # 入口：装配 + 优雅退出
├── cmd/mock-upstream/main.go    # 离线假上游：可故意失败，演示 fallback
└── internal/
    ├── config/                  # 环境变量配置（12-Factor）
    ├── provider/
    │   ├── provider.go          # Provider 接口 + RetryableError（错误的"可重试"语义）
    │   ├── openai.go            # OpenAI 兼容客户端：非流式 + SSE 解析 ★
    │   └── chain.go             # 重试 + fallback 链 ★
    ├── handler/chat.go          # /v1/chat/completions：非流式 + SSE 转发 ★
    └── server/server.go         # gin 装配 + 中间件
```

## 三终端演示（不需要任何真实 key）

```bash
# 终端1：上游A（故意先失败2次）
go run ./cmd/mock-upstream --addr :9090 --label 上游A --fail 2
# 终端2：上游B
go run ./cmd/mock-upstream --addr :9091 --label 上游B
# 终端3：代理，上游A在前B在后
UPSTREAMS="A=http://localhost:9090/v1,B=http://localhost:9091/v1" \
ATTEMPTS=2 BACKOFF_MS=300 go run ./cmd/proxy
```

```bash
# 流式（-N 禁用 curl 缓冲，肉眼看到逐字吐出）
curl -N -X POST localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"demo","stream":true,"messages":[{"role":"user","content":"你好"}]}'

# 非流式
curl -s -X POST localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"demo","messages":[{"role":"user","content":"你好"}]}'

# 观察日志：前2次请求 A 失败(500) → 重试 → fallback 到 B，回复带 [上游B] 标识
```

## 必须吃透的五个点

1. **流式不能设 `http.Client.Timeout`**——它管到"body 读完"，长流会被砍断。分层方案：`Transport.ResponseHeaderTimeout` 管等头，整体时长交给 `context.WithTimeout`（见 `provider/openai.go` 和 `handler/chat.go` 的注释）。
2. **首字节规则**：SSE 一旦吐出第一个 chunk，就不能再 fallback/重试（重放会让客户端看到内容跳变）。所以 chain 的 Stream 在拿到上游 200 前可切换，之后只能把错误塞进流里。
3. **每个 chunk 都要 `Flush()`**——不 Flush 就不是流式，是"最后一次性吐出"。`X-Accel-Buffering: no` 是给 nginx 的。
4. **错误的可重试语义用类型表达**：`RetryableError`（网络错/5xx/429）才重试和切换；400 这种换谁都没用。判断用 `errors.As`，不是字符串。
5. **ctx 贯穿全链**：客户端断开 → `c.Request.Context()` 取消 → 上游请求自动中止，不浪费上游配额。

## 测试

```bash
go test ./... -v    # fallback 矩阵、流式端到端、错误映射、网络错误切换
```

测试不用真实网络：上游全部是 `httptest.Server`，毫秒级退避——这也是 table-driven + httptest 的示范。

## 练习（做不出来回 08 模块补课）

**9.1〔徒手〕** 给 `mock-upstream` 加 `--max-conns` 场景不够刺激？直接给 proxy 加**按 API key 限流**中间件（每 key 10 QPS，复用 `07-backend-patterns` 的 `x/time/rate` + `sync.Map` 想法）。

**9.2〔改造〕** `handler` 的非流式路径目前是"解码再重编码"。改成**透传原始 body**（`io.Copy` 上游响应到客户端），想想什么时候透传更好（性能/保真），什么时候需要解码（要改写内容/统计 token）。

**9.3〔思考〕** 如果上游中途断流，客户端已经看到半截回答。设计一下：怎么让客户端知道"这次回答不完整"？（思路：错误事件里带 `partial: true` 和已收到的字节数；或在上游响应头里带可校验的长度）。

**9.4〔结业〕** 不看源码，徒手实现：`Chain.Stream` 的 fallback 逻辑 + handler 的 SSE 转发循环。写完跑 `go test ./internal/provider/...` 对照行为。

## 对照真实网关

学完后打开 `/Users/weijian/Desktop/hwj/litellm-gateway/go-gateway/`，找到并读懂这几个文件（每个都比本模块的对应物复杂一档）：

| 本模块 | 真实网关 |
|---|---|
| `internal/provider/chain.go` | `internal/handlers/chat_completions.go`（fallback 编排） |
| `internal/provider/openai.go` 的 Stream | `internal/handlers/streaming.go` / `responses_streaming.go` |
| `internal/server/server.go` 的 requestID/日志 | `internal/middleware/request_id.go` / `logging.go` |
| `internal/handler/chat.go` 的错误映射 | `internal/handlers/errors.go` |
