package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{ch: make(chan struct{}, n)}
}

func (s *Semaphore) Wait() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Signal() {
	<-s.ch
}

func producerConsumer() {
	fmt.Println("--- 生产者-消费者问题 ---")
	const bufferSize = 5
	const itemCount = 10

	buffer := make([]int, 0, bufferSize)
	var mu sync.Mutex
	empty := NewSemaphore(bufferSize)
	full := NewSemaphore(0)
	done := make(chan bool, 2)

	go func() {
		for i := 1; i <= itemCount; i++ {
			empty.Wait()
			mu.Lock()
			buffer = append(buffer, i)
			fmt.Printf("  生产者: 生产%d, 缓冲区=%v\n", i, buffer)
			mu.Unlock()
			full.Signal()
			time.Sleep(50 * time.Millisecond)
		}
		done <- true
	}()

	go func() {
		for i := 1; i <= itemCount; i++ {
			full.Wait()
			mu.Lock()
			item := buffer[0]
			buffer = buffer[1:]
			fmt.Printf("  消费者: 消费%d, 缓冲区=%v\n", item, buffer)
			mu.Unlock()
			empty.Signal()
			time.Sleep(80 * time.Millisecond)
		}
		done <- true
	}()

	<-done
	<-done
	fmt.Println("  生产者-消费者完成!")
	fmt.Println()
}

func readersWriters() {
	fmt.Println("--- 读者-写者问题(读者优先) ---")
	var mu sync.Mutex
	var readCount int
	var writeMu sync.Mutex
	var wg sync.WaitGroup

	read := func(id int) {
		defer wg.Done()
		mu.Lock()
		readCount++
		if readCount == 1 {
			writeMu.Lock()
		}
		mu.Unlock()

		fmt.Printf("  读者%d 正在读取... (当前读者数=%d)\n", id, readCount)
		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		readCount--
		if readCount == 0 {
			writeMu.Unlock()
		}
		mu.Unlock()
		fmt.Printf("  读者%d 读取完毕\n", id)
	}

	write := func(id int) {
		defer wg.Done()
		writeMu.Lock()
		fmt.Printf("  写者%d 正在写入...\n", id)
		time.Sleep(150 * time.Millisecond)
		fmt.Printf("  写者%d 写入完毕\n", id)
		writeMu.Unlock()
	}

	wg.Add(5)
	go read(1)
	go read(2)
	go write(1)
	go read(3)
	go write(2)
	wg.Wait()
	fmt.Println()
}

func diningPhilosophers() {
	fmt.Println("--- 哲学家进餐问题 ---")
	const n = 5
	forks := make([]sync.Mutex, n)
	philosophers := []string{"孔子", "老子", "庄子", "孟子", "荀子"}

	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			left := id
			right := (id + 1) % n

			if id == n-1 {
				left, right = right, left
			}

			forks[left].Lock()
			fmt.Printf("  %s 拿起左边筷子%d\n", philosophers[id], left)
			time.Sleep(50 * time.Millisecond)

			forks[right].Lock()
			fmt.Printf("  %s 拿起右边筷子%d -> 开始进餐\n", philosophers[id], right)
			time.Sleep(100 * time.Millisecond)

			forks[right].Unlock()
			forks[left].Unlock()
			fmt.Printf("  %s 放下筷子，进餐完毕\n", philosophers[id])
		}(i)
	}
	wg.Wait()
	fmt.Println()
}

func barberShop() {
	fmt.Println("--- 理发师问题 ---")
	const chairs = 3
	var mu sync.Mutex
	waiting := 0
	customerReady := make(chan struct{}, 1)
	barberReady := make(chan struct{}, 1)
	var wg sync.WaitGroup

	go func() {
		for {
			<-customerReady
			mu.Lock()
			waiting--
			mu.Unlock()
			fmt.Println("  理发师: 开始理发")
			time.Sleep(200 * time.Millisecond)
			fmt.Println("  理发师: 理发完毕")
			barberReady <- struct{}{}
		}
	}()

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mu.Lock()
			if waiting < chairs {
				waiting++
				fmt.Printf("  顾客%d: 坐下等待 (等待人数=%d)\n", id, waiting)
				mu.Unlock()
				customerReady <- struct{}{}
				<-barberReady
				fmt.Printf("  顾客%d: 理发完毕离开\n", id)
			} else {
				fmt.Printf("  顾客%d: 没有空位，离开\n", id)
				mu.Unlock()
			}
		}(i)
		time.Sleep(80 * time.Millisecond)
	}
	wg.Wait()
	fmt.Println()
}

func main() {
	fmt.Println("=== 经典同步问题 ===")
	fmt.Println()

	producerConsumer()
	readersWriters()
	diningPhilosophers()
	barberShop()

	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. 临界区: 访问共享资源的代码段，需互斥进入")
	fmt.Println("2. 信号量: P操作(wait)减1，V操作(signal)加1")
	fmt.Println("3. 生产者-消费者: 同步(满/空) + 互斥(缓冲区)")
	fmt.Println("4. 读者-写者: 读者可同时读，写者必须独占")
	fmt.Println("   - 读者优先: 可能写者饥饿")
	fmt.Println("   - 写者优先: 读者可能饥饿")
	fmt.Println("   - 读写公平: 按请求顺序")
	fmt.Println("5. 哲学家进餐: 破坏循环等待(奇偶编号拿筷子顺序不同)")
	fmt.Println("6. 理发师问题: 有限等待+服务唤醒")
	fmt.Println("7. 管程(Monitor): 封装共享数据+条件变量，比信号量更安全")
	fmt.Println("8. Go的sync.Mutex对应二值信号量，sync.RWMutex对应读写锁")
}
