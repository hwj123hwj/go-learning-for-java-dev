package main

import (
	"fmt"
	"sync"
	"time"
)

type TicketSystem struct {
	tickets int
}

func (t *TicketSystem) Buy() bool {
	if t.tickets > 0 {
		time.Sleep(1 * time.Microsecond)
		t.tickets--
		return true
	}
	return false
}

func main() {
	fmt.Println("=== 抢票问题演示（无锁）===")
	fmt.Println()

	system := &TicketSystem{tickets: 100}
	var wg sync.WaitGroup
	var successCount int
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if system.Buy() {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("初始票数: 100\n")
	fmt.Printf("成功购票: %d\n", successCount)
	fmt.Printf("剩余票数: %d\n", system.tickets)
	fmt.Printf("耗时: %v\n", elapsed)
	fmt.Println()
	fmt.Println("问题：超卖！成功购票数 > 初始票数")
	fmt.Println("原因：多个goroutine同时读取tickets > 0，都认为有票")
	fmt.Println()
	fmt.Println("Java对比：")
	fmt.Println("  Java多线程同样存在这个问题")
	fmt.Println("  解决方案：synchronized 或 ReentrantLock 或 AtomicInteger")
}
