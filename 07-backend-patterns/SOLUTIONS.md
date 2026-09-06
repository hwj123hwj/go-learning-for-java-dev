# 后端业务模式 · 参考答案

## 限流

**1.1** Allow()：只有 **1 个** true。桶容量 1，瞬时 5 个请求只有一个能拿到令牌（其余直接 false，不等待——这是 Allow 的语义："现在有没有令牌"）。Wait()：5 个都会排队等令牌，按每秒 2 个的速度依次放行；500ms 超时的话，大约只有第 1~2 个能等到，其余收到 ctx 取消。选择：API 网关用 Allow 快速拒绝（返回 429），后台任务用 Wait 削峰。

**1.2** 按键限流的核心是"sync.Map + 双检"（和 `09_sync/syncmap` 呼应）：

```go
type KeyLimiter struct {
    limit rate.Limit
    burst int
    m     sync.Map // key -> *rate.Limiter
}

func (k *KeyLimiter) Allow(key string) bool {
    if v, ok := k.m.Load(key); ok {
        return v.(*rate.Limiter).Allow()
    }
    l := rate.NewLimiter(k.limit, k.burst)
    actual, loaded := k.m.LoadOrStore(key, l) // 并发安全的双检
    if loaded {                               // 别人先建了，用别人的
        l = actual.(*rate.Limiter)
    }
    return l.Allow()
}
```

生产还要考虑：key 无限增长（加 LRU/定期清理）、不同 key 不同配额（value 换成结构体）、limiter 冷启动。

**1.3**

```go
type TokenBucket struct {
    ch     chan struct{}
    refill int // 每秒补充数
}

func NewTokenBucket(capacity, refillPerSec int) *TokenBucket {
    b := &TokenBucket{ch: make(chan struct{}, capacity), refill: refillPerSec}
    for i := 0; i < capacity; i++ {
        b.ch <- struct{}{}
    }
    go func() {
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()
        for range ticker.C {
            for i := 0; i < b.refill; i++ {
                select {
                case b.ch <- struct{}{}:
                default: // 桶满了，丢弃令牌
                }
            }
        }
    }()
    return b
}

func (b *TokenBucket) Allow() bool {
    select {
    case <-b.ch:
        return true
    default:
        return false
    }
}
```

差异：本实现靠后台 goroutine 定时补（有误差、常驻 goroutine）；`x/time/rate` **不靠定时器**，每次 Allow 时用"上次剩余令牌 + 经过时间 × 速率"惰性计算，零后台开销、精度高。面试答出"惰性计算"是加分项。

## 缓存三连

**2.1**

| 问题 | 一句话定义 | 标准方案 |
|---|---|---|
| 穿透 | 查**不存在**的 key，缓存永远不命中，全部打到 DB | 布隆过滤器；或缓存空值（短 TTL） |
| 击穿 | **热点** key 过期瞬间，大量并发同时回源 | singleflight 合并回源；或热点永不过期 + 主动刷新 |
| 雪崩 | **大量** key 同时过期（或 Redis 整体挂了），DB 被冲垮 | 随机 TTL 抖动；多级缓存；集群高可用 |

记忆锚点：穿透是"数据不存在"、击穿是"一个热点"、雪崩是"一大片"。

**2.2** 布隆过滤器只预热了上线时 DB 里已有的 key。新注册用户的 key 不在过滤器里 → `Test` 返回 false → 被当成"不存在的 key"拦截。修复：**写路径同步更新**——新增用户写 DB 成功后同时 `filter.AddString(key)`。布隆过滤器**不能删**（位是共享的，删不了单 key）；要删就换计数布隆过滤器（Counting Bloom）或布谷鸟过滤器（Cuckoo Filter）。

**2.3** singleflight 的语义是"合并**调用**"，成败一起扛：50 个请求拿到**同一个 error**。对"key 不存在"这类正常结果这是对的；但"DB 挂了"时，下一批 50 个并发又会触发一次真实查询——DB 刚恢复就被打。修复方向：
- **error 也缓存短 TTL**（如 100ms~1s），期间直接快速失败，不让请求穿透到 DB；
- 或上熔断器（连续 N 次失败后直接拒绝一段时间，半开探测恢复）。
- 这正是 `sentinel error + cache-aside` 组合的实战变体。

**2.4**

```go
type SafeTTLCache struct {
    mu   sync.RWMutex
    data map[string]cacheItem
}
type cacheItem struct {
    value    string
    expireAt time.Time
}

func (c *SafeTTLCache) Get(key string) (string, bool) {
    c.mu.RLock()
    item, ok := c.data[key]
    c.mu.RUnlock()
    if !ok || time.Now().After(item.expireAt) {
        return "", false
    }
    return item.value, true
}

func (c *SafeTTLCache) Set(key, value string, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = cacheItem{value: value, expireAt: time.Now().Add(ttl)}
}
```

顺删过期项：Get 发现过期后升级写锁删除（记得 double-check，可能已被别人删了）。和 `06-concurrency` 06 题是同一模式：**读锁查、未命中升级写锁 + 双检**。

## 分布式锁

**3.1** 第 5 秒锁自动释放，B 可以拿锁——第 5~8 秒 **A、B 同时认为自己持锁**，互斥被破坏。Extend/续期（watchdog）就是解决这个：持锁期间定期把 TTL 往后推，业务完成才释放。Redisson 在 Java 里默认每 10 秒（TTL 30 秒的 1/3）自动续期，业务崩溃则续期也停止、锁自然过期——Go 的 redsync 需要你自己起 goroutine 调 `Extend`，没有内置 watchdog。

**3.2** A 执行 `DEL lockkey` 删掉的是 **B 的锁**！之后 C 又能拿锁 → 互斥彻底崩坏。标准修复：

```
加锁:  SET lockkey <随机token> NX PX 5000
解锁:  Lua 脚本（原子）:
       if redis.call("GET", KEYS[1]) == ARGV[1] then
           return redis.call("DEL", KEYS[1])
       else return 0 end
```

必须是 Lua 的原因：**GET 比对和 DEL 必须原子执行**。分两条命令发的话，"比对通过"和"DEL"之间锁可能恰好过期、别人刚拿锁——又是删别人的锁。redsync 的 `Unlock` 内部就是这段 Lua。一句话：分布式锁的正确性依赖"检查+动作"的原子性，Redis 只能用脚本保证。

**3.3** 本质局限：Redis 的主从复制是**异步**的，锁的持久性和互斥性没有强保证——它是个"高可用但非严格"的锁。RedLock：向 N 个独立 Redis 节点加锁，过半成功才算持有（不依赖主从同步）。争议（Martin Kleppmann vs antirez）：RedLock 依赖各节点时钟基本一致，进程暂停/GC 停顿仍可能导致两个客户端同时持锁，而它又没有像 fencing token 那样的下游防护——为了"较强正确性"付出 5 倍延迟，不如直接用**共识存储**。工程结论（也是你面试可以说的）：要 CP 语义用 etcd/ZooKeeper（或 Consul）的锁；用 Redis 锁就承认它是"防绝大多数并发"的性能优化，配合幂等/版本号兜底。
