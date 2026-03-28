package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== sync.Map ===")
	fmt.Println()

	basicSyncMap()
	fmt.Println()

	comparePerformance()
	fmt.Println()

	compareWithJava()
}

func basicSyncMap() {
	fmt.Println("--- 基础用法 ---")

	var m sync.Map

	m.Store("name", "Alice")
	m.Store("age", 25)

	if val, ok := m.Load("name"); ok {
		fmt.Printf("name: %v\n", val)
	}

	m.Range(func(key, value interface{}) bool {
		fmt.Printf("%v: %v\n", key, value)
		return true
	})

	m.Delete("age")
}

func comparePerformance() {
	fmt.Println("--- 性能对比 ---")

	n := 100000

	var syncMap sync.Map
	start := time.Now()
	for i := 0; i < n; i++ {
		syncMap.Store(i, i)
	}
	for i := 0; i < n; i++ {
		syncMap.Load(i)
	}
	syncMapTime := time.Since(start)

	normalMap := make(map[int]int)
	var mu sync.RWMutex
	start = time.Now()
	for i := 0; i < n; i++ {
		mu.Lock()
		normalMap[i] = i
		mu.Unlock()
	}
	for i := 0; i < n; i++ {
		mu.RLock()
		_ = normalMap[i]
		mu.RUnlock()
	}
	mapTime := time.Since(start)

	fmt.Printf("sync.Map: %v\n", syncMapTime)
	fmt.Printf("map+RWMutex: %v\n", mapTime)

	fmt.Println()
	fmt.Println("结论:")
	fmt.Println("  - 写多场景：map+RWMutex更快")
	fmt.Println("  - 读多写少：sync.Map可能更快")
	fmt.Println("  - key稳定：sync.Map更适合")
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java ConcurrentHashMap:")
	fmt.Println("  ConcurrentHashMap<String, Integer> map = new ConcurrentHashMap<>();")
	fmt.Println("  map.put(\"key\", 1);")
	fmt.Println("  Integer val = map.get(\"key\");")
	fmt.Println()
	fmt.Println("Go sync.Map:")
	fmt.Println("  var m sync.Map")
	fmt.Println("  m.Store(\"key\", 1)")
	fmt.Println("  val, ok := m.Load(\"key\")")
	fmt.Println()
	fmt.Println("⚠️ 重要区别:")
	fmt.Println("  sync.Map不是ConcurrentHashMap的直接替代！")
	fmt.Println("  sync.Map适用场景有限：读多写少、key稳定")
	fmt.Println("  写多场景性能反而差于map+RWMutex")
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - 大多数场景用map+Mutex更合适")
	fmt.Println("  - 只有确认sync.Map更适合时才使用")
	fmt.Println("  - 不要无脑使用sync.Map")
}
