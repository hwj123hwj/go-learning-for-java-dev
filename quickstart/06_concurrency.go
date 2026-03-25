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

	go worker(1, ch)
	go worker(2, ch)
	go worker(3, ch)

	for i := 0; i < 3; i++ {
		msg := <-ch
		fmt.Println("收到:", msg)
	}

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
}
