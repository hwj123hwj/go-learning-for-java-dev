package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redsync/redsync/v4"
	goredisv9 "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("=== 分布式锁（redsync + go-redis/v9）===")
	fmt.Println()
	fmt.Println("前提：需要本地 Redis，默认连接 localhost:6379")
	fmt.Println()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// 检查连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis连接失败: %v\n请先启动Redis: brew services start redis", err)
	}
	fmt.Println("Redis连接成功")
	fmt.Println()

	pool := goredisv9.NewPool(client)
	rs := redsync.New(pool)

	basicLock(rs)
	fmt.Println()

	concurrentLock(rs)
	fmt.Println()

	compareWithJava()
}

func basicLock(rs *redsync.Redsync) {
	fmt.Println("--- 基础加锁/解锁 ---")

	mutex := rs.NewMutex("demo:basic-lock",
		redsync.WithExpiry(10*time.Second), // 锁过期时间
		redsync.WithTries(3),               // 最多重试3次
		redsync.WithRetryDelay(100*time.Millisecond),
	)

	// 加锁
	if err := mutex.Lock(); err != nil {
		log.Fatalf("加锁失败: %v", err)
	}
	fmt.Println("加锁成功")

	// 模拟业务操作
	time.Sleep(100 * time.Millisecond)
	fmt.Println("执行业务逻辑...")

	// 释放锁
	if ok, err := mutex.Unlock(); !ok || err != nil {
		log.Fatalf("解锁失败: %v", err)
	}
	fmt.Println("解锁成功")
}

func concurrentLock(rs *redsync.Redsync) {
	fmt.Println("--- 并发竞争演示 ---")
	fmt.Println("5个goroutine竞争同一把锁，演示互斥效果")
	fmt.Println()

	results := make(chan string, 5)

	for i := 1; i <= 5; i++ {
		go func(id int) {
			mutex := rs.NewMutex("demo:concurrent-lock",
				redsync.WithExpiry(5*time.Second),
				redsync.WithTries(10),
				redsync.WithRetryDelay(50*time.Millisecond),
			)

			if err := mutex.Lock(); err != nil {
				results <- fmt.Sprintf("goroutine %d: 获取锁失败 - %v", id, err)
				return
			}

			msg := fmt.Sprintf("goroutine %d: 持有锁，执行业务", id)
			time.Sleep(80 * time.Millisecond)

			mutex.Unlock()
			results <- msg + " -> 释放锁"
		}(i)
	}

	for i := 0; i < 5; i++ {
		fmt.Println(" ", <-results)
	}

	fmt.Println()
	fmt.Println("效果：每次只有一个goroutine持有锁，其余等待重试")
	fmt.Println()
	fmt.Println("redsync关键特性:")
	fmt.Println("  - WithExpiry:     锁自动过期，防止持有者崩溃后死锁")
	fmt.Println("  - WithTries:      获取失败时自动重试")
	fmt.Println("  - WithRetryDelay: 重试间隔")
	fmt.Println("  - Unlock() 返回 bool：验证是否是自己持有的锁再释放（防误删）")
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java Redisson:")
	fmt.Println("  RLock lock = redisson.getLock(\"my-lock\");")
	fmt.Println("  lock.lock(10, TimeUnit.SECONDS);")
	fmt.Println("  try { /* 业务 */ } finally { lock.unlock(); }")
	fmt.Println()
	fmt.Println("Go redsync:")
	fmt.Println("  mutex := rs.NewMutex(\"my-lock\", redsync.WithExpiry(10*time.Second))")
	fmt.Println("  mutex.Lock()")
	fmt.Println("  defer mutex.Unlock()")
	fmt.Println()
	fmt.Println("核心差异:")
	fmt.Println("  - Redisson 功能更丰富（自动watchdog续期、可重入、红锁）")
	fmt.Println("  - redsync 专注分布式互斥锁，更轻量")
	fmt.Println("  - redsync 的 Unlock() 会验证锁归属，防止误删他人的锁")
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - 不手写分布式锁，使用 redsync")
	fmt.Println("  - WithExpiry 必须设置，防止持有者宕机导致死锁")
	fmt.Println("  - 用 defer mutex.Unlock() 确保释放")
}
