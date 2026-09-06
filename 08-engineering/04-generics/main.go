package main

import (
	"fmt"
	"strings"
)

// Number 是"类型约束"：类型集合的并集。
// ~int 表示"底层类型是 int 的所有类型"（含 type MyInt int），
// 类似 Java 的 <T extends Number>，但更精确。
type Number interface {
	~int | ~int64 | ~float64
}

// ─── 泛型工具函数：Java 里只能写 Stream API 或重载，Go 一个函数搞定 ──────────

func Map[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, len(xs))
	for i, v := range xs {
		out[i] = f(v)
	}
	return out
}

func Filter[T any](xs []T, keep func(T) bool) []T {
	var out []T
	for _, v := range xs {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

func Sum[T Number](xs []T) T {
	var total T // 零值可用：数字类型的零值就是 0
	for _, v := range xs {
		total += v
	}
	return total
}

// ─── 泛型容器：编译期类型安全，没有 Java 的类型擦除烦恼 ────────────────────────

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

func (s *Stack[T]) Pop() (T, bool) {
	var zero T // 泛型零值：无法写 T{}，但可以声明 var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

func main() {
	fmt.Println("=== 泛型：类型参数 ───────────")

	words := []string{"go", "rust", "java", "python"}

	// Map：[]string → []int，编译期确定类型，IDE 补全完整
	lengths := Map(words, func(s string) int { return len(s) })
	fmt.Println("Map 取长度:  ", lengths)

	upper := Map(words, strings.ToUpper)
	fmt.Println("Map 转大写:  ", upper)

	// Filter
	short := Filter(words, func(s string) bool { return len(s) <= 4 })
	fmt.Println("Filter 短词: ", short)

	// Sum：约束保证只接受数字
	fmt.Println("Sum 长度总和:", Sum(lengths))
	fmt.Println("Sum []float64:", Sum([]float64{1.5, 2.5, 3.0}))

	// Stack[T]：泛型容器
	fmt.Println()
	fmt.Println("=== 泛型容器 Stack[T] ─────────")
	stack := &Stack[string]{}
	stack.Push("first")
	stack.Push("second")
	for {
		v, ok := stack.Pop()
		if !ok {
			break
		}
		fmt.Println("pop:", v)
	}

	fmt.Println()
	compareWithJava()
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java（类型擦除，运行时没有 T）:")
	fmt.Println(`  List<String> list = new ArrayList<>();  // 编译后就是 List
  T first(List<T> list) { ... }             // 不能 new T()、不能 T.class`)
	fmt.Println()
	fmt.Println("Go（编译期生成特化代码，运行时类型真实存在）:")
	fmt.Println(`  func Map[T, U any](xs []T, f func(T) U) []U
  var zero T          // 有零值，绕开"不能 new T"的问题`)
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - 从具体代码开始重复两三次再抽象成泛型，别一开始就泛型")
	fmt.Println("  - 约束写窄不写宽：需要比较就 comparable，需要算术就自定义 Number")
	fmt.Println("  - 泛型没有方法级类型参数（Java 有 <T> 方法），容器+函数是主要用途")
	fmt.Println("  - 标准库范例：slices.Contains / maps.Keys / atomic.Pointer[T]")
}
