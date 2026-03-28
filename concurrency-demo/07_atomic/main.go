package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Println("=== atomic原子操作 ===")
	fmt.Println()

	atomicCounter()
	fmt.Println()

	compareWithMutex()
	fmt.Println()

	compareWithJava()
}

func atomicCounter() {
	fmt.Println("--- atomic.Int64计数器 ---")

	var counter atomic.Int64
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("计数结果: %d (期望10000)\n", counter.Load())
	fmt.Printf("耗时: %v\n", elapsed)
}

func compareWithMutex() {
	fmt.Println("--- atomic vs Mutex性能对比 ---")

	var counter atomic.Int64
	var mu sync.Mutex
	var muCounter int
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}
	wg.Wait()
	atomicTime := time.Since(start)

	start = time.Now()
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			muCounter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	mutexTime := time.Since(start)

	fmt.Printf("atomic耗时: %v\n", atomicTime)
	fmt.Printf("Mutex耗时: %v\n", mutexTime)
	fmt.Printf("atomic快 %.1f%%\n", float64(mutexTime-atomicTime)/float64(mutexTime)*100)

	fmt.Println()
	fmt.Println("结论:")
	fmt.Println("  - 单一数值读写，atomic性能更好")
	fmt.Println("  - 复合操作（多个变量），必须用Mutex")
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java AtomicInteger:")
	fmt.Println(`  AtomicInteger counter = new AtomicInteger(0);
  counter.incrementAndGet();  // 原子递增
  counter.compareAndSet(expect, update);  // CAS操作`)
	fmt.Println()
	fmt.Println("Go atomic.Int64 (Go 1.19+):")
	fmt.Println(`  var counter atomic.Int64
  counter.Add(1)           // 原子递增
  counter.CompareAndSwap(old, new)  // CAS操作`)
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - Go 1.19+推荐使用atomic.Int64、atomic.Bool等类型安全封装")
	fmt.Println("  - 单一数值读写用atomic")
	fmt.Println("  - 复合操作用Mutex")
	fmt.Println("  - atomic适合计数器、标志位等简单场景")
}
