package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	fmt.Println("=== 限流器（golang.org/x/time/rate）===")
	fmt.Println()

	demoAllow()
	fmt.Println()

	demoWait()
	fmt.Println()

	demoHTTPMiddleware()
	fmt.Println()

	compareWithJava()
}

func demoAllow() {
	fmt.Println("--- Allow()：非阻塞限流 ---")
	fmt.Println("每秒2个令牌，桶容量3，模拟10个请求")
	fmt.Println()

	// rate.Limit(2) = 每秒2个令牌；burst=3 表示最多积累3个令牌（允许突发）
	limiter := rate.NewLimiter(2, 3)

	pass, reject := 0, 0
	for i := 1; i <= 10; i++ {
		if limiter.Allow() {
			fmt.Printf("  请求 %2d [%s]: 通过\n", i, time.Now().Format("15:04:05.000"))
			pass++
		} else {
			fmt.Printf("  请求 %2d [%s]: 被限流\n", i, time.Now().Format("15:04:05.000"))
			reject++
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\n通过: %d  被限流: %d\n", pass, reject)
	fmt.Println("说明：桶初始有3个令牌（burst），之后每秒补充2个")
}

func demoWait() {
	fmt.Println("--- Wait()：阻塞等待令牌 ---")
	fmt.Println("每秒5个令牌，5个并发请求，Wait会阻塞直到拿到令牌")
	fmt.Println()

	limiter := rate.NewLimiter(5, 5)
	var wg sync.WaitGroup
	start := time.Now()

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			if err := limiter.Wait(ctx); err != nil {
				fmt.Printf("  请求 %2d: 超时或取消 - %v\n", id, err)
				return
			}
			fmt.Printf("  请求 %2d: 获得令牌，耗时 %v\n", id, time.Since(start).Round(time.Millisecond))
		}(i)
	}
	wg.Wait()

	fmt.Println()
	fmt.Println("说明：Wait()不会丢弃请求，而是排队等待，配合context可设超时")
}

// 模拟HTTP处理器
func handleRequest(id int) string {
	time.Sleep(5 * time.Millisecond)
	return fmt.Sprintf("response-%d", id)
}

func demoHTTPMiddleware() {
	fmt.Println("--- HTTP中间件限流模拟 ---")
	fmt.Println("每秒3个令牌，模拟HTTP请求经过限流中间件")
	fmt.Println()

	limiter := rate.NewLimiter(3, 3)

	var wg sync.WaitGroup
	pass, reject := 0, 0
	var mu sync.Mutex

	for i := 1; i <= 9; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// 中间件逻辑：非阻塞检查
			if !limiter.Allow() {
				fmt.Printf("  请求 %2d: 429 Too Many Requests\n", id)
				mu.Lock()
				reject++
				mu.Unlock()
				return
			}
			resp := handleRequest(id)
			fmt.Printf("  请求 %2d: 200 OK -> %s\n", id, resp)
			mu.Lock()
			pass++
			mu.Unlock()
		}(i)
		time.Sleep(80 * time.Millisecond)
	}
	wg.Wait()

	fmt.Printf("\n200 OK: %d  429限流: %d\n", pass, reject)
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java Guava RateLimiter:")
	fmt.Println("  RateLimiter limiter = RateLimiter.create(10.0);")
	fmt.Println("  limiter.acquire();        // 阻塞，无法设置超时")
	fmt.Println("  limiter.tryAcquire();     // 非阻塞")
	fmt.Println("  limiter.tryAcquire(100, TimeUnit.MILLISECONDS); // 带超时")
	fmt.Println()
	fmt.Println("Go golang.org/x/time/rate:")
	fmt.Println("  limiter := rate.NewLimiter(10, 100)")
	fmt.Println("  limiter.Allow()           // 非阻塞，返回bool")
	fmt.Println("  limiter.Wait(ctx)         // 阻塞，ctx可设超时/取消")
	fmt.Println("  limiter.Reserve()         // 预留令牌，返回等待时间")
	fmt.Println()
	fmt.Println("核心差异:")
	fmt.Println("  - Go的Wait()天然支持context，超时/取消更灵活")
	fmt.Println("  - Guava支持SmoothWarmingUp（平滑预热），x/time/rate暂不支持")
	fmt.Println("  - Go的rate.Inf表示不限速，rate.Every()从时间间隔算速率")
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - HTTP中间件用Allow()快速拒绝")
	fmt.Println("  - 消息队列消费者用Wait(ctx)平滑消费")
	fmt.Println("  - Burst设为速率的1~2倍，允许短暂突发")
}
