package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("=== goroutine泄漏演示与防范 ===")
	fmt.Println()

	leakDemo()
	fmt.Println()

	fixWithCloseChannel()
	fmt.Println()

	fixWithContext()
	fmt.Println()

	pprofTip()
}

func leakDemo() {
	fmt.Println("--- 问题：goroutine泄漏 ---")

	ch := make(chan int)

	go func() {
		fmt.Println("goroutine等待接收...")
		<-ch
		fmt.Println("goroutine收到数据")
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("泄漏后goroutine数量: %d\n", runtime.NumGoroutine())
	fmt.Println("问题：main函数结束，但goroutine还在等待，造成泄漏")
	fmt.Println()
	fmt.Println("Java对比：")
	fmt.Println("  Java线程泄漏同样存在，但Go工具链（pprof）更易排查")
}

func fixWithCloseChannel() {
	fmt.Println("--- 解决方案1：close channel退出 ---")

	ch := make(chan int)

	go func() {
		defer fmt.Println("goroutine正常退出")
		for v := range ch {
			fmt.Println("收到:", v)
		}
	}()

	ch <- 1
	ch <- 2
	close(ch)

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("正常退出后goroutine数量: %d\n", runtime.NumGoroutine())

	fmt.Println()
	fmt.Println("Go惯用法：")
	fmt.Println("  - 使用close(channel)通知所有接收者退出")
	fmt.Println("  - for range channel会自动感知close")
}

func fixWithContext() {
	fmt.Println("--- 解决方案2：context取消 ---")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer fmt.Println("goroutine通过context退出")
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Println("工作中...")
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	time.Sleep(500 * time.Millisecond)
	cancel()

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("context取消后goroutine数量: %d\n", runtime.NumGoroutine())

	fmt.Println()
	fmt.Println("Go惯用法：")
	fmt.Println("  - context是控制goroutine生命周期的标准方式")
	fmt.Println("  - context.WithTimeout/WithCancel/WithDeadline")
	fmt.Println("  - context取消会自动传播到子goroutine")
}

func pprofTip() {
	fmt.Println("--- pprof排查goroutine泄漏 ---")
	fmt.Println()
	fmt.Println("排查方法:")
	fmt.Println("  1. 启动pprof: go tool pprof http://localhost:6060/debug/pprof/goroutine")
	fmt.Println("  2. 查看goroutine数量: runtime.NumGoroutine()")
	fmt.Println("  3. 导出goroutine栈: http://localhost:6060/debug/pprof/goroutine?debug=1")
	fmt.Println()
	fmt.Println("代码中启用pprof:")
	fmt.Println(`  import _ "net/http/pprof"
  go func() {
      http.ListenAndServe("localhost:6060", nil)
  }()`)
}
