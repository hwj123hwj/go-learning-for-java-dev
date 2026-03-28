package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]string)}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func main() {
	fmt.Println("=== sync.RWMutex vs ReentrantReadWriteLock ===")
	fmt.Println()
	rwMutexDemo()
	fmt.Println()
	performanceCompare()
	fmt.Println()
	compareWithJava()
}

func rwMutexDemo() {
	fmt.Println("--- 读写锁演示 ---")
	cache := NewCache()
	cache.Set("name", "Alice")
	cache.Set("city", "Beijing")

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if val, ok := cache.Get("name"); ok {
				fmt.Printf("Reader %d: name=%s\n", id, val)
			}
		}(i)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		cache.Set("name", "Bob")
		fmt.Println("Writer: updated name")
	}()
	wg.Wait()
}

func performanceCompare() {
	fmt.Println("--- 并发读场景性能对比（读多写少）---")
	data := make(map[int]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	var mu sync.Mutex
	var rwmu sync.RWMutex
	const readers = 20
	const ops = 5000
	var wg sync.WaitGroup

	// Mutex并发读
	start := time.Now()
	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				mu.Lock()
				_ = data[i%1000]
				mu.Unlock()
			}
		}(r)
	}
	wg.Wait()
	mutexTime := time.Since(start)

	// RWMutex并发读
	start = time.Now()
	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				rwmu.RLock()
				_ = data[i%1000]
				rwmu.RUnlock()
			}
		}(r)
	}
	wg.Wait()
	rwmuTime := time.Since(start)

	fmt.Printf("Mutex   并发读(%d goroutine x %d ops): %v\n", readers, ops, mutexTime)
	fmt.Printf("RWMutex 并发读(%d goroutine x %d ops): %v\n", readers, ops, rwmuTime)
	if rwmuTime < mutexTime {
		fmt.Printf("RWMutex 快 %.1f%%\n", float64(mutexTime-rwmuTime)/float64(mutexTime)*100)
	} else {
		fmt.Println("本次Mutex更快（读操作极短时goroutine调度开销影响更大，属正常）")
	}
	fmt.Println()
	fmt.Println("结论:")
	fmt.Println("  - 读多写少且读操作耗时较长时，RWMutex性能更好")
	fmt.Println("  - 写多场景RWMutex反而更慢，用Mutex即可")
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java ReentrantReadWriteLock:")
	fmt.Println("  rwlock.readLock().lock();  try { ... } finally { rwlock.readLock().unlock(); }")
	fmt.Println("  rwlock.writeLock().lock(); try { ... } finally { rwlock.writeLock().unlock(); }")
	fmt.Println()
	fmt.Println("Go sync.RWMutex:")
	fmt.Println("  rwmu.RLock(); defer rwmu.RUnlock()  // 读")
	fmt.Println("  rwmu.Lock();  defer rwmu.Unlock()   // 写")
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - 只有真正读多写少时才用RWMutex")
	fmt.Println("  - 读操作用RLock/RUnlock，写操作用Lock/Unlock")
	fmt.Println("  - 用defer确保释放")
}
