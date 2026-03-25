package main

import "fmt"

func main() {
	fmt.Println("=== 变量声明对比 Java vs Go ===")
	fmt.Println("")

	name := "Alice"
	age := 25
	height := 1.75
	isActive := true
	fmt.Printf("简短声明: name=%s, age=%d, height=%.2f, isActive=%v\n", name, age, height, isActive)

	var javaStyle string = "类似Java"
	var count int = 100
	fmt.Printf("显式声明: %s, count=%d\n", javaStyle, count)

	var (
		a = 1
		b = 2
		c = 3
	)
	fmt.Printf("多变量声明: a=%d, b=%d, c=%d\n", a, b, c)

	const Pi = 3.14159
	const Status = "active"
	fmt.Printf("常量: Pi=%.5f, Status=%s\n", Pi, Status)

	fmt.Println("\n=== 基本数据类型 ===")
	fmt.Printf("int: %d, float64: %.2f, string: %s, bool: %v\n", 42, 3.14, "hello", true)

	s1 := "Hello"
	s2 := "Go"
	result := s1 + " " + s2
	fmt.Printf("字符串拼接: %s\n", result)
	fmt.Printf("字符串长度: %d\n", len(result))
	fmt.Printf("字符串切片[0:5]: %s\n", result[0:5])
}
