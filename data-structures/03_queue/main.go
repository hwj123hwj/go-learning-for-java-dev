package main

import "fmt"

type Queue struct {
	data []int
}

func NewQueue() *Queue {
	return &Queue{data: make([]int, 0)}
}

func (q *Queue) Enqueue(val int) {
	q.data = append(q.data, val)
}

func (q *Queue) Dequeue() (int, bool) {
	if len(q.data) == 0 {
		return 0, false
	}
	val := q.data[0]
	q.data = q.data[1:]
	return val, true
}

func (q *Queue) Front() (int, bool) {
	if len(q.data) == 0 {
		return 0, false
	}
	return q.data[0], true
}

func (q *Queue) IsEmpty() bool {
	return len(q.data) == 0
}

func (q *Queue) Size() int {
	return len(q.data)
}

type CircularQueue struct {
	data  []int
	front int
	rear  int
	size  int
	cap   int
}

func NewCircularQueue(k int) *CircularQueue {
	return &CircularQueue{
		data:  make([]int, k),
		front: 0,
		rear:  0,
		size:  0,
		cap:   k,
	}
}

func (cq *CircularQueue) Enqueue(val int) bool {
	if cq.size == cq.cap {
		return false
	}
	cq.data[cq.rear] = val
	cq.rear = (cq.rear + 1) % cq.cap
	cq.size++
	return true
}

func (cq *CircularQueue) Dequeue() (int, bool) {
	if cq.size == 0 {
		return 0, false
	}
	val := cq.data[cq.front]
	cq.front = (cq.front + 1) % cq.cap
	cq.size--
	return val, true
}

func (cq *CircularQueue) Front() (int, bool) {
	if cq.size == 0 {
		return 0, false
	}
	return cq.data[cq.front], true
}

func (cq *CircularQueue) Rear() (int, bool) {
	if cq.size == 0 {
		return 0, false
	}
	return cq.data[(cq.rear-1+cq.cap)%cq.cap], true
}

func (cq *CircularQueue) IsEmpty() bool {
	return cq.size == 0
}

func (cq *CircularQueue) IsFull() bool {
	return cq.size == cq.cap
}

type Deque struct {
	data []int
}

func NewDeque() *Deque {
	return &Deque{data: make([]int, 0)}
}

func (d *Deque) PushFront(val int) {
	d.data = append([]int{val}, d.data...)
}

func (d *Deque) PushBack(val int) {
	d.data = append(d.data, val)
}

func (d *Deque) PopFront() (int, bool) {
	if len(d.data) == 0 {
		return 0, false
	}
	val := d.data[0]
	d.data = d.data[1:]
	return val, true
}

func (d *Deque) PopBack() (int, bool) {
	if len(d.data) == 0 {
		return 0, false
	}
	val := d.data[len(d.data)-1]
	d.data = d.data[:len(d.data)-1]
	return val, true
}

func (d *Deque) Front() (int, bool) {
	if len(d.data) == 0 {
		return 0, false
	}
	return d.data[0], true
}

func (d *Deque) Back() (int, bool) {
	if len(d.data) == 0 {
		return 0, false
	}
	return d.data[len(d.data)-1], true
}

func (d *Deque) IsEmpty() bool {
	return len(d.data) == 0
}

func maxSlidingWindow(nums []int, k int) []int {
	deque := NewDeque()
	result := []int{}
	for i := 0; i < len(nums); i++ {
		for !deque.IsEmpty() {
			back, _ := deque.Back()
			if nums[i] >= nums[back] {
				deque.PopBack()
			} else {
				break
			}
		}
		deque.PushBack(i)
		front, _ := deque.Front()
		if front <= i-k {
			deque.PopFront()
		}
		if i >= k-1 {
			f, _ := deque.Front()
			result = append(result, nums[f])
		}
	}
	return result
}

func main() {
	fmt.Println("=== 队列 (Queue) ===")
	fmt.Println()

	q := NewQueue()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	fmt.Printf("Enqueue 1,2,3 -> Size: %d\n", q.Size())
	if v, ok := q.Dequeue(); ok {
		fmt.Printf("Dequeue: %d\n", v)
	}
	if v, ok := q.Front(); ok {
		fmt.Printf("Front: %d\n", v)
	}

	fmt.Println()
	fmt.Println("--- 循环队列 ---")
	cq := NewCircularQueue(3)
	fmt.Printf("Enqueue 1: %v\n", cq.Enqueue(1))
	fmt.Printf("Enqueue 2: %v\n", cq.Enqueue(2))
	fmt.Printf("Enqueue 3: %v\n", cq.Enqueue(3))
	fmt.Printf("Enqueue 4 (满): %v\n", cq.Enqueue(4))
	fmt.Printf("IsFull: %v\n", cq.IsFull())
	if v, ok := cq.Front(); ok {
		fmt.Printf("Front: %d\n", v)
	}
	if v, ok := cq.Rear(); ok {
		fmt.Printf("Rear: %d\n", v)
	}
	if v, ok := cq.Dequeue(); ok {
		fmt.Printf("Dequeue: %d\n", v)
	}
	fmt.Printf("Enqueue 4: %v\n", cq.Enqueue(4))
	if v, ok := cq.Rear(); ok {
		fmt.Printf("Rear: %d\n", v)
	}

	fmt.Println()
	fmt.Println("--- 双端队列 ---")
	dq := NewDeque()
	dq.PushBack(2)
	dq.PushFront(1)
	dq.PushBack(3)
	fmt.Printf("PushFront(1), PushBack(2), PushBack(3)\n")
	if f, ok := dq.Front(); ok {
		fmt.Printf("Front: %d\n", f)
	}
	if b, ok := dq.Back(); ok {
		fmt.Printf("Back: %d\n", b)
	}

	fmt.Println()
	fmt.Println("--- 滑动窗口最大值 ---")
	nums := []int{1, 3, -1, -3, 5, 3, 6, 7}
	k := 3
	fmt.Printf("  输入: %v, k=%d\n", nums, k)
	fmt.Printf("  结果: %v\n", maxSlidingWindow(nums, k))

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. 队列FIFO特性: BFS、层序遍历、任务调度")
	fmt.Println("2. 循环队列: 解决假溢出，front=(front+1)%cap, rear=(rear+1)%cap")
	fmt.Println("3. 队满条件: (rear+1)%cap == front (牺牲一个存储单元)")
	fmt.Println("4. 双端队列: 同时支持两端入队出队，滑动窗口最大值经典应用")
	fmt.Println("5. 链式队列: 无需预分配空间，无假溢出问题")
	fmt.Println("6. 时间复杂度: Enqueue/Dequeue 均O(1)")
}
