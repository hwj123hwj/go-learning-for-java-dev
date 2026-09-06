# Go并发编程实战（Java开发者版）

> 通过对比Java并发机制，深入理解Go并发编程的核心概念和最佳实践。

## ⚠️ 重要提示

**Java对比仅作为理解入口**，每个章节都会明确指出"Go惯用法"，避免把Go工具强行套用Java思维。

## 目录结构

```
concurrency-demo/
├── 01_goroutine/           # goroutine vs Thread
├── 02_goroutine_leak/      # goroutine泄漏演示
├── 03_channel/             # Channel基础
├── 04_select/              # select多路复用
├── 05_mutex/               # sync.Mutex vs ReentrantLock
├── 06_rwlock/              # sync.RWMutex
├── 07_atomic/              # atomic原子操作
├── 08_ticket/              # 抢票业务实战
├── 09_sync/                # sync包核心组件
└── 10_context/             # Context超时控制
```

## 快速开始

```bash
cd concurrency-demo
go run ./01_goroutine/main.go
```

## 核心对比速查

| Java | Go | 说明 |
|------|-----|------|
| Thread | goroutine | goroutine更轻量（KB vs MB） |
| BlockingQueue | channel | channel是Go并发核心 |
| wait/notify | select | select更优雅处理多路IO |
| synchronized | sync.Mutex | Mutex不可重入！ |
| ReentrantLock | sync.Mutex | Mutex不可重入！ |
| ReentrantReadWriteLock | sync.RWMutex | 读多写少场景 |
| AtomicInteger | atomic.Int64 | Go 1.19+推荐类型安全封装 |
| CountDownLatch | sync.WaitGroup | 等待一组goroutine |
| 双检锁单例 | sync.Once | Once是Go单例标准做法 |
| 对象池 | sync.Pool | 临时对象池，不适合连接池 |
| ConcurrentHashMap | sync.Map | 适用场景有限！ |

## 学习路径

1. **goroutine基础** - 理解轻量级并发
2. **goroutine泄漏** - 掌握正确退出方式
3. **Channel** - 理解"用通信共享内存"
4. **select** - 多路复用核心
5. **Mutex** - 互斥锁（注意不可重入）
6. **atomic** - 轻量级原子操作
7. **抢票实战** - 综合应用
8. **sync包** - WaitGroup/Once/Pool/Map
9. **Context** - 超时控制

## Go惯用法原则

1. **不要通过共享内存来通信，要通过通信来共享内存** - channel优先
2. **goroutine泄漏比内存泄漏更难发现** - 必须有退出机制
3. **Mutex不可重入** - 与Java最大的区别
4. **atomic适合单一数值，Mutex适合复合操作**
5. **sync.Map适用场景有限** - 大多数情况用map+Mutex
6. **Context作为第一个参数传递** - 不要存储在struct中
