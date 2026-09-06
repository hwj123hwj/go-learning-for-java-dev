# 并发十讲 · 练习题

> 用法：遮住源码和答案做题。四类题型：**预测**（写出程序输出）、**改造**（改既有代码）、**徒手**（不看源码从零实现）、**找 bug**（指出问题并修复）。
> 做完对 [SOLUTIONS.md](./SOLUTIONS.md)。每题都做对 = 这个模块真的吃透了。

---

## 01 goroutine 基础

**1.1〔预测〕** 下面程序输出什么？为什么？

```go
func main() {
    go fmt.Println("hello from goroutine")
    fmt.Println("main done")
}
```

**1.2〔找 bug〕** 下面代码想等 3 个 goroutine 全部完成，有两处错误，找出来：

```go
var wg sync.WaitGroup

for i := 0; i < 3; i++ {
    go func() {
        wg.Add(1)
        defer wg.Done()
        fmt.Println("worker", i)
    }()
}
wg.Wait()
```

**1.3〔徒手〕** 不用 `sync.WaitGroup`，只用 channel 实现"等 3 个 goroutine 全部完成"。写出完整可运行的 main。

---

## 02 goroutine 泄漏

**2.1〔找 bug〕** 下面哪个 goroutine 泄漏了？指出具体卡在哪一行，并说明为什么：

```go
func f1() {
    ch := make(chan int)
    go func() { ch <- 1 }()          // A
    <-ch                              // B
}

func f2() {
    ch := make(chan int)
    go func() { <-ch }()              // C
    // 函数直接返回，没有发送、没有 close    // D
}
```

**2.2〔改造〕** 写一个 `worker` goroutine：每秒打印一次心跳，直到收到停止信号退出。给出 channel 版和 context 版两种写法，并说明哪种更推荐、为什么。

**2.3〔思考〕** 泄漏的 goroutine 会阻止进程退出吗？`runtime.NumGoroutine()` 能不能用来验证泄漏？怎么在真实项目里定位泄漏点？

---

## 03 channel

**3.1〔预测〕** 三段代码分别发生什么（正常输出 / 卡住 / panic）？

```go
// A
ch := make(chan int)
ch <- 1
fmt.Println(<-ch)

// B
ch := make(chan int, 1)
ch <- 1
fmt.Println(<-ch)

// C
ch := make(chan int, 2)
ch <- 1; ch <- 2; ch <- 3
fmt.Println(len(ch))
```

**3.2〔预测〕** 向已关闭的 channel 发送数据、从已关闭的 channel 接收数据、重复 close，分别发生什么？

**3.3〔徒手〕** 用"带缓冲 channel 当信号量"实现：10 个并发任务，同一时刻最多 3 个在执行。写出完整代码，要求结束时打印总耗时（应该约为 10/3 个任务的时长）。

---

## 04 select

**4.1〔预测〕** 两个 case 都 ready 时 select 选哪个？

```go
ch1 := make(chan string, 1)
ch2 := make(chan string, 1)
ch1 <- "a"; ch2 <- "b"
for i := 0; i < 4; i++ {
    select {
    case v := <-ch1: fmt.Println(v); ch1 <- "a"
    case v := <-ch2: fmt.Println(v); ch2 <- "b"
    }
}
```

**4.2〔改造〕** 下面代码用轮询 + sleep 实现"等结果或超时"，改写成 select 写法（不许用 time.Sleep）：

```go
result := make(chan string)
go func() { time.Sleep(2 * time.Second); result <- "done" }()

for {
    select {} // ← 你的实现
}
```

**4.3〔徒手〕** 实现 `merge(ch1, ch2 chan int) chan int`：把两个 channel 的数据汇入一个返回的 channel，**两个输入都关闭后**返回的 channel 才关闭。这是 fan-in 模式，是 pipeline 的基础。

---

## 05 Mutex

**5.1〔预测〕** 两个 goroutine 各自对共享 `int` 做 10000 次 `counter++`（无锁），最终 counter 的可能范围是多少？为什么不是 20000？

**5.2〔找 bug〕** Go 的 `sync.Mutex` 不可重入。下面代码运行时发生什么？如果这是 Java，同样的逻辑（synchronized 方法互调）是正常的——说明 Go 为什么这样设计：

```go
func (a *Account) Deduct(amount int) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.balance -= amount
}

func (a *Account) DeductWithFee(amount, fee int) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.Deduct(amount) // 同一把锁，第二次 Lock
    a.Deduct(fee)
}
```

**5.3〔找 bug〕** 下面是更隐蔽的死锁：两个账户互相转账。两个 goroutine 各自持有一把锁再去要另一把，怎么修复？（提示：全局锁排序）

```go
func (a *Account) Transfer(b *Account, amount int) {
    a.mu.Lock()
    defer a.mu.Unlock()
    b.mu.Lock()
    defer b.mu.Unlock()
    a.balance -= amount
    b.balance += amount
}
// goroutine 1: a.Transfer(b, 100)  → 先锁 a 再锁 b
// goroutine 2: b.Transfer(a, 50)   → 先锁 b 再锁 a
```

---

## 06 RWMutex

**6.1〔思考〕** 什么场景 RWMutex 比 Mutex 快？什么场景反而更慢？`06_rwlock` 里的性能对比数据说明了什么？

**6.2〔找 bug〕** 下面的 `Config` 想支持读锁里触发写（锁升级），会发生什么？

```go
func (c *Config) GetOrLoad(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    if v, ok := c.data[key]; ok { return v }
    v := loadFromRemote(key)      // 内部会调用 c.Set()
    return v
}
```

**6.3〔徒手〕** 用 `sync.RWMutex` 实现线程安全的配置表 `Set(key, val)` / `Get(key) string`，要求 `Get` 返回的是**副本**而不是可能被并发写坏的共享值。

---

## 07 atomic

**7.1〔预测〕** `atomic.Int64` 的 `Add(1)` 和 `x = x.Load() + 1` 有什么区别？后者错在哪一步？

**7.2〔思考〕** `07_atomic` 的 benchmark 里 atomic 比 Mutex 快多少？什么情况必须用 Mutex（atomic 做不到）？给一个具体例子：一个结构体有 `count` 和 `total` 两个字段要保持一致。

**7.3〔徒手〕** 用 `atomic.Value`（或泛型 `atomic.Pointer[T]`）实现"配置快照"：写侧整体替换配置对象，读侧无锁拿到一致快照。这就是你网关里配置热更新的原理。

---

## 08 抢票实战（超卖问题）

**8.1〔预测〕** `no_lock` 版在"剩余 1 张票、2 个 goroutine 同时 Buy()"时，哪两行代码之间被调度打断会导致超卖？把那次交错写出来。

**8.2〔对比〕** `atomic_lock`（CAS 循环）和 `mutex_lock` 各自的优缺点？抢票场景你会选哪个，为什么？（提示：CAS 失败会自旋重试，竞争激烈时 CPU 空转）

**8.3〔徒手〕** 用"带缓冲 channel 初始化 N 张票"实现票池：`Buy()` 从 channel 取票，取不到即售罄。写出完整代码，和 `channel_lock` 版对比语义差异。

---

## 09 sync 包（WaitGroup / Once / Pool / Map）

**9.1〔预测〕** `sync.Once.Do(f)` 里 f panic 了，程序恢复后第二次调用 `Do(g)` 还会执行 g 吗？先猜，再写代码验证。

**9.2〔思考〕** 为什么 `sync.Pool` 不适合当数据库连接池？（提示：Pool 会被 GC 清空。那 Pool 的正确用途是什么？）

**9.3〔思考〕** `sync.Map` 相比 `map + RWMutex`，什么访问模式下才占优势？`syncmap` 的 benchmark 里它是全面胜出吗？把结论写成一句话，面试直接用。

---

## 10 Context

**10.1〔预测〕**

```go
parent, _ := context.WithTimeout(ctx, 1*time.Second)  // 1s
child, _ := context.WithTimeout(parent, 2*time.Second) // 2s
```

child 实际多久后 Done？`ctx.Err()` 是什么？如果 child 写成 200ms 呢？总结一条规则。

**10.2〔找 bug〕** 下面代码两处问题（一处泄漏、一处 API 设计错误）：

```go
func (s *Service) Fetch(id string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    resp, err := http.Get("https://api.example.com/" + id)
    if err != nil { return "", err }
    defer resp.Body.Close()
    data, err := io.ReadAll(resp.Body)
    return string(data), err
}
```

**10.3〔改造〕** 把上面的 `Fetch` 改成：context 由**调用方**传入，超时加在请求上（`http.NewRequestWithContext`），cancel 正确释放，返回 `[]byte`。这是你网关里每个上游请求的标准写法。
