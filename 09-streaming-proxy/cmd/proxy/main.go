// cmd/proxy 是迷你 LLM 网关入口：装配依赖 + 优雅退出（和 08 模块的 todoapp 同一套收尾）。
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

	"streaming-proxy/internal/config"
	"streaming-proxy/internal/provider"
	"streaming-proxy/internal/server"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		log.Error("invalid config", "err", err)
		os.Exit(1)
	}

	// 上游链：配置顺序即 fallback 顺序
	var providers []provider.Provider
	for _, u := range cfg.Upstreams {
		providers = append(providers, provider.NewOpenAICompat(u.Name, u.BaseURL, u.APIKey))
		log.Info("upstream registered", "name", u.Name, "base_url", u.BaseURL)
	}
	chain := provider.NewChain(log, cfg.Attempts, cfg.Backoff, providers...)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           server.New(cfg, chain, log),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("proxy listening", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "err", err)
	}
}
