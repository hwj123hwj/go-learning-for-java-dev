package main

import (
	"fmt"
	"sync"
)

type Buffer struct {
	Data []byte
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("创建新Buffer")
		return &Buffer{Data: make([]byte, 1024)}
	},
}

func main() {
	fmt.Println("=== sync.Pool ===")
	fmt.Println()

	basicPool()
	fmt.Println()

	compareWithJava()
}

func basicPool() {
	fmt.Println("--- 基础用法 ---")

	buf1 := bufferPool.Get().(*Buffer)
	fmt.Printf("获取Buffer: %p, len=%d\n", buf1, len(buf1.Data))
	buf1.Data = buf1.Data[:0]
	bufferPool.Put(buf1)

	buf2 := bufferPool.Get().(*Buffer)
	fmt.Printf("再次获取: %p (可能复用)\n", buf2)
	bufferPool.Put(buf2)

	fmt.Println()
	fmt.Println("注意：")
	fmt.Println("  - sync.Pool中的对象可能被GC回收")
	fmt.Println("  - 不适合做连接池（数据库连接、Redis连接等）")
	fmt.Println("  - 适合临时对象的复用，减少GC压力")
}

func compareWithJava() {
	fmt.Println("--- Java对比 ---")
	fmt.Println()
	fmt.Println("Java没有直接对应的标准库")
	fmt.Println("常见做法：")
	fmt.Println("  - Apache Commons Pool：通用对象池")
	fmt.Println("  - 数据库连接池：HikariCP、Druid")
	fmt.Println("  - 线程池：ExecutorService")
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - sync.Pool用于临时对象复用")
	fmt.Println("  - 连接池用专门的库（如sql.DB自带连接池）")
	fmt.Println("  - 只有GC压力大时才考虑Pool")
	fmt.Println()
	fmt.Println("⚠️ 重要警告:")
	fmt.Println("  sync.Pool不适合做连接池！")
	fmt.Println("  原因：Pool中的对象可能随时被GC回收")
	fmt.Println("  连接池推荐：database/sql、go-redis等库自带")
}
