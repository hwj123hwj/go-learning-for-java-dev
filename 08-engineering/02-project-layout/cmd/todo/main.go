// cmd/todo 是程序唯一入口：组装所有依赖并启动 HTTP 服务。
// main 里只有"装配"，没有业务逻辑——依赖全部显式传入，
// 对比 Java：@SpringBootApplication + 组件扫描自动装配。
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todoapp/internal/httpapi"
	"todoapp/internal/storage"
	"todoapp/internal/todo"
)

func main() {
	// 依赖装配（真实项目里 config 也从这里读）：
	// logger → storage → service → handler → server，一条线没有隐藏依赖
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := storage.NewMemoryStorage()
	svc := todo.NewService(store, logger)
	handler := httpapi.NewHandler(svc, logger)

	srv := &http.Server{
		Addr:              ":8090",
		Handler:           handler.Router(),
		ReadHeaderTimeout: 5 * time.Second, // 防 slowloris 攻击的最低配置
	}

	// 优雅退出：Ctrl+C / SIGTERM（docker stop / k8s 滚动更新发的就是这个）
	// 先停止接新请求，给存量请求 5 秒处理完，再退出。
	// 你网关作为 LaunchAgent 常驻服务，靠的就是这段。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done() // 阻塞在这里，直到收到退出信号
	logger.Info("shutting down gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
	}
	logger.Info("bye")
}
