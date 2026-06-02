package main

import "fmt"

type Node struct {
	Val  int
	Next *Node
}

type LinkedList struct {
	Head *Node
	Size int
}

func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

func (l *LinkedList) InsertHead(val int) {
	node := &Node{Val: val, Next: l.Head}
	l.Head = node
	l.Size++
}

func (l *LinkedList) InsertTail(val int) {
	node := &Node{Val: val}
	if l.Head == nil {
		l.Head = node
	} else {
		cur := l.Head
		for cur.Next != nil {
			cur = cur.Next
		}
		cur.Next = node
	}
	l.Size++
}

func (l *LinkedList) InsertAt(index, val int) bool {
	if index < 0 || index > l.Size {
		return false
	}
	if index == 0 {
		l.InsertHead(val)
		return true
	}
	cur := l.Head
	for i := 0; i < index-1; i++ {
		cur = cur.Next
	}
	node := &Node{Val: val, Next: cur.Next}
	cur.Next = node
	l.Size++
	return true
}

func (l *LinkedList) DeleteHead() bool {
	if l.Head == nil {
		return false
	}
	l.Head = l.Head.Next
	l.Size--
	return true
}

func (l *LinkedList) DeleteTail() bool {
	if l.Head == nil {
		return false
	}
	if l.Head.Next == nil {
		l.Head = nil
		l.Size--
		return true
	}
	cur := l.Head
	for cur.Next.Next != nil {
		cur = cur.Next
	}
	cur.Next = nil
	l.Size--
	return true
}

func (l *LinkedList) DeleteAt(index int) bool {
	if index < 0 || index >= l.Size {
		return false
	}
	if index == 0 {
		return l.DeleteHead()
	}
	cur := l.Head
	for i := 0; i < index-1; i++ {
		cur = cur.Next
	}
	cur.Next = cur.Next.Next
	l.Size--
	return true
}

func (l *LinkedList) Search(val int) int {
	cur := l.Head
	index := 0
	for cur != nil {
		if cur.Val == val {
			return index
		}
		cur = cur.Next
		index++
	}
	return -1
}

func (l *LinkedList) Reverse() {
	var prev *Node
	cur := l.Head
	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}
	l.Head = prev
}

func (l *LinkedList) Print() {
	cur := l.Head
	for cur != nil {
		fmt.Printf("%d -> ", cur.Val)
		cur = cur.Next
	}
	fmt.Println("nil")
}

func (l *LinkedList) ToSlice() []int {
	result := make([]int, 0, l.Size)
	cur := l.Head
	for cur != nil {
		result = append(result, cur.Val)
		cur = cur.Next
	}
	return result
}

type DoubleNode struct {
	Val  int
	Prev *DoubleNode
	Next *DoubleNode
}

type DoublyLinkedList struct {
	Head *DoubleNode
	Tail *DoubleNode
	Size int
}

func NewDoublyLinkedList() *DoublyLinkedList {
	return &DoublyLinkedList{}
}

func (dl *DoublyLinkedList) InsertHead(val int) {
	node := &DoubleNode{Val: val}
	if dl.Head == nil {
		dl.Head = node
		dl.Tail = node
	} else {
		node.Next = dl.Head
		dl.Head.Prev = node
		dl.Head = node
	}
	dl.Size++
}

func (dl *DoublyLinkedList) InsertTail(val int) {
	node := &DoubleNode{Val: val}
	if dl.Tail == nil {
		dl.Head = node
		dl.Tail = node
	} else {
		node.Prev = dl.Tail
		dl.Tail.Next = node
		dl.Tail = node
	}
	dl.Size++
}

func (dl *DoublyLinkedList) DeleteHead() bool {
	if dl.Head == nil {
		return false
	}
	if dl.Head == dl.Tail {
		dl.Head = nil
		dl.Tail = nil
	} else {
		dl.Head = dl.Head.Next
		dl.Head.Prev = nil
	}
	dl.Size--
	return true
}

func (dl *DoublyLinkedList) DeleteTail() bool {
	if dl.Tail == nil {
		return false
	}
	if dl.Head == dl.Tail {
		dl.Head = nil
		dl.Tail = nil
	} else {
		dl.Tail = dl.Tail.Prev
		dl.Tail.Next = nil
	}
	dl.Size--
	return true
}

func (dl *DoublyLinkedList) PrintForward() {
	cur := dl.Head
	for cur != nil {
		fmt.Printf("%d <-> ", cur.Val)
		cur = cur.Next
	}
	fmt.Println("nil")
}

func (dl *DoublyLinkedList) PrintBackward() {
	cur := dl.Tail
	for cur != nil {
		fmt.Printf("%d <-> ", cur.Val)
		cur = cur.Prev
	}
	fmt.Println("nil")
}

func main() {
	fmt.Println("=== 单链表 ===")
	fmt.Println()

	ll := NewLinkedList()
	ll.InsertTail(1)
	ll.InsertTail(2)
	ll.InsertTail(3)
	ll.InsertTail(4)
	ll.InsertTail(5)
	fmt.Print("初始链表: ")
	ll.Print()

	ll.InsertHead(0)
	fmt.Print("头插0: ")
	ll.Print()

	ll.InsertAt(3, 99)
	fmt.Print("位置3插入99: ")
	ll.Print()

	idx := ll.Search(99)
	fmt.Printf("查找99: 位置=%d\n", idx)

	ll.Reverse()
	fmt.Print("反转链表: ")
	ll.Print()

	ll.DeleteHead()
	fmt.Print("删除头: ")
	ll.Print()

	ll.DeleteTail()
	fmt.Print("删除尾: ")
	ll.Print()

	ll.DeleteAt(1)
	fmt.Print("删除位置1: ")
	ll.Print()

	fmt.Println()
	fmt.Println("=== 双链表 ===")
	fmt.Println()

	dl := NewDoublyLinkedList()
	dl.InsertTail(1)
	dl.InsertTail(2)
	dl.InsertTail(3)
	dl.InsertHead(0)
	fmt.Print("正向遍历: ")
	dl.PrintForward()
	fmt.Print("反向遍历: ")
	dl.PrintBackward()

	dl.DeleteHead()
	fmt.Print("删除头后正向: ")
	dl.PrintForward()

	dl.DeleteTail()
	fmt.Print("删除尾后反向: ")
	dl.PrintBackward()

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. 单链表: 头插法/尾插法建表，逆置，查找/删除第i个节点")
	fmt.Println("2. 双链表: 支持双向遍历，插入/删除需修改两个指针")
	fmt.Println("3. 时间复杂度:")
	fmt.Println("   - 按值查找: O(n)")
	fmt.Println("   - 按序号查找: O(n)")
	fmt.Println("   - 头插/尾插(有尾指针): O(1)")
	fmt.Println("   - 删除/插入(已知位置): O(1), 查找位置O(n)")
	fmt.Println("4. 空间复杂度: O(n)")
}
