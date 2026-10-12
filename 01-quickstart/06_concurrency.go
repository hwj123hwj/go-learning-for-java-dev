package main

import (
	"fmt"
	"time"
)

func worker(id int, ch chan string) {
	time.Sleep(time.Second)
	ch <- fmt.Sprintf("Worker %d 完成!", id)
}

func main() {
	fmt.Println("=== Go并发: Goroutine + Channel ===")
	fmt.Println("")

	ch := make(chan string)

	// 从启动 goroutine 前开始计时，量到的是并行墙钟时间而非累加时间
	start := time.Now()

	go worker(1, ch)
	go worker(2, ch)
	go worker(3, ch)

	for i := 0; i < 3; i++ {
		msg := <-ch
		fmt.Println("收到:", msg)
	}

	// 3 个 worker 各睡 1 秒，并行执行总耗时约 1 秒，而不是串行的 3 秒
	fmt.Printf("3个worker总耗时: %v（若串行需3秒）\n", time.Since(start))

	fmt.Println("\n=== Buffered Channel ===")
	bufCh := make(chan int, 3)
	for i := 1; i <= 3; i++ {
		bufCh <- i
	}
	close(bufCh)
	fmt.Print("Buffered Channel值: ")
	for v := range bufCh {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	fmt.Println("\n=== Select (类似Java的NIO Selector) ===")
	intCh := make(chan int, 1)
	strCh := make(chan string, 1)

	intCh <- 42
	strCh <- "hello"

	select {
	case v := <-intCh:
		fmt.Printf("收到int: %d\n", v)
	case v := <-strCh:
		fmt.Printf("收到string: %s\n", v)
	default:
		fmt.Println("没有数据")
	}

	fmt.Println("\n=== 无缓冲 vs 有缓冲：发送方会不会阻塞 ===")

	// 无缓冲：发送和接收必须当场对上，主 goroutine 故意晚 2 秒才接收，
	// 子 goroutine 的"发送后"就会被压到 2 秒后才打印
	fmt.Println("--- 无缓冲 channel：发送方必须等到接收方 ---")
	t0 := time.Now()
	unbufCh := make(chan string)
	go func() {
		fmt.Printf("[%v] 子goroutine：发送前\n", time.Since(t0).Round(time.Millisecond))
		unbufCh <- "msg"
		fmt.Printf("[%v] 子goroutine：发送后（此刻才有人接收）\n", time.Since(t0).Round(time.Millisecond))
	}()
	time.Sleep(2 * time.Second)
	fmt.Printf("[%v] 主goroutine：开始接收 %q\n", time.Since(t0).Round(time.Millisecond), <-unbufCh)
	time.Sleep(100 * time.Millisecond) // 留出调度时间，确保子 goroutine 打印"发送后"

	// 有缓冲：容量 1 能暂存一个值，发送方扔进缓冲区立刻返回，
	// 即使主 goroutine 2 秒后才来接收，"发送后"也在 0 时刻就打印
	fmt.Println("--- 有缓冲 channel（容量1）：发送方扔进箱子就走 ---")
	t1 := time.Now()
	bufCh2 := make(chan string, 1)
	fmt.Printf("[%v] 发送前\n", time.Since(t1).Round(time.Millisecond))
	bufCh2 <- "msg"
	fmt.Printf("[%v] 发送后（无人接收也立刻返回）\n", time.Since(t1).Round(time.Millisecond))
	time.Sleep(2 * time.Second)
	fmt.Printf("[%v] 主goroutine：才来接收 %q\n", time.Since(t1).Round(time.Millisecond), <-bufCh2)
}
