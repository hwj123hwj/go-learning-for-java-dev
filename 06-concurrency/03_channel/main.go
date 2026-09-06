package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Channel基础 ===")
	fmt.Println()

	unbufferedChannel()
	fmt.Println()

	bufferedChannel()
	fmt.Println()

	compareWithJava()
}

func unbufferedChannel() {
	fmt.Println("--- 无缓冲channel ---")

	ch := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- "hello"
	}()

	msg := <-ch
	fmt.Println("收到:", msg)

	fmt.Println()
	fmt.Println("特点:")
	fmt.Println("  - 发送和接收必须同时准备好")
	fmt.Println("  - 同步通信，发送方阻塞直到接收方准备好")
	fmt.Println("  - 类似Java的SynchronousQueue")
}

func bufferedChannel() {
	fmt.Println("--- 有缓冲channel ---")

	ch := make(chan string, 2)

	ch <- "消息1"
	ch <- "消息2"

	fmt.Println("缓冲区已满，再发送会阻塞")

	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("接收:", <-ch)
	}()

	ch <- "消息3"

	fmt.Println("发送成功")

	fmt.Println()
	fmt.Println("特点:")
	fmt.Println("  - 缓冲区未满时，发送不阻塞")
	fmt.Println("  - 缓冲区满时，发送阻塞")
	fmt.Println("  - 类似Java的ArrayBlockingQueue")
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java BlockingQueue:")
	fmt.Println(`  BlockingQueue<String> queue = new ArrayBlockingQueue<>(10);
  queue.put("message");  // 阻塞发送
  String msg = queue.take();  // 阻塞接收`)
	fmt.Println()
	fmt.Println("Go channel:")
	fmt.Println(`  ch := make(chan string, 10)
  ch <- "message"    // 阻塞发送
  msg := <-ch        // 阻塞接收`)
	fmt.Println()
	fmt.Println("核心差异:")
	fmt.Println("  1. Go channel是语言内置，Java需要类库支持")
	fmt.Println("  2. Go channel可以close，Java BlockingQueue不能")
	fmt.Println("  3. Go channel支持select多路复用")
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - channel用于goroutine间通信")
	fmt.Println("  - 不要把channel当作队列滥用")
	fmt.Println("  - 无缓冲channel用于同步，有缓冲channel用于缓冲")
}
