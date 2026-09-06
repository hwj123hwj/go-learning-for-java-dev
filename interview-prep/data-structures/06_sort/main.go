package main

import "fmt"

func bubbleSort(arr []int) []int {
	n := len(arr)
	a := make([]int, n)
	copy(a, arr)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-1-i; j++ {
			if a[j] > a[j+1] {
				a[j], a[j+1] = a[j+1], a[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
	return a
}

func selectionSort(arr []int) []int {
	n := len(arr)
	a := make([]int, n)
	copy(a, arr)
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if a[j] < a[minIdx] {
				minIdx = j
			}
		}
		a[i], a[minIdx] = a[minIdx], a[i]
	}
	return a
}

func insertionSort(arr []int) []int {
	n := len(arr)
	a := make([]int, n)
	copy(a, arr)
	for i := 1; i < n; i++ {
		key := a[i]
		j := i - 1
		for j >= 0 && a[j] > key {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = key
	}
	return a
}

func shellSort(arr []int) []int {
	n := len(arr)
	a := make([]int, n)
	copy(a, arr)
	for gap := n / 2; gap > 0; gap /= 2 {
		for i := gap; i < n; i++ {
			key := a[i]
			j := i
			for j >= gap && a[j-gap] > key {
				a[j] = a[j-gap]
				j -= gap
			}
			a[j] = key
		}
	}
	return a
}

func quickSort(arr []int) []int {
	a := make([]int, len(arr))
	copy(a, arr)
	var qsort func(lo, hi int)
	qsort = func(lo, hi int) {
		if lo >= hi {
			return
		}
		pivot := a[lo]
		i, j := lo, hi
		for i < j {
			for i < j && a[j] >= pivot {
				j--
			}
			a[i] = a[j]
			for i < j && a[i] <= pivot {
				i++
			}
			a[j] = a[i]
		}
		a[i] = pivot
		qsort(lo, i-1)
		qsort(i+1, hi)
	}
	qsort(0, len(a)-1)
	return a
}

func mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := mergeSort(arr[:mid])
	right := mergeSort(arr[mid:])
	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

func heapSort(arr []int) []int {
	n := len(arr)
	a := make([]int, n)
	copy(a, arr)

	var heapify func(size, i int)
	heapify = func(size, i int) {
		largest := i
		left := 2*i + 1
		right := 2*i + 2
		if left < size && a[left] > a[largest] {
			largest = left
		}
		if right < size && a[right] > a[largest] {
			largest = right
		}
		if largest != i {
			a[i], a[largest] = a[largest], a[i]
			heapify(size, largest)
		}
	}

	for i := n/2 - 1; i >= 0; i-- {
		heapify(n, i)
	}
	for i := n - 1; i > 0; i-- {
		a[0], a[i] = a[i], a[0]
		heapify(i, 0)
	}
	return a
}

func countingSort(arr []int) []int {
	if len(arr) == 0 {
		return arr
	}
	maxVal := arr[0]
	for _, v := range arr {
		if v > maxVal {
			maxVal = v
		}
	}
	count := make([]int, maxVal+1)
	for _, v := range arr {
		count[v]++
	}
	result := make([]int, 0, len(arr))
	for i := 0; i <= maxVal; i++ {
		for count[i] > 0 {
			result = append(result, i)
			count[i]--
		}
	}
	return result
}

func main() {
	fmt.Println("=== 排序算法 ===")
	fmt.Println()

	arr := []int{49, 38, 65, 97, 76, 13, 27, 49}
	fmt.Printf("原始数组: %v\n", arr)
	fmt.Println()

	fmt.Printf("冒泡排序: %v\n", bubbleSort(arr))
	fmt.Printf("选择排序: %v\n", selectionSort(arr))
	fmt.Printf("插入排序: %v\n", insertionSort(arr))
	fmt.Printf("希尔排序: %v\n", shellSort(arr))
	fmt.Printf("快速排序: %v\n", quickSort(arr))
	fmt.Printf("归并排序: %v\n", mergeSort(arr))
	fmt.Printf("堆排序:   %v\n", heapSort(arr))
	fmt.Printf("计数排序: %v\n", countingSort(arr))

	fmt.Println()
	fmt.Println("=== 考研排序算法对比 ===")
	fmt.Println()
	fmt.Println("算法       | 平均时间   | 最坏时间   | 最好时间   | 空间      | 稳定性")
	fmt.Println("-----------|-----------|-----------|-----------|----------|-------")
	fmt.Println("冒泡排序   | O(n^2)    | O(n^2)    | O(n)      | O(1)     | 稳定")
	fmt.Println("选择排序   | O(n^2)    | O(n^2)    | O(n^2)    | O(1)     | 不稳定")
	fmt.Println("插入排序   | O(n^2)    | O(n^2)    | O(n)      | O(1)     | 稳定")
	fmt.Println("希尔排序   | O(n^1.3)  | O(n^2)    | O(n)      | O(1)     | 不稳定")
	fmt.Println("快速排序   | O(nlogn)  | O(n^2)    | O(nlogn)  | O(logn)  | 不稳定")
	fmt.Println("归并排序   | O(nlogn)  | O(nlogn)  | O(nlogn)  | O(n)     | 稳定")
	fmt.Println("堆排序     | O(nlogn)  | O(nlogn)  | O(nlogn)  | O(1)     | 不稳定")
	fmt.Println("计数排序   | O(n+k)    | O(n+k)    | O(n+k)    | O(k)     | 稳定")
	fmt.Println()
	fmt.Println("考研重点:")
	fmt.Println("1. 快速排序: 枢轴选取影响性能，最坏O(n^2)出现在已排序数组")
	fmt.Println("2. 归并排序: 唯一稳定O(nlogn)，但需O(n)额外空间")
	fmt.Println("3. 堆排序: 原地O(nlogn)，但不稳定，建堆O(n)")
	fmt.Println("4. 稳定性: 冒泡/插入/归并/计数 稳定; 选择/希尔/快排/堆排 不稳定")
	fmt.Println("5. 初始状态影响: 冒泡/插入 最好O(n); 选择/快排 最坏O(n^2)")
	fmt.Println("6. 基数排序: 按位排序，O(d(n+r))，d为位数，r为基数")
}
