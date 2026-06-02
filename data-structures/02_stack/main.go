package main

import "fmt"

type Stack struct {
	data []int
}

func NewStack() *Stack {
	return &Stack{data: make([]int, 0)}
}

func (s *Stack) Push(val int) {
	s.data = append(s.data, val)
}

func (s *Stack) Pop() (int, bool) {
	if len(s.data) == 0 {
		return 0, false
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val, true
}

func (s *Stack) Peek() (int, bool) {
	if len(s.data) == 0 {
		return 0, false
	}
	return s.data[len(s.data)-1], true
}

func (s *Stack) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *Stack) Size() int {
	return len(s.data)
}

func isValidParentheses(s string) bool {
	stack := NewStack()
	pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}
	for _, ch := range s {
		switch ch {
		case '(', '[', '{':
			stack.Push(int(ch))
		case ')', ']', '}':
			if top, ok := stack.Pop(); !ok || top != int(pairs[ch]) {
				return false
			}
		}
	}
	return stack.IsEmpty()
}

func nextGreaterElement(nums []int) []int {
	n := len(nums)
	result := make([]int, n)
	for i := range result {
		result[i] = -1
	}
	stack := NewStack()
	for i := 0; i < n; i++ {
		for !stack.IsEmpty() {
			top, _ := stack.Peek()
			if nums[i] > nums[top] {
				result[top] = nums[i]
				stack.Pop()
			} else {
				break
			}
		}
		stack.Push(i)
	}
	return result
}

func dailyTemperatures(temps []int) []int {
	n := len(temps)
	result := make([]int, n)
	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && temps[i] > temps[stack[len(stack)-1]] {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result[top] = i - top
		}
		stack = append(stack, i)
	}
	return result
}

type MinStack struct {
	data    []int
	minData []int
}

func NewMinStack() *MinStack {
	return &MinStack{}
}

func (ms *MinStack) Push(val int) {
	ms.data = append(ms.data, val)
	if len(ms.minData) == 0 || val <= ms.minData[len(ms.minData)-1] {
		ms.minData = append(ms.minData, val)
	}
}

func (ms *MinStack) Pop() (int, bool) {
	if len(ms.data) == 0 {
		return 0, false
	}
	val := ms.data[len(ms.data)-1]
	ms.data = ms.data[:len(ms.data)-1]
	if val == ms.minData[len(ms.minData)-1] {
		ms.minData = ms.minData[:len(ms.minData)-1]
	}
	return val, true
}

func (ms *MinStack) GetMin() (int, bool) {
	if len(ms.minData) == 0 {
		return 0, false
	}
	return ms.minData[len(ms.minData)-1], true
}

func evaluatePostfix(expr []string) int {
	stack := NewStack()
	for _, token := range expr {
		switch token {
		case "+", "-", "*", "/":
			b, _ := stack.Pop()
			a, _ := stack.Pop()
			var result int
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				result = a / b
			}
			stack.Push(result)
		default:
			var num int
			fmt.Sscanf(token, "%d", &num)
			stack.Push(num)
		}
	}
	val, _ := stack.Pop()
	return val
}

func main() {
	fmt.Println("=== 栈 (Stack) ===")
	fmt.Println()

	s := NewStack()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	fmt.Printf("Push 1,2,3 -> Size: %d\n", s.Size())
	if v, ok := s.Peek(); ok {
		fmt.Printf("Peek: %d\n", v)
	}
	if v, ok := s.Pop(); ok {
		fmt.Printf("Pop: %d\n", v)
	}
	if v, ok := s.Pop(); ok {
		fmt.Printf("Pop: %d\n", v)
	}
	fmt.Printf("Size after 2 pops: %d\n", s.Size())

	fmt.Println()
	fmt.Println("--- 括号匹配 ---")
	tests := []string{"()", "()[]{}", "(]", "([)]", "{[]}", "((()))"}
	for _, t := range tests {
		fmt.Printf("  \"%s\" -> %v\n", t, isValidParentheses(t))
	}

	fmt.Println()
	fmt.Println("--- 下一个更大元素 ---")
	nums := []int{4, 5, 2, 25, 7, 8}
	fmt.Printf("  输入: %v\n", nums)
	fmt.Printf("  结果: %v\n", nextGreaterElement(nums))

	fmt.Println()
	fmt.Println("--- 每日温度 ---")
	temps := []int{73, 74, 75, 71, 69, 72, 76, 73}
	fmt.Printf("  温度: %v\n", temps)
	fmt.Printf("  等待天数: %v\n", dailyTemperatures(temps))

	fmt.Println()
	fmt.Println("--- 最小栈 ---")
	ms := NewMinStack()
	ms.Push(5)
	ms.Push(3)
	ms.Push(7)
	if v, ok := ms.GetMin(); ok {
		fmt.Printf("  Push 5,3,7 -> Min: %d\n", v)
	}
	ms.Pop()
	if v, ok := ms.GetMin(); ok {
		fmt.Printf("  Pop 7 -> Min: %d\n", v)
	}
	ms.Pop()
	if v, ok := ms.GetMin(); ok {
		fmt.Printf("  Pop 3 -> Min: %d\n", v)
	}

	fmt.Println()
	fmt.Println("--- 后缀表达式求值 ---")
	expr := []string{"3", "4", "+", "5", "*"}
	fmt.Printf("  3 4 + 5 * = %d\n", evaluatePostfix(expr))

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. 栈的LIFO特性: 递归、函数调用、括号匹配、表达式求值")
	fmt.Println("2. 单调栈: 下一个更大/更小元素，O(n)时间")
	fmt.Println("3. 最小栈: 辅助栈同步维护最小值，O(1)取最小")
	fmt.Println("4. 顺序栈 vs 链式栈: 顺序栈用数组，链式栈用链表")
	fmt.Println("5. 时间复杂度: Push/Pop/Peek 均O(1)")
}
