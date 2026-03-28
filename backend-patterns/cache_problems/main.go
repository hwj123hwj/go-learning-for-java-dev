package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"golang.org/x/sync/singleflight"
)

func main() {
	fmt.Println("=== 缓存问题解决方案 ===")
	fmt.Println()

	demoCachePenetration()
	fmt.Println()

	demoCacheBreakdown()
	fmt.Println()

	demoCacheAvalanche()
	fmt.Println()

	compareWithJava()
}

// ─── 缓存穿透：布隆过滤器 ────────────────────────────────────────────────────

// 模拟数据库：只有 user:1 ~ user:100 存在
var db = func() map[string]string {
	m := make(map[string]string)
	for i := 1; i <= 100; i++ {
		key := fmt.Sprintf("user:%d", i)
		m[key] = fmt.Sprintf("Alice-%d", i)
	}
	return m
}()

var dbQueryCount int
var dbMu sync.Mutex

func queryDB(key string) (string, bool) {
	dbMu.Lock()
	dbQueryCount++
	dbMu.Unlock()
	time.Sleep(1 * time.Millisecond) // 模拟DB延迟
	val, ok := db[key]
	return val, ok
}

func demoCachePenetration() {
	fmt.Println("--- 缓存穿透：布隆过滤器 ---")
	fmt.Println()
	fmt.Println("问题：查询不存在的key，每次都打到数据库")
	fmt.Println()

	// 初始化布隆过滤器，预估100万元素，误判率1%
	filter := bloom.NewWithEstimates(1000000, 0.01)

	// 预热：将DB中所有存在的key加入过滤器
	for key := range db {
		filter.AddString(key)
	}
	fmt.Printf("布隆过滤器预热完成，加载 %d 个key\n", len(db))
	fmt.Println()

	// 模拟请求：50个存在的key + 50个不存在的key
	dbQueryCount = 0
	hits, blocked := 0, 0

	for i := 1; i <= 50; i++ {
		key := fmt.Sprintf("user:%d", i) // 存在
		if !filter.TestString(key) {
			blocked++
		} else {
			queryDB(key)
			hits++
		}
	}
	for i := 101; i <= 150; i++ {
		key := fmt.Sprintf("user:%d", i) // 不存在
		if !filter.TestString(key) {
			blocked++ // 布隆过滤器拦截，不查DB
		} else {
			// 假阳性：过滤器误判为存在，仍会查DB
			queryDB(key)
			hits++
		}
	}

	fmt.Printf("请求总数:       100\n")
	fmt.Printf("布隆过滤器拦截: %d（不存在的key，无需查DB）\n", blocked)
	fmt.Printf("实际查询DB:     %d次\n", dbQueryCount)
	fmt.Printf("假阳性（误判）: %d次\n", hits-50)
	fmt.Println()
	fmt.Println("特点：无假阴性（存在的key一定通过），有少量假阳性（可接受）")
	fmt.Println("Go惯用法：启动时预热过滤器，写入时同步更新过滤器")
}

// ─── 缓存击穿：singleflight ──────────────────────────────────────────────────

// 模拟内存缓存
type MemCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func (c *MemCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *MemCache) Set(key, val string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = val
}

var memCache = &MemCache{data: make(map[string]string)}
var sfGroup singleflight.Group
var realDBCallCount int
var realDBMu sync.Mutex

func getUserWithSingleflight(id string) (string, error) {
	// 先查缓存
	if val, ok := memCache.Get(id); ok {
		return val, nil
	}

	// 缓存未命中，用singleflight合并对同一key的并发请求
	val, err, shared := sfGroup.Do(id, func() (interface{}, error) {
		realDBMu.Lock()
		realDBCallCount++
		realDBMu.Unlock()

		time.Sleep(50 * time.Millisecond) // 模拟DB查询
		result, ok := queryDB(id)
		if !ok {
			return "", fmt.Errorf("not found")
		}
		// 写入缓存
		memCache.Set(id, result)
		return result, nil
	})

	_ = shared // shared=true 表示此次结果被多个调用者共享
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

func demoCacheBreakdown() {
	fmt.Println("--- 缓存击穿：singleflight ---")
	fmt.Println()
	fmt.Println("问题：热点key过期，大量并发请求同时打到数据库")
	fmt.Println()

	realDBCallCount = 0
	var wg sync.WaitGroup

	// 模拟50个goroutine同时请求同一个热点key
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getUserWithSingleflight("user:1")
		}()
	}
	wg.Wait()

	fmt.Printf("并发请求数:   50\n")
	fmt.Printf("实际查询DB:   %d次（singleflight合并了其余请求）\n", realDBCallCount)
	fmt.Println()
	fmt.Println("效果：50个并发请求只触发1次DB查询，其余49个等待结果共享")
	fmt.Println("Go惯用法：singleflight.Group 是处理缓存击穿的标准方案")
	fmt.Println("对比Java：Java用 synchronized/双检锁，Go用 singleflight 更简洁")
}

// ─── 缓存雪崩：随机TTL ───────────────────────────────────────────────────────

type CacheItem struct {
	value     string
	expireAt  time.Time
}

type TTLCache struct {
	mu   sync.RWMutex
	data map[string]CacheItem
}

func (c *TTLCache) Set(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheItem{value: value, expireAt: time.Now().Add(ttl)}
}

func (c *TTLCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.data[key]
	if !ok || time.Now().After(item.expireAt) {
		return "", false
	}
	return item.value, true
}

// randomTTL 在基础TTL上加随机抖动，避免大量key同时过期
func randomTTL(base time.Duration, jitter time.Duration) time.Duration {
	return base + time.Duration(rand.Int63n(int64(jitter)))
}

func demoCacheAvalanche() {
	fmt.Println("--- 缓存雪崩：随机TTL ---")
	fmt.Println()
	fmt.Println("问题：大量key设置相同TTL，同时过期导致DB被压垮")
	fmt.Println()

	cache := &TTLCache{data: make(map[string]CacheItem)}
	baseTTL := 1 * time.Second
	jitter := 500 * time.Millisecond

	// 写入100个key，每个TTL加随机抖动
	for i := 1; i <= 100; i++ {
		key := fmt.Sprintf("product:%d", i)
		ttl := randomTTL(baseTTL, jitter)
		cache.Set(key, fmt.Sprintf("data-%d", i), ttl)
	}

	// 统计1.0s、1.2s、1.5s时的过期数量
	checkpoints := []time.Duration{1000, 1200, 1500}
	start := time.Now()

	for _, ms := range checkpoints {
		target := start.Add(time.Duration(ms) * time.Millisecond)
		time.Sleep(time.Until(target))

		expired := 0
		for i := 1; i <= 100; i++ {
			key := fmt.Sprintf("product:%d", i)
			if _, ok := cache.Get(key); !ok {
				expired++
			}
		}
		fmt.Printf("经过 %dms：已过期 %d/100 个key\n", ms, expired)
	}

	fmt.Println()
	fmt.Println("效果：key分散在不同时间过期，DB压力平滑分布")
	fmt.Println("Go惯用法：base TTL + rand.Int63n(jitter) 是最简单有效的方案")
	fmt.Println("对比Java：逻辑相同，Go的time.Duration运算更直观")
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("问题       Java方案                    Go方案")
	fmt.Println("穿透       Guava BloomFilter            bits-and-blooms/bloom")
	fmt.Println("击穿       synchronized / 双检锁        singleflight.Group")
	fmt.Println("雪崩       随机TTL（逻辑相同）          随机TTL（time.Duration）")
	fmt.Println()
	fmt.Println("Go核心优势：")
	fmt.Println("  - singleflight 内置于标准扩展库，无需额外框架")
	fmt.Println("  - time.Duration 类型安全，不会混淆单位")
	fmt.Println("  - bloom 库轻量，开箱即用")
}
