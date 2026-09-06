# 后端常见业务问题解决方案

> Go语言实现后端常见业务问题的最佳实践，对比Java解决方案。

## 目录结构

```
backend-patterns/
├── distributed_lock/    # 分布式锁
├── rate_limiter/        # 限流器
└── cache_problems/      # 缓存问题（穿透/击穿/雪崩）
```

## 快速开始

```bash
cd backend-patterns
go mod tidy
go run ./distributed_lock/main.go
```

## 核心对比

| 问题 | Java方案 | Go方案 |
|------|---------|--------|
| 分布式锁 | Redisson | go-redis + redsync |
| 限流 | Guava RateLimiter | golang.org/x/time/rate |
| 缓存击穿 | synchronized | singleflight |
| 缓存穿透 | 布隆过滤器 | bloomfilter |
| 缓存雪崩 | 随机TTL | 随机TTL |

## Go惯用法原则

1. **不要手写分布式锁** - 使用成熟库（redsync）
2. **限流用标准库** - golang.org/x/time/rate
3. **singleflight防击穿** - golang.org/x/sync/singleflight
4. **布隆过滤器防穿透** - 适合海量数据判断
