package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== sync.WaitGroup ===")
	fmt.Println()

	basicWaitGroup()
	fmt.Println()

	compareWithJava()
}

func basicWaitGroup() {
	fmt.Println("--- 基础用法 ---")

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(id*100) * time.Millisecond)
			fmt.Printf("任务 %d 完成\n", id)
		}(i)
	}

	fmt.Println("等待所有任务完成...")
	wg.Wait()
	fmt.Println("所有任务已完成")
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java CountDownLatch:")
	fmt.Println(`  CountDownLatch latch = new CountDownLatch(5);
  for (int i = 0; i < 5; i++) {
      new Thread(() -> {
          try {
              Thread.sleep(100);
              System.out.println("任务完成");
          } finally {
              latch.countDown();
          }
      }).start();
  }
  latch.await();
  System.out.println("所有任务完成");`)
	fmt.Println()
	fmt.Println("Go sync.WaitGroup:")
	fmt.Println(`  var wg sync.WaitGroup
  for i := 0; i < 5; i++ {
      wg.Add(1)
      go func() {
          defer wg.Done()
          time.Sleep(100 * time.Millisecond)
          fmt.Println("任务完成")
      }()
  }
  wg.Wait()
  fmt.Println("所有任务完成")`)
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - Add在goroutine外调用")
	fmt.Println("  - Done用defer确保调用")
	fmt.Println("  - Wait阻塞直到计数归零")
}
