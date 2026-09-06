// Package server 组装 gin 引擎：中间件 + 路由。
// 你网关的 main.go + middleware 包干的就是这些事。
package server

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"streaming-proxy/internal/config"
	"streaming-proxy/internal/handler"
	"streaming-proxy/internal/provider"
)

func New(cfg *config.Config, chain *provider.Chain, log *slog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // 生产模式：关掉调试路由和彩色输出

	r := gin.New() // 不用 Default()：中间件显式装配，自己装日志而不是 gin 自带的
	r.Use(
		gin.Recovery(), // panic → 500，防单请求打挂进程
		requestID(),    // 每请求一个 ID，贯穿日志
		accessLog(log), // 结构化访问日志
	)
	if cfg.ProxyAPIKey != "" {
		r.Use(bearerAuth(cfg.ProxyAPIKey)) // 可选：客户端鉴权
	}

	handler.NewHealth().Register(r)
	handler.NewChat(chain).Register(r)
	return r
}

// ─── 中间件 ─────────────────────────────────────────────────────────────────

func requestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		b := make([]byte, 8)
		_, _ = rand.Read(b)
		id := hex.EncodeToString(b)
		c.Set("request_id", id)
		c.Header("X-Request-ID", id) // 响应头带回去，客户端报障时直接贴 ID
		c.Next()
	}
}

// accessLog 用 slog 输出结构化访问日志——和你网关的 middleware/logging.go 同构。
func accessLog(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // 先放行，等 handler 处理完再记日志

		log.Info("request",
			"request_id", c.GetString("request_id"),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
		)
	}
}

// bearerAuth 校验 Authorization: Bearer <key>。
func bearerAuth(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const prefix = "Bearer "
		auth := c.GetHeader("Authorization")
		if len(auth) <= len(prefix) || auth[:len(prefix)] != prefix || auth[len(prefix):] != key {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": gin.H{"message": "missing or invalid bearer token"}})
			return
		}
		c.Next()
	}
}
