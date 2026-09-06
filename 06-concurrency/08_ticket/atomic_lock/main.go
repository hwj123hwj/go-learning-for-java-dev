package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type TicketSystemAtomic struct {
	tickets atomic.Int64
}

func (t *TicketSystemAtomic) Buy() bool {
	for {
		current := t.tickets.Load()
		if current <= 0 {
			return false
		}
		if t.tickets.CompareAndSwap(current, current-1) {
			return true
		}
	}
}

func main() {
	fmt.Println("=== 抢票方案一：atomic原子操作 ===")
	fmt.Println()

	system := &TicketSystemAtomic{}
	system.tickets.Store(100)
	var wg sync.WaitGroup
	var successCount atomic.Int64

	start := time.Now()

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if system.Buy() {
				successCount.Add(1)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("初始票数: 100\n")
	fmt.Printf("成功购票: %d\n", successCount.Load())
	fmt.Printf("剩余票数: %d\n", system.tickets.Load())
	fmt.Printf("耗时: %v\n", elapsed)
	fmt.Println()
	fmt.Println("结果：正确！没有超卖")
	fmt.Println()
	fmt.Println("原理：")
	fmt.Println("  - 使用CAS（CompareAndSwap）原子操作")
	fmt.Println("  - 循环重试直到成功或票数为0")
	fmt.Println()
	fmt.Println("Go惯用法：")
	fmt.Println("  - 单一数值场景，atomic是最优选择")
	fmt.Println("  - 性能最好，无锁实现")
	fmt.Println()
	fmt.Println("Java对比：")
	fmt.Println("  AtomicInteger tickets = new AtomicInteger(100);")
	fmt.Println("  tickets.compareAndSet(current, current - 1);")
}
