package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== goroutine vs Java Thread 对比 ===")
	fmt.Println()

	basicGoroutine()
	fmt.Println()

	lightweightDemo()
	fmt.Println()

	compareWithJava()
}

func basicGoroutine() {
	fmt.Println("--- 基础用法 ---")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("goroutine 1 执行")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("goroutine 2 执行")
	}()

	wg.Wait()
	fmt.Println("所有goroutine完成")
}

func lightweightDemo() {
	fmt.Println("--- 轻量级演示：创建100万个goroutine ---")

	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(1000000)

	for i := 0; i < 1000000; i++ {
		go func() {
			defer wg.Done()
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("创建100万个goroutine耗时: %v\n", elapsed)
	fmt.Printf("当前goroutine数量: %d\n", runtime.NumGoroutine())
	fmt.Printf("堆内存使用: %.2f MB\n", float64(m.HeapAlloc)/1024/1024)

	fmt.Println()
	fmt.Println("对比Java：")
	fmt.Println("  Java Thread: 每个线程约1MB栈空间，100万线程需要约1TB内存")
	fmt.Println("  Go goroutine: 初始栈仅2KB，可动态增长，100万goroutine仅需约2GB")
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java Thread:")
	fmt.Println("  new Thread(() -> { ... }).start();")
	fmt.Println("  // 或使用线程池")
	fmt.Println("  executor.submit(() -> { ... });")
	fmt.Println()
	fmt.Println("Go goroutine:")
	fmt.Println("  go func() { ... }()")
	fmt.Println("  // 或")
	fmt.Println("  go someFunction()")
	fmt.Println()
	fmt.Println("核心差异:")
	fmt.Println("  1. 内存: goroutine初始2KB vs Thread约1MB")
	fmt.Println("  2. 调度: Go运行时调度 vs OS线程调度")
	fmt.Println("  3. 创建成本: goroutine极低 vs Thread较高")
	fmt.Println("  4. 切换成本: goroutine用户态 vs Thread内核态")
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - goroutine非常轻量，可大量创建")
	fmt.Println("  - 但必须注意goroutine泄漏问题")
	fmt.Println("  - 使用sync.WaitGroup等待goroutine完成")
	fmt.Println("  - 使用context或done channel控制退出")
}
