# Go并发编程与业务实战 Spec

## Why

帮助Java开发者深入理解Go并发编程的核心概念，通过对比Java的并发机制（synchronized、ReentrantLock、线程池等），结合实际业务场景（抢票、秒杀等），掌握Go并发编程的最佳实践。

> **注意**：Java对比仅作为理解入口，每个章节需明确指出"Go惯用法"，避免把Go工具强行套用Java思维。

## What Changes

### 新增模块
- `concurrency-demo/` - Go并发编程实战
  - 对比Java synchronized/ReentrantLock
  - goroutine + channel + select + sync包 + atomic包
  - 抢票/秒杀业务场景

- `backend-patterns/` - 后端常见业务问题
  - 分布式锁（基于redsync库）
  - 限流器
  - 缓存穿透/击穿/雪崩

## Impact

- 学习路径：quickstart → user-api → concurrency-demo → backend-patterns
- 目标：掌握Go并发编程和常见后端业务解决方案

---

## ADDED Requirements

### Requirement: 并发基础对比
系统应提供Go并发基础与Java对比的学习材料。

#### Scenario: goroutine vs Thread
- **WHEN** 学习并发基础
- **THEN** 理解goroutine与Java Thread的区别（调度方式、内存占用、创建成本）
- **Go惯用法**：goroutine非常轻量，可大量创建，但必须注意goroutine泄漏问题

#### Scenario: goroutine泄漏
- **WHEN** goroutine因channel阻塞或逻辑未退出而永久挂起
- **THEN** 演示泄漏场景，及通过context/close channel正确退出的方式
- **对比**：Java线程泄漏同样存在，但Go工具链（pprof）更易排查

### Requirement: Channel与select
系统应提供channel基础及select多路复用的学习材料。

#### Scenario: channel基础
- **WHEN** 学习channel
- **THEN** 理解无缓冲/有缓冲channel，对比Java BlockingQueue

#### Scenario: select多路复用
- **WHEN** 需要同时监听多个channel
- **THEN** 使用select处理多路channel，包含default分支（非阻塞）和超时场景
- **Go惯用法**：select是Go并发的核心，相当于Java中wait/notify + 条件判断的组合

### Requirement: 锁机制对比
系统应提供Go锁机制与Java对比的实战代码。

#### Scenario: sync.Mutex vs ReentrantLock
- **WHEN** 需要互斥访问共享资源
- **THEN** 使用sync.Mutex实现，对比Java ReentrantLock
- **注意**：Go的Mutex不可重入，与Java ReentrantLock不同，演示中需体现此差异

#### Scenario: sync.RWMutex vs ReentrantReadWriteLock
- **WHEN** 读多写少场景
- **THEN** 使用sync.RWMutex优化性能

#### Scenario: 死锁场景
- **WHEN** 演示死锁
- **THEN** 区分两种死锁：全局死锁（Go运行时可检测并panic）与局部死锁（部分goroutine阻塞，运行时无法检测）
- **Go惯用法**：用context + select设置超时避免死锁，而非依赖运行时检测

### Requirement: sync/atomic原子操作
系统应演示atomic包的使用，作为轻量级并发原语。

#### Scenario: atomic vs Mutex
- **WHEN** 对单个数值进行并发读写（如计数器、标志位）
- **THEN** 使用sync/atomic代替Mutex，性能更优
- **对比Java**：对应AtomicInteger、AtomicLong等
- **Go惯用法**：Go 1.19+推荐使用atomic.Int64、atomic.Bool等类型安全的封装

### Requirement: 抢票业务实战
系统应实现抢票业务，演示并发问题的解决方案。

#### Scenario: 无锁抢票（演示问题）
- **WHEN** 多个goroutine同时抢票
- **THEN** 演示超卖问题

#### Scenario: atomic抢票（轻量方案）
- **WHEN** 票务库存是单一数值的原子递减
- **THEN** 使用atomic.Int64实现，对比Mutex方案的性能差异

#### Scenario: 加锁抢票（解决问题）
- **WHEN** 使用sync.Mutex保护库存
- **THEN** 正确处理并发，避免超卖

#### Scenario: Channel抢票（Go特色，教学向）
- **WHEN** 使用channel串行化请求
- **THEN** 演示Go特有的并发模式
- **注意**：此方案为教学目的，不代表性能优于Mutex方案，需在代码注释中说明适用场景

### Requirement: sync包核心组件
系统应演示sync包的核心组件使用。

#### Scenario: sync.WaitGroup
- **WHEN** 需要等待多个goroutine完成
- **THEN** 使用WaitGroup同步，对比Java CountDownLatch

#### Scenario: sync.Once
- **WHEN** 需要单例初始化
- **THEN** 使用sync.Once保证只执行一次，对比Java双检锁单例

#### Scenario: sync.Pool
- **WHEN** 需要对象复用
- **THEN** 使用sync.Pool减少GC压力
- **注意**：sync.Pool中的对象可能被GC回收，不适合需要持久化的连接池场景

#### Scenario: sync.Map
- **WHEN** 读多写少、key相对稳定的并发Map场景
- **THEN** 使用sync.Map，对比Java ConcurrentHashMap
- **重要**：sync.Map不是ConcurrentHashMap的直接替代，写多读少场景性能反而差于`map+RWMutex`，需在示例中对比说明

### Requirement: Context超时控制
系统应演示Context的使用场景。

#### Scenario: 请求超时控制
- **WHEN** 需要控制请求超时
- **THEN** 使用context.WithTimeout

#### Scenario: 请求取消
- **WHEN** 需要取消正在进行的操作
- **THEN** 使用context.WithCancel

#### Scenario: Context链式传递
- **WHEN** goroutine树中需要统一取消
- **THEN** 演示父context取消自动传播到子context

### Requirement: 后端常见业务问题
系统应提供常见后端业务问题的解决方案。

#### Scenario: 分布式锁
- **WHEN** 多实例部署需要互斥
- **THEN** 使用go-redis + redsync库实现Redis分布式锁，演示锁超时、自动续期
- **注意**：不从头手写分布式锁，直接使用成熟库，对比Java Redisson

#### Scenario: 限流器
- **WHEN** 需要保护系统不被压垮
- **THEN** 实现令牌桶/漏桶限流，可参考golang.org/x/time/rate

#### Scenario: 缓存问题
- **WHEN** 遇到缓存穿透/击穿/雪崩
- **THEN** 提供对应解决方案（布隆过滤器、singleflight、随机TTL）
