package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	fmt.Println("=== sync.Mutex vs ReentrantLock ===")
	fmt.Println()
	basicMutex()
	fmt.Println()
	nonReentrantDemo()
	fmt.Println()
	deadlockDemo()
	fmt.Println()
	compareWithJava()
}

func basicMutex() {
	fmt.Println("--- 基础互斥锁 ---")
	counter := &Counter{}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Printf("最终计数: %d (期望1000)\n", counter.Value())
}

func nonReentrantDemo() {
	fmt.Println("--- Mutex不可重入演示 ---")
	fmt.Println()
	fmt.Println("Go的Mutex不可重入：同一goroutine重复Lock会永久阻塞")
	fmt.Println("fatal error 无法被 recover() 捕获，与Java异常机制不同")
	fmt.Println()
	fmt.Println("错误示例（勿执行）：")
	fmt.Println("  mu.Lock()")
	fmt.Println("  mu.Lock()  // 永久阻塞，recover()救不回来")
	fmt.Println()
	fmt.Println("Java ReentrantLock可重入，同一线程多次lock不会阻塞")
	fmt.Println()
	fmt.Println("Go正确做法：拆分内部函数，避免同一goroutine重复加锁：")
	fmt.Println("  func (c *Counter) lockedAdd() { c.value++ }")
	fmt.Println("  func (c *Counter) Increment() {")
	fmt.Println("      c.mu.Lock(); defer c.mu.Unlock(); c.lockedAdd()")
	fmt.Println("  }")
}

func deadlockDemo() {
	fmt.Println("--- 死锁场景演示 ---")
	fmt.Println()
	fmt.Println("1. 全局死锁（所有goroutine阻塞）:")
	fmt.Println("   Go运行时可检测: fatal error: all goroutines are asleep - deadlock!")
	fmt.Println()
	fmt.Println("2. 局部死锁（部分goroutine互相等待对方channel）:")
	fmt.Println("   Go运行时无法检测，程序不崩溃，goroutine永久阻塞")

	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	done := make(chan struct{})

	go func() {
		select {
		case <-ch1:
			ch2 <- struct{}{}
		case <-ctx.Done():
			fmt.Println("   goroutine A: 超时退出（模拟局部死锁的阻塞goroutine）")
		}
	}()

	go func() {
		select {
		case <-ch2:
			ch1 <- struct{}{}
		case <-ctx.Done():
			fmt.Println("   goroutine B: 超时退出（模拟局部死锁的阻塞goroutine）")
			close(done)
		}
	}()

	<-done
	fmt.Println()
	fmt.Println("   结论：局部死锁程序不崩溃，需用pprof排查泄漏的goroutine")
	fmt.Println()
	fmt.Println("预防措施:")
	fmt.Println("  - 用context + select设置超时，如上面演示")
	fmt.Println("  - 避免在锁内调用可能阻塞的外部函数")
	fmt.Println("  - 保持临界区小而简单")
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java ReentrantLock（可重入）：")
	fmt.Println("  lock.lock()")
	fmt.Println("  try { lock.lock(); lock.unlock(); } finally { lock.unlock(); }")
	fmt.Println()
	fmt.Println("Go sync.Mutex（不可重入）：")
	fmt.Println("  mu.Lock(); defer mu.Unlock()")
	fmt.Println("  // 同goroutine再次Lock永久阻塞，非panic")
	fmt.Println()
	fmt.Println("核心差异:")
	fmt.Println("  1. Go Mutex不可重入，Java ReentrantLock可重入")
	fmt.Println("  2. Go没有tryLock，Java有tryLock(timeout)")
	fmt.Println("  3. Go用defer unlock，Java用try-finally")
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - 临界区要小，避免在锁内调用外部函数")
	fmt.Println("  - 用defer unlock确保释放")
	fmt.Println("  - 需要重入时，重构代码拆分函数")
}
