package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== select多路复用 ===")
	fmt.Println()

	basicSelect()
	fmt.Println()

	nonBlockingSelect()
	fmt.Println()

	timeoutSelect()
	fmt.Println()

	compareWithJava()
}

func basicSelect() {
	fmt.Println("--- 基础select ---")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "来自ch1"
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- "来自ch2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("ch1收到:", msg1)
		case msg2 := <-ch2:
			fmt.Println("ch2收到:", msg2)
		}
	}

	fmt.Println()
	fmt.Println("特点:")
	fmt.Println("  - 同时监听多个channel")
	fmt.Println("  - 随机选择一个就绪的case执行")
}

func nonBlockingSelect() {
	fmt.Println("--- 非阻塞select (default分支) ---")

	ch := make(chan string, 1)

	select {
	case msg := <-ch:
		fmt.Println("收到:", msg)
	default:
		fmt.Println("没有数据，执行default分支")
	}

	ch <- "hello"

	select {
	case msg := <-ch:
		fmt.Println("收到:", msg)
	default:
		fmt.Println("没有数据，执行default分支")
	}

	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - default分支实现非阻塞操作")
	fmt.Println("  - 常用于轮询或心跳检测")
}

func timeoutSelect() {
	fmt.Println("--- 超时select ---")

	ch := make(chan string)

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch <- "慢响应"
	}()

	select {
	case msg := <-ch:
		fmt.Println("收到:", msg)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("超时！100ms内未收到数据")
	}

	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - time.After()返回一个定时channel")
	fmt.Println("  - 配合select实现超时控制")
	fmt.Println("  - 比Java的Future.get(timeout)更灵活")
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java wait/notify多路等待:")
	fmt.Println(`  synchronized(obj1) {
      obj1.wait();  // 只能等待一个
  }
  // 多条件需要复杂的条件判断`)
	fmt.Println()
	fmt.Println("Go select:")
	fmt.Println(`  select {
  case msg1 := <-ch1:
      // 处理ch1
  case msg2 := <-ch2:
      // 处理ch2
  case <-time.After(100 * time.Millisecond):
      // 超时处理
  }`)
	fmt.Println()
	fmt.Println("核心差异:")
	fmt.Println("  1. select是Go语言内置特性")
	fmt.Println("  2. select天然支持多路复用")
	fmt.Println("  3. select + channel = Go并发核心模式")
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - select是Go并发的核心")
	fmt.Println("  - 相当于Java中wait/notify + 条件判断的组合")
	fmt.Println("  - 配合context实现优雅退出")
}
