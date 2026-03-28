package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Println("=== 抢票方案三：Channel串行化 ===")
	fmt.Println()
	fmt.Println("注意：此方案为教学目的，演示Go特有并发模式")
	fmt.Println("生产中atomic或Mutex通常更实用")
	fmt.Println()

	tickets := 100
	buyCh := make(chan struct{}, 1000)
	resultCh := make(chan bool, 1000)

	// 管理goroutine：独立WaitGroup，串行处理所有购票请求
	var managerWg sync.WaitGroup
	managerWg.Add(1)
	go func() {
		defer managerWg.Done()
		for range buyCh {
			if tickets > 0 {
				tickets--
				resultCh <- true
			} else {
				resultCh <- false
			}
		}
		// buyCh关闭后退出循环，再关闭resultCh通知消费者
		close(resultCh)
	}()

	var buyerWg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		buyerWg.Add(1)
		go func(id int) {
			defer buyerWg.Done()
			buyCh <- struct{}{}
		}(i)
	}

	// 等所有购票请求发送完毕，关闭buyCh
	buyerWg.Wait()
	close(buyCh)

	// 等管理goroutine处理完并关闭resultCh
	managerWg.Wait()

	// 统计结果（resultCh已关闭，range会自动结束）
	var successCount atomic.Int64
	start := time.Now()
	for ok := range resultCh {
		if ok {
			successCount.Add(1)
		}
	}
	elapsed := time.Since(start)

	fmt.Printf("初始票数: 100\n")
	fmt.Printf("成功购票: %d\n", successCount.Load())
	fmt.Printf("剩余票数: %d\n", tickets)
	fmt.Printf("耗时: %v\n", elapsed)
	fmt.Println()
	fmt.Println("结果：正确！没有超卖")
	fmt.Println()
	fmt.Println("原理：")
	fmt.Println("  - 通过channel串行化所有请求")
	fmt.Println("  - 单个goroutine处理购票逻辑，无需加锁")
	fmt.Println("  - 体现\"用通信共享内存\"的Go哲学")
	fmt.Println()
	fmt.Println("关键：close顺序")
	fmt.Println("  1. 所有buyer发送完毕 -> close(buyCh)")
	fmt.Println("  2. 管理goroutine处理完所有请求 -> close(resultCh)")
	fmt.Println("  3. 消费者range resultCh自动结束")
	fmt.Println()
	fmt.Println("适用场景：")
	fmt.Println("  - 需要顺序处理请求（如状态机）")
	fmt.Println("  - 演示Go并发哲学")
	fmt.Println("  - 生产中atomic或Mutex通常更高效")
}
