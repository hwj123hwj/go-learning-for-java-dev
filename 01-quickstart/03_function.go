package main

import (
	"errors"
	"fmt"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("除数不能为0")
	}
	return a / b, nil
}

func add(nums ...int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func main() {
	fmt.Println("=== 函数与错误处理 ===")
	fmt.Println("")

	result, err := divide(10, 2)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("10 / 2 = %d\n", result)
	}

	result, err = divide(10, 0)
	if err != nil {
		fmt.Printf("捕获错误: %v\n", err)
	}

	fmt.Printf("\n可变参数函数 add(1,2,3,4,5) = %d\n", add(1, 2, 3, 4, 5))
}
