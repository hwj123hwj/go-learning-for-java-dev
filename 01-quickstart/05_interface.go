package main

import "fmt"

type Animal interface {
	Speak() string
	Move() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "汪汪!"
}

func (d Dog) Move() string {
	return "跑"
}

type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "喵喵!"
}

func (c Cat) Move() string {
	return "跳"
}

type Bird struct {
	Name string
}

func (b Bird) Speak() string {
	return "吱吱!"
}

func (b Bird) Move() string {
	return "飞"
}

func main() {
	fmt.Println("=== 接口与多态 ===")
	fmt.Println("")

	animals := []Animal{
		Dog{Name: "旺财"},
		Cat{Name: "咪咪"},
		Bird{Name: "小小鸟"},
	}

	for i, animal := range animals {
		fmt.Printf("动物%d: %s\n", i+1, animal.Speak())
	}

	fmt.Println("\n=== Go接口的duck typing特性 ===")
	fmt.Println("Go接口是隐式实现的，不需要显式声明implements")
	fmt.Println("只要你实现了接口的所有方法，就自动实现了该接口")

	fmt.Println("\n=== 空接口 interface{} ===")
	var anything interface{} = 42
	fmt.Printf("空接口可以存储任何值: %v, 类型: %T\n", anything, anything)
	anything = "hello"
	fmt.Printf("空接口可以存储任何值: %v, 类型: %T\n", anything, anything)
}
