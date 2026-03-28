# Checklist

## 模块一：Go并发基础
- [x] goroutine vs Thread对比代码可运行
- [x] goroutine泄漏演示可复现，退出方案代码可运行
- [x] Channel基础示例可运行
- [x] select多路复用示例可运行（含default分支和超时场景）
- [x] sync.Mutex示例正确演示互斥
- [x] Mutex不可重入特性有说明和演示
- [x] 全局死锁 vs 局部死锁场景有区分说明
- [x] sync.RWMutex示例正确演示读写锁

## 模块二：atomic原子操作
- [x] atomic.Int64计数器示例可运行
- [x] atomic vs Mutex有benchmark对比数据
- [x] 说明了atomic适用场景（单一数值）vs Mutex适用场景（复合操作）

## 模块三：抢票业务实战
- [x] 无锁抢票能演示超卖问题
- [x] atomic抢票能解决超卖问题
- [x] Mutex抢票能解决超卖问题
- [x] atomic vs Mutex有性能对比
- [x] Channel抢票能正确处理并发，且注释说明了教学性质

## 模块四：sync包核心组件
- [x] WaitGroup示例正确同步
- [x] Once示例只执行一次
- [x] Pool示例正确复用对象，且说明了不适合连接池的原因
- [x] Map示例并发安全
- [x] sync.Map vs map+RWMutex有benchmark对比，说明了适用场景

## 模块五：Context超时控制
- [x] WithTimeout示例正确超时
- [x] WithCancel示例正确取消
- [x] 链式传递正确工作（父cancel传播到子goroutine）

## 模块六：后端常见业务问题
- [x] 分布式锁使用redsync库实现，非手写
- [x] 分布式锁演示了锁超时与续期
- [x] 限流器正确限流
- [x] 缓存穿透：布隆过滤器方案正确
- [x] 缓存击穿：singleflight方案正确
- [x] 缓存雪崩：随机TTL方案正确

## 文档完整性
- [x] 每个示例都有Java对比说明
- [x] 每个Java对比章节都有"Go惯用法"说明
- [x] README文档完整
