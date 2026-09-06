package provider

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// Chain 是一条上游链：按顺序尝试，失败自动 fallback。
// 你网关的"多供应商 fallback"核心就是这段逻辑的完整版。
type Chain struct {
	providers []Provider
	attempts  int           // 每个上游的非流式尝试次数
	backoff   time.Duration // 同一上游两次尝试之间的退避
	log       *slog.Logger
}

func NewChain(log *slog.Logger, attempts int, backoff time.Duration, ps ...Provider) *Chain {
	if attempts < 1 {
		attempts = 1
	}
	return &Chain{providers: ps, attempts: attempts, backoff: backoff, log: log}
}

// Complete 非流式：每个上游试 attempts 次（RetryableError 才重试），
// 用尽后切下一个上游。
func (c *Chain) Complete(ctx context.Context, req ChatRequest) (string, error) {
	var lastErr error
	for _, p := range c.providers {
		for attempt := 1; attempt <= c.attempts; attempt++ {
			text, err := p.Complete(ctx, req)
			if err == nil {
				return text, nil
			}
			lastErr = err

			var re *RetryableError
			if !errors.As(err, &re) {
				// 不可重试（如 400 参数错）：换上游也没用，立刻返回
				return "", err
			}
			c.log.Warn("upstream failed, will retry or fallback",
				"provider", p.Name(), "attempt", attempt, "err", err)

			if attempt < c.attempts {
				select {
				case <-time.After(c.backoff):
				case <-ctx.Done(): // 退避也要响应取消，否则客户端跑了还在傻等
					return "", ctx.Err()
				}
			}
		}
	}
	return "", fmt.Errorf("all %d providers failed: %w", len(c.providers), lastErr)
}

// Stream 流式：拿到上游响应头之前可以换；吐出首字节之后只能硬着头皮走完。
// 所以流式不做同上游重试（重放已发出的字节没有意义），直接 fallback 下一个。
func (c *Chain) Stream(ctx context.Context, req ChatRequest, onDelta func(string)) (string, error) {
	var lastErr error
	for _, p := range c.providers {
		text, err := p.Stream(ctx, req, onDelta)
		if err == nil {
			return text, nil
		}
		lastErr = err

		var re *RetryableError
		if !errors.As(err, &re) {
			// 首字节已吐出（或请求本身非法）：不能 fallback，原样上抛
			return "", err
		}
		c.log.Warn("stream upstream failed, fallback to next",
			"provider", p.Name(), "err", err)
	}
	return "", fmt.Errorf("all %d providers failed: %w", len(c.providers), lastErr)
}
