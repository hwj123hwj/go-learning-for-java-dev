package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	fmt.Println("=== slog 结构化日志 ===")
	fmt.Println()

	slog.Info("默认 logger（log 包兼容格式）")
	fmt.Println()

	// 结构化 JSON logger：生产标配，日志平台按字段检索/聚合
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // 级别开关在这里，改一行就能调 verbose
	}))
	slog.SetDefault(logger) // 全局默认（第三方库的 slog 调用也会走它）

	logger.Info("user logged in",
		"user_id", 42,
		"login_method", "password",
		"latency_ms", 128,
	)

	// 子 logger：公共字段挂在 logger 上，不必每行重复传
	// 对比 Java：MDC.put("component", "api") 之后 log.info(...)
	apiLog := logger.With("component", "api", "region", "cn")
	apiLog.Info("request handled")
	apiLog.Warn("slow query", "took_ms", 900)
	apiLog.Debug("cache hit", "key", "user:42") // 只有 Level>=Debug 才输出

	// 带错误的结构化日志：err 作为字段传值，保持可检索
	// apiLog.Error("call upstream failed", "err", err)

	// InfoContext：把 ctx 传下去——接入 trace 后自动携带 trace_id
	ctx := context.Background()
	apiLog.InfoContext(ctx, "with context")

	fmt.Println()
	compareWithJava()
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java SLF4J（字符串插值 + MDC）:")
	fmt.Println(`  MDC.put("component", "api");
  log.info("user {} logged in via {}", userId, method);`)
	fmt.Println()
	fmt.Println("Go slog（键值对，天然结构化）:")
	fmt.Println(`  apiLog := logger.With("component", "api")
  apiLog.Info("user logged in", "user_id", id, "method", method)`)
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - 键值对优先于格式化字符串：日志平台能按 user_id 精确检索")
	fmt.Println("  - error 一律 \"err\", err 传值，不 fmt 成字符串")
	fmt.Println("  - 公共字段用 With 挂子 logger（对应 MDC），别每行重复")
	fmt.Println("  - 生产 JSON + Info，开发 Text + Debug，改 HandlerOptions 即可")
	fmt.Println("  - 敏感信息（密码/token）不打日志——网关的 redactor 就干这个")
}
