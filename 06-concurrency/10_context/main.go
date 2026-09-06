package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Context超时控制 ===")
	fmt.Println()

	withTimeout()
	fmt.Println()

	withCancel()
	fmt.Println()

	chainPropagation()
	fmt.Println()

	compareWithJava()
}

func withTimeout() {
	fmt.Println("--- context.WithTimeout ---")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	select {
	case <-time.After(200 * time.Millisecond):
		fmt.Println("操作完成")
	case <-ctx.Done():
		fmt.Printf("超时取消！耗时: %v\n", time.Since(start))
		fmt.Printf("错误: %v\n", ctx.Err())
	}
}

func withCancel() {
	fmt.Println("--- context.WithCancel ---")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("主动取消")
		cancel()
	}()

	start := time.Now()
	select {
	case <-time.After(200 * time.Millisecond):
		fmt.Println("操作完成")
	case <-ctx.Done():
		fmt.Printf("收到取消信号！耗时: %v\n", time.Since(start))
		fmt.Printf("错误: %v\n", ctx.Err())
	}
}

func chainPropagation() {
	fmt.Println("--- 链式传播 ---")

	parentCtx, parentCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer parentCancel()

	childCtx, childCancel := context.WithCancel(parentCtx)
	defer childCancel()

	go func() {
		select {
		case <-childCtx.Done():
			fmt.Printf("子goroutine收到取消: %v\n", childCtx.Err())
		}
	}()

	time.Sleep(150 * time.Millisecond)
	fmt.Printf("父context状态: %v\n", parentCtx.Err())
	fmt.Printf("子context状态: %v\n", childCtx.Err())
	fmt.Println()
	fmt.Println("结论：父cancel自动传播到子")
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java超时控制:")
	fmt.Println(`  ExecutorService executor = Executors.newCachedThreadPool();
  Future<?> future = executor.submit(() -> {
      // 长时间操作
  });
  try {
      future.get(100, TimeUnit.MILLISECONDS);
  } catch (TimeoutException e) {
      future.cancel(true);  // 中断线程
  }`)
	fmt.Println()
	fmt.Println("Go Context:")
	fmt.Println(`  ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
  defer cancel()
  
  select {
  case <-doWork(ctx):
      fmt.Println("完成")
  case <-ctx.Done():
      fmt.Println("超时:", ctx.Err())
  }`)
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - context作为函数第一个参数传递")
	fmt.Println("  - 不要在struct中存储context")
	fmt.Println("  - 总是用defer cancel()确保释放资源")
	fmt.Println("  - context取消会自动传播到所有子goroutine")
}
