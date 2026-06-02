package main

import "fmt"

type Partition struct {
	Start int
	Size  int
	Free  bool
}

func firstFit(partitions []Partition, request int) (int, bool) {
	for i := range partitions {
		if partitions[i].Free && partitions[i].Size >= request {
			return allocate(partitions, i, request)
		}
	}
	return -1, false
}

func bestFit(partitions []Partition, request int) (int, bool) {
	bestIdx := -1
	minDiff := 1 << 30
	for i := range partitions {
		if partitions[i].Free && partitions[i].Size >= request {
			diff := partitions[i].Size - request
			if diff < minDiff {
				minDiff = diff
				bestIdx = i
			}
		}
	}
	if bestIdx == -1 {
		return -1, false
	}
	return allocate(partitions, bestIdx, request)
}

func worstFit(partitions []Partition, request int) (int, bool) {
	worstIdx := -1
	maxSize := 0
	for i := range partitions {
		if partitions[i].Free && partitions[i].Size >= request {
			if partitions[i].Size > maxSize {
				maxSize = partitions[i].Size
				worstIdx = i
			}
		}
	}
	if worstIdx == -1 {
		return -1, false
	}
	return allocate(partitions, worstIdx, request)
}

func allocate(partitions []Partition, idx, request int) (int, bool) {
	if partitions[idx].Size == request {
		partitions[idx].Free = false
		return partitions[idx].Start, true
	}
	newPartition := Partition{
		Start: partitions[idx].Start + request,
		Size:  partitions[idx].Size - request,
		Free:  true,
	}
	partitions[idx].Size = request
	partitions[idx].Free = false
	result := []Partition{}
	result = append(result, partitions[:idx+1]...)
	result = append(result, newPartition)
	result = append(result, partitions[idx+1:]...)
	copy(partitions, result)
	return partitions[idx].Start, true
}

func printPartitions(partitions []Partition) {
	fmt.Println("内存分区表:")
	fmt.Printf("%-8s %-8s %-8s\n", "起始", "大小", "状态")
	for _, p := range partitions {
		status := "已分配"
		if p.Free {
			status = "空闲"
		}
		fmt.Printf("%-8d %-8d %-8s\n", p.Start, p.Size, status)
	}
	fmt.Println()
}

func main() {
	fmt.Println("=== 内存管理: 连续分配 ===")
	fmt.Println()

	fmt.Println("--- 首次适应(First Fit) ---")
	ffPartitions := []Partition{
		{0, 100, true},
		{100, 50, true},
		{150, 200, true},
		{350, 80, true},
	}
	fmt.Println("初始状态:")
	printPartitions(ffPartitions)

	if addr, ok := firstFit(ffPartitions, 60); ok {
		fmt.Printf("分配60KB -> 起始地址%d\n", addr)
	}
	printPartitions(ffPartitions)

	if addr, ok := firstFit(ffPartitions, 30); ok {
		fmt.Printf("分配30KB -> 起始地址%d\n", addr)
	}
	printPartitions(ffPartitions)

	fmt.Println("--- 最佳适应(Best Fit) ---")
	bfPartitions := []Partition{
		{0, 100, true},
		{100, 50, true},
		{150, 200, true},
		{350, 80, true},
	}
	if addr, ok := bestFit(bfPartitions, 40); ok {
		fmt.Printf("分配40KB -> 起始地址%d (选择最小能满足的分区)\n", addr)
	}
	printPartitions(bfPartitions)

	fmt.Println("--- 最坏适应(Worst Fit) ---")
	wfPartitions := []Partition{
		{0, 100, true},
		{100, 50, true},
		{150, 200, true},
		{350, 80, true},
	}
	if addr, ok := worstFit(wfPartitions, 40); ok {
		fmt.Printf("分配40KB -> 起始地址%d (选择最大空闲分区)\n", addr)
	}
	printPartitions(wfPartitions)

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. 首次适应: 按地址递增查找，低地址碎片多，查找快")
	fmt.Println("2. 最佳适应: 选最小能满足的，产生很多小碎片(外部碎片)")
	fmt.Println("3. 最坏适应: 选最大空闲区，避免小碎片，但大作业可能无法分配")
	fmt.Println("4. 内部碎片: 分配给进程但未用完的空间(分页)")
	fmt.Println("5. 外部碎片: 总空闲足够但无法分配(分区)")
	fmt.Println("6. 紧凑(Compaction): 移动进程消除外部碎片，开销大")
	fmt.Println("7. 分页: 页大小固定，有内部碎片，无外部碎片")
	fmt.Println("8. 分段: 按逻辑段划分，有外部碎片，无内部碎片")
	fmt.Println("9. 虚拟地址: 逻辑地址 -> 物理地址，通过页表映射")
}
