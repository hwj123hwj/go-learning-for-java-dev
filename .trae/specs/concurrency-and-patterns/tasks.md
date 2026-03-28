# Tasks

## 模块一：Go并发基础（对比Java）

- [x] Task 1: 创建concurrency-demo项目结构
  - [x] 创建目录和go.mod
  - [x] 创建README.md说明文档

- [x] Task 2: goroutine vs Thread对比
  - [x] 创建goroutine基础示例
  - [x] 对比Java Thread创建方式
  - [x] 演示goroutine轻量级特性（创建百万goroutine vs 线程内存占用）

- [x] Task 3: goroutine泄漏演示与防范
  - [x] 演示goroutine因channel阻塞导致泄漏的场景
  - [x] 演示通过close channel正确退出
  - [x] 演示通过context取消正确退出
  - [x] 说明pprof如何排查goroutine泄漏

- [x] Task 4: Channel基础
  - [x] 创建无缓冲channel示例
  - [x] 创建有缓冲channel示例
  - [x] 对比Java BlockingQueue

- [x] Task 5: select多路复用
  - [x] 创建监听多个channel的select示例
  - [x] 创建带default分支的非阻塞select示例
  - [x] 创建select + time.After实现超时的示例
  - [x] 说明select与Java wait/notify的对比

- [x] Task 6: sync.Mutex vs ReentrantLock
  - [x] 创建互斥锁示例
  - [x] 对比Java synchronized和ReentrantLock
  - [x] 演示Go Mutex不可重入特性（区别于Java ReentrantLock）
  - [x] 演示全局死锁（运行时可检测）vs 局部死锁（运行时无法检测）

- [x] Task 7: sync.RWMutex vs ReentrantReadWriteLock
  - [x] 创建读写锁示例
  - [x] 演示读多写少场景性能优化
  - [x] 对比Java读写锁

## 模块二：atomic原子操作

- [x] Task 8: sync/atomic基础
  - [x] 创建atomic.Int64计数器示例
  - [x] 对比Java AtomicInteger/AtomicLong
  - [x] 演示atomic vs Mutex的性能差异（benchmark）
  - [x] 说明适用场景：单一数值读写用atomic，复合操作用Mutex

## 模块三：抢票业务实战

- [x] Task 9: 抢票问题演示（无锁）
  - [x] 创建票务结构
  - [x] 多goroutine并发抢票
  - [x] 演示超卖问题

- [x] Task 10: 抢票解决方案一（atomic）
  - [x] 使用atomic.Int64原子递减库存
  - [x] 验证不会超卖
  - [x] 与Mutex方案做benchmark对比

- [x] Task 11: 抢票解决方案二（Mutex）
  - [x] 使用sync.Mutex保护库存
  - [x] 验证不会超卖
  - [x] 对比Java synchronized解决方案

- [x] Task 12: 抢票解决方案三（Channel，教学向）
  - [x] 使用channel串行化请求
  - [x] 演示Go特有并发模式
  - [x] 代码注释中说明：此方案教学目的为主，生产中原子操作通常更优

## 模块四：sync包核心组件

- [x] Task 13: sync.WaitGroup
  - [x] 创建WaitGroup示例
  - [x] 对比Java CountDownLatch
  - [x] 演示批量等待场景

- [x] Task 14: sync.Once
  - [x] 创建单例初始化示例
  - [x] 对比Java双检锁单例
  - [x] 演示线程安全初始化

- [x] Task 15: sync.Pool
  - [x] 创建对象池示例
  - [x] 演示GC优化效果
  - [x] 说明Pool中对象可能被GC回收，不适合连接池场景

- [x] Task 16: sync.Map
  - [x] 创建并发安全Map示例
  - [x] 对比Java ConcurrentHashMap
  - [x] 对比map+RWMutex，说明sync.Map适用场景（读多写少、key稳定）
  - [x] benchmark演示写多场景下sync.Map性能劣势

## 模块五：Context超时控制

- [x] Task 17: Context基础
  - [x] 创建context.WithTimeout示例
  - [x] 创建context.WithCancel示例
  - [x] 创建context.WithDeadline示例

- [x] Task 18: Context实战
  - [x] 模拟HTTP请求超时
  - [x] 模拟数据库查询超时
  - [x] 链式传递context（父cancel自动传播到子goroutine）

## 模块六：后端常见业务问题

- [x] Task 19: 创建backend-patterns项目结构
  - [x] 创建目录和go.mod
  - [x] 创建README.md说明文档

- [x] Task 20: 分布式锁
  - [x] 引入go-redis和redsync库
  - [x] 使用redsync实现Redis分布式锁
  - [x] 演示锁超时与自动续期
  - [x] 对比Java Redisson，说明为何不从头手写

- [x] Task 21: 限流器
  - [x] 使用golang.org/x/time/rate实现令牌桶限流
  - [x] 实现漏桶限流
  - [x] 对比Java Guava RateLimiter

- [x] Task 22: 缓存问题解决
  - [x] 缓存穿透：布隆过滤器方案
  - [x] 缓存击穿：golang.org/x/sync/singleflight方案
  - [x] 缓存雪崩：随机TTL + 降级方案

# Task Dependencies
- Task 3 依赖 Task 2（泄漏需先理解goroutine基础）
- Task 5 依赖 Task 4（select需先理解channel）
- Task 9-12 依赖 Task 6、Task 8（抢票需先理解锁和atomic）
- Task 13-16 依赖 Task 2-5（sync包需先理解基础并发）
- Task 17-18 依赖 Task 2（Context需理解goroutine）
- Task 20-22 依赖 Task 6（业务问题需理解锁）
