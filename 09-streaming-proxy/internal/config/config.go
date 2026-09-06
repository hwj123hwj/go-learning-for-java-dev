// Package config 从环境变量读配置（12-Factor：配置注入，不写死）。
// 对比 Java 的 application.yml + @Value——Go 里 os.Getenv + 显式默认值就够了。
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Upstream struct {
	Name    string // 展示名，出现在日志和错误里
	BaseURL string // OpenAI 兼容根地址，含 /v1
	APIKey  string // 可空
}

type Config struct {
	Addr        string        // 监听地址
	Upstreams   []Upstream    // 至少一个
	Attempts    int           // 每上游非流式尝试次数
	Backoff     time.Duration // 重试退避
	ProxyAPIKey string        // 客户端鉴权 key，空 = 不鉴权（仅本地学习用）
}

// Load 读取并校验配置。
// 格式：
//   UPSTREAMS="主=https://api.openai.com/v1,备=http://localhost:4001/v1"
//     （别名=URL，逗号分隔；API key 从 UPSTREAM_KEYS 按位置对应，可省略）
//   PROXY_ADDR=:8080  ATTEMPTS=2  BACKOFF_MS=300  PROXY_API_KEY=
func Load() (*Config, error) {
	cfg := &Config{
		Addr:      getenv("PROXY_ADDR", ":8080"),
		Attempts:  getenvInt("ATTEMPTS", 2),
		Backoff:   time.Duration(getenvInt("BACKOFF_MS", 300)) * time.Millisecond,
		ProxyAPIKey: os.Getenv("PROXY_API_KEY"),
	}

	raw := getenv("UPSTREAMS", "")
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("UPSTREAMS 未设置，格式：别名=URL[,别名=URL...]")
	}
	keys := strings.Split(os.Getenv("UPSTREAM_KEYS"), ",")

	for i, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		name, url, found := strings.Cut(part, "=")
		if !found {
			return nil, fmt.Errorf("UPSTREAMS 第 %d 段缺少 '='：%q（应为 别名=URL）", i+1, part)
		}
		u := Upstream{Name: strings.TrimSpace(name), BaseURL: strings.TrimSpace(url)}
		if i < len(keys) && strings.TrimSpace(keys[i]) != "" {
			u.APIKey = strings.TrimSpace(keys[i])
		}
		cfg.Upstreams = append(cfg.Upstreams, u)
	}
	if len(cfg.Upstreams) == 0 {
		return nil, fmt.Errorf("UPSTREAMS 解析后为空")
	}
	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
