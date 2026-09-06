package main

import "fmt"

func main() {
	fmt.Println("=== Slice (动态数组，对比Java List) ===")
	slice := []int{1, 2, 3, 4, 5}
	fmt.Printf("Slice: %v\n", slice)
	slice = append(slice, 6, 7)
	fmt.Printf("Append后: %v\n", slice)
	fmt.Printf("切片操作[1:4]: %v\n", slice[1:4])

	fmt.Println("\n=== Map (HashMap) ===")
	userMap := map[string]int{
		"Alice":   25,
		"Bob":     30,
		"Charlie": 28,
	}
	fmt.Printf("Map: %v\n", userMap)
	delete(userMap, "Bob")
	fmt.Printf("Delete Bob后: %v\n", userMap)
	age, ok := userMap["Alice"]
	if ok {
		fmt.Printf("Alice的年龄: %d\n", age)
	}

	fmt.Println("\n=== 循环对比 ===")
	nums := []int{1, 2, 3, 4, 5}

	fmt.Print("Java风格for: ")
	for i := 0; i < len(nums); i++ {
		fmt.Printf("%d ", nums[i])
	}
	fmt.Println()

	fmt.Print("For-range (类似Java forEach): ")
	for i, v := range nums {
		fmt.Printf("[%d:%d] ", i, v)
	}
	fmt.Println()

	fmt.Print("只取value: ")
	for _, v := range nums {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	fmt.Println("\n=== Map遍历 ===")
	for k, v := range userMap {
		fmt.Printf("%s -> %d\n", k, v)
	}
}
