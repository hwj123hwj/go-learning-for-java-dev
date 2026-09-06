# 并发十讲 · 参考答案

> 先做题，再对答案。答案只给关键点和核心代码，跑不通的题回源码找原因。

## 01 goroutine 基础

**1.1** 大概率只输出 `main done`。main 结束 = 整个进程退出，不会等 goroutine（和 Java 的非 daemon 线程相反，和 daemon 线程一致）。这就是 `01_goroutine` 里必须用 WaitGroup 的原因。

**1.2** 两处错误：
1. `wg.Add(1)` 写在了 goroutine **内部** —— main 可能先跑到 `wg.Wait()`，此时计数器还是 0，Wait 直接返回，worker 是否执行全看运气。`Add` 必须在 `go` 之前调用。
2. 循环变量捕获：本仓库 go.mod 声明 go 1.22+，`i` 每次迭代独立，能正确打印 0/1/2；但在 go 1.21 及以前，闭包捕获的是同一个 `i`，通常打印三个 `3`。面试常考这个分界线。

**1.3** 用"完成信号 channel"：

```go
done := make(chan struct{})
for i := 0; i < 3; i++ {
    go func(id int) {
        fmt.Println("worker", id)
        done <- struct{}{}
    }(i)
}
for i := 0; i < 3; i++ {
    <-done
}
```

要点：`struct{}` 零大小，只当信号用；收 N 次。WaitGroup 版更好读，channel 版是理解 pipeline 的基础——两种都要会。

## 02 goroutine 泄漏

**2.1** `f2` 泄漏。goroutine C 卡在 `<-ch`：无缓冲 channel 没有接收者也没有 close，永远阻塞。`f1` 正常（A 发送时有 B 在接收）。判定口诀：**goroutine 阻塞在一个永远不会有人再碰的 channel 上 = 泄漏**。

**2.2**

```go
// channel 版
stop := make(chan struct{})
go func() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-stop:
            return
        case <-ticker.C:
            fmt.Println("heartbeat")
        }
    }
}()
close(stop) // 停止

// context 版
go func(ctx context.Context) {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            fmt.Println("heartbeat")
        }
    }
}(ctx)
```

推荐 context：取消信号可以沿调用链**传播**（父取消 → 所有子一起退出），还能组合超时，是 Go 的标准做法；裸 stop channel 只适合单点控制。

**2.3** 不会阻止进程退出（进程退出时所有 goroutine 一起销毁），但长期运行的服务里会累积：占内存、持有引用导致对象无法 GC、可能占满连接。`runtime.NumGoroutine()` 持续上涨是泄漏信号；定位用 pprof：`import _ "net/http/pprof"` 后 `go tool pprof http://localhost:port/debug/pprof/goroutine`，直接看到每个 goroutine 卡在哪一行。

## 03 channel

**3.1** A：main 阻在 `ch <- 1`（无缓冲发送需要就绪的接收者）→ 全部 goroutine 睡眠 → `fatal error: all goroutines are asleep - deadlock!`。B：正常打印 1。C：缓冲 2 装下 1、2 后，第三次发送阻塞 → 同样 deadlock 崩溃。

**3.2** 发送到已关闭 channel：**panic**。从已关闭 channel 接收：立即返回缓冲区剩余数据，取完后返回零值 + `ok=false`，永不阻塞。重复 close：**panic**。规则：**close 由发送方负责，且只 close 一次**；接收方用 `v, ok := <-ch` 判断。

**3.3**

```go
sem := make(chan struct{}, 3)
var wg sync.WaitGroup
start := time.Now()
for i := 1; i <= 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        sem <- struct{}{}         // 满了就阻塞 = 等信号量
        defer func() { <-sem }()  // 释放
        time.Sleep(300 * time.Millisecond)
        fmt.Println("task", id)
    }(i)
}
wg.Wait()
fmt.Println(time.Since(start)) // ≈ 1.2s（4 批 × 300ms），而不是 3s
```

## 04 select

**4.1** 输出 a/b 的随机序列。多个 case 同时 ready 时 select **等概率随机**选一个——这是刻意设计，避免某个 case 饥饿。所以不能依赖 select 的"顺序"。

**4.2**

```go
select {
case v := <-result:
    fmt.Println("结果:", v)
case <-time.After(1 * time.Second):
    fmt.Println("超时")
}
```

**4.3**

```go
func merge(ch1, ch2 chan int) chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(2)
    for _, ch := range []chan int{ch1, ch2} {
        go func(c chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }
    go func() { wg.Wait(); close(out) }() // 两个都读完才关
    return out
}
```

关键：**关闭 out 的职责交给独立 goroutine**，等两个转发 goroutine 都退出后再 close。这是 fan-in 的标准模板。

## 05 Mutex

**5.1** 0~20000 之间的任意值（通常略小于 20000）。`counter++` 是读-改-写三步，两个 goroutine 可能都读到同一个旧值，各加一次只写回一份 → 丢失更新。这就是 `08_ticket` 超卖的微观原理。

**5.2** `DeductWithFee` 死锁：持锁状态下调 `Deduct` 再次 `Lock` 同一把锁，Go 的 Mutex **不可重入**（不记录持有者），永久阻塞。Java 的 synchronized/ObjectMonitor 记录持有线程所以可重入。Go 的设计哲学：隐式可重入让锁的边界不可见，容易掩盖设计问题；Go 要求你把"不需要锁的内部版本"显式拆出来（如 `deductLocked` 私有方法）。

**5.3** 经典 ABBA 死锁：g1 持 a 等 b，g2 持 b 等 a，永远互等。修复：**全局统一的加锁顺序**——按地址排序：

```go
func (a *Account) Transfer(b *Account, amount int) {
    first, second := a, b
    if &a.balance > &b.balance { first, second = b, a } // 指针比较定序
    first.mu.Lock()
    defer first.mu.Unlock()
    second.mu.Lock()
    defer second.mu.Unlock()
    // ...按方向转账
}
```

实际工程更常用：按 ID 排序，或者干脆一把全局锁/串行化通道。

## 06 RWMutex

**6.1** 读多写少（如配置、缓存元数据）时 RWMutex 让读并发、明显快于 Mutex。但临界区很短或写频繁时更慢：RWMutex 本身状态更复杂（要维护读者计数、写者排队），每次 Lock/Unlock 开销更大。结论：**先用 Mutex，benchmark 证明读是瓶颈再换 RWMutex**。

**6.2** 死锁。RLock 未释放时调 Set 要拿写锁；写锁要等所有读锁释放；而自己这个"读锁"又在等 Set 返回——自我死锁。RWMutex 没有锁升级。修复：GetOrLoad 不要在读锁里做加载——先 RLock 查，未命中则释放读锁、拿写锁、**double-check** 再加载（和 Java 双检锁同构）。

**6.3**

```go
type Config struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Config) Set(key, val string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = val
}

func (c *Config) Get(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key] // string 是不可变值，返回即副本
}
```

陷阱：如果存的是 map/slice/指针，Get 返回的引用仍可能被并发写坏——要么存值类型，要么返回深拷贝。

## 07 atomic

**7.1** `Add(1)` 是原子的读-改-写（底层 CAS 循环），任何时刻不会丢更新。`x = x.Load() + 1` 分两步：Load 和 Store 之间可能被插入——两个 goroutine 都 Load 到 5，都 Store 6，丢一次。错误不在 Load 也不在 Store，在**两步之间的间隙**。

**7.2** `07_atomic` 的 benchmark 里单变量计数 atomic 比 Mutex 快数倍（无互斥队列、无调度开销）。但 `{count, total}` 要**同时一致**：atomic 只能保证单个变量，两个 atomic 之间仍会交错（读到的 count 和 total 可能是"半新半旧"）。必须 Mutex 保护复合状态，或把整个结构体当快照存进一个 `atomic.Pointer`。

**7.3**

```go
type Config struct { RPS int; Timeout time.Duration }

var cfgPtr atomic.Pointer[Config]

func UpdateConfig(c *Config) { cfgPtr.Store(c) }          // 整体替换
func CurrentConfig() *Config  { return cfgPtr.Load() }    // 无锁读，快照永远一致
```

写侧每次构造**新的完整对象**再 Store（copy-on-write），读侧 Load 拿到的快照不可变——这就是配置热更新、你网关里路由表热加载的标准做法。

## 08 抢票实战

**8.1** `no_lock` 的 Buy 核心是：

```go
if t.remaining > 0 {   // ① 读
    t.remaining--      // ② 写（本身又分读、减、写）
    return true
}
```

g1、g2 都在 remaining=1 时执行到 ① 之间交错：g1 读到 1 通过检查 → 调度 → g2 读到 1（还没人减）也通过检查 → 两人都执行 ② → 1 张票卖出 2 张。

**8.2** CAS 版：不加锁不阻塞，低竞争时最快；竞争激烈时大量自旋空转烧 CPU，且只保护单个变量。Mutex 版：竞争时让等待者休眠（不空转），适合高竞争；但每次 Lock/Unlock 有固定开销。抢票这种瞬时高竞争场景，Mutex/Channel 更稳；CAS 适合低竞争的简单计数——这也是 `07_atomic` benchmark 想让你建立的直觉。

**8.3**

```go
pool := make(chan struct{}, 100) // 100 张票
for i := 0; i < 100; i++ { pool <- struct{}{} }

func Buy() bool {
    select {
    case <-pool:  // 取到一张
        return true
    default:      // 池子空了，不阻塞
        return false
    }
}
```

语义差异：channel 版"票在池子里"，售罄是天然语义；`channel_lock` 版是拿 channel 当锁用（保护共享计数），两种用法都值得吃透。

## 09 sync 包

**9.1** **不会执行**。验证：

```go
var once sync.Once
func() {
    defer func() { recover() }()
    once.Do(func() { panic("boom") })
}()
once.Do(func() { fmt.Println("second") }) // 不打印
```

原因：Once 实现里 `done` 标记是用 `defer o.done.Store(1)` 设置的，f panic 时 defer 照样执行，Once 已标记完成。所以 Once 的语义是"最多调用一次"，**不保证成功一次"。

**9.2** Pool 在每次 GC 时会被清空（对象随时消失），也不限制总量、没有 Close 生命周期。数据库连接昂贵、必须全局有界、要显式关闭——这三点 Pool 全不满足。Pool 的正确用途：复用**廉价、可丢弃**的临时对象，典型如 `bytes.Buffer`（标准库 fmt/print 就这么干）。

**9.3** 一句话：**key 集合稳定、读多写少（一次写入多次读）或各 goroutine 操作不相交 key 时 sync.Map 占优；写多、key 频繁增删时反而比 map+RWMutex 慢**。benchmark 里 sync.Map 并非全面胜出，读多场景才拉开差距。

## 10 Context

**10.1** child 在 **1 秒** 后 Done，`ctx.Err() == context.DeadlineExceeded`。child 写 200ms 则 200ms 到期。规则：**子 context 的截止时间 = min(自己的, 父链上所有的)**，且父取消无条件传播给所有子孙——子永远活不过父。

**10.2** 两处问题：
1. `cancel` 创建了却从没调用（缺 `defer cancel()`）——timer 和 context 树节点在超时前一直不释放，`go vet` 会报 `lostcancel`。
2. 更严重：`http.Get` **根本不接 context**！WithTimeout 建了个 ctx 却没用它，超时完全不生效。必须 `http.NewRequestWithContext(ctx, ...)` + `client.Do(req)`。

**10.3**

```go
func (s *Service) Fetch(ctx context.Context, id string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second) // 超时叠加在调用方 ctx 之上
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, http.MethodGet,
        "https://api.example.com/"+id, nil)
    if err != nil {
        return nil, fmt.Errorf("build request: %w", err)
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("call upstream: %w", err) // %w 保留错误链
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}
```

三个要点：ctx 由调用方传入（函数只叠加自己的超时）；`defer cancel()` 立刻跟上；错误用 `%w` 包装。这就是你网关里每个上游调用的形状。
