package main

import (
	"fmt"
	"sync"
	"time"
)

type TicketSystemMutex struct {
	mu     sync.Mutex
	tickets int
}

func (t *TicketSystemMutex) Buy() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.tickets > 0 {
		t.tickets--
		return true
	}
	return false
}

func main() {
	fmt.Println("=== 抢票方案二：sync.Mutex ===")
	fmt.Println()

	system := &TicketSystemMutex{tickets: 100}
	var wg sync.WaitGroup
	var successCount int
	var countMu sync.Mutex

	start := time.Now()

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if system.Buy() {
				countMu.Lock()
				successCount++
				countMu.Unlock()
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
	fmt.Println("结果：正确！没有超卖")
	fmt.Println()
	fmt.Println("原理：")
	fmt.Println("  - Mutex保护临界区")
	fmt.Println("  - 同一时刻只有一个goroutine能操作tickets")
	fmt.Println()
	fmt.Println("Go惯用法：")
	fmt.Println("  - 生产中最常用的方案")
	fmt.Println("  - 代码简单易懂")
	fmt.Println("  - 适合复杂业务逻辑")
	fmt.Println()
	fmt.Println("Java对比：")
	fmt.Println("  synchronized (this) {")
	fmt.Println("      if (tickets > 0) {")
	fmt.Println("          tickets--;")
	fmt.Println("          return true;")
	fmt.Println("      }")
	fmt.Println("      return false;")
	fmt.Println("  }")
}
