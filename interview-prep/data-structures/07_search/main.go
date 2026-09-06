package main

import "fmt"

func binarySearch(arr []int, target int) int {
	lo, hi := 0, len(arr)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return -1
}

func binarySearchFirst(arr []int, target int) int {
	lo, hi := 0, len(arr)-1
	result := -1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if arr[mid] == target {
			result = mid
			hi = mid - 1
		} else if arr[mid] < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return result
}

func binarySearchLast(arr []int, target int) int {
	lo, hi := 0, len(arr)-1
	result := -1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if arr[mid] == target {
			result = mid
			lo = mid + 1
		} else if arr[mid] < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return result
}

func binarySearchInsertPos(arr []int, target int) int {
	lo, hi := 0, len(arr)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if arr[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

func hashSearch(table map[int]int, key int) (int, bool) {
	val, ok := table[key]
	return val, ok
}

func hashLinearProbe(keys []int, tableSize int) map[int]int {
	table := make(map[int]int)
	for _, k := range keys {
		idx := k % tableSize
		for {
			if _, exists := table[idx]; !exists {
				table[idx] = k
				break
			}
			idx = (idx + 1) % tableSize
		}
	}
	return table
}

type BSTNode struct {
	Key  int
	Left *BSTNode
	Right *BSTNode
}

func bstSearch(root *BSTNode, key int) *BSTNode {
	if root == nil || root.Key == key {
		return root
	}
	if key < root.Key {
		return bstSearch(root.Left, key)
	}
	return bstSearch(root.Right, key)
}

func bstInsert(root *BSTNode, key int) *BSTNode {
	if root == nil {
		return &BSTNode{Key: key}
	}
	if key < root.Key {
		root.Left = bstInsert(root.Left, key)
	} else if key > root.Key {
		root.Right = bstInsert(root.Right, key)
	}
	return root
}

func main() {
	fmt.Println("=== 查找算法 ===")
	fmt.Println()

	fmt.Println("--- 二分查找 ---")
	arr := []int{1, 3, 5, 5, 5, 7, 9, 11, 13, 15}
	fmt.Printf("有序数组: %v\n", arr)
	fmt.Printf("查找7: 位置=%d\n", binarySearch(arr, 7))
	fmt.Printf("查找5(第一个): 位置=%d\n", binarySearchFirst(arr, 5))
	fmt.Printf("查找5(最后一个): 位置=%d\n", binarySearchLast(arr, 5))
	fmt.Printf("查找6(插入位置): 位置=%d\n", binarySearchInsertPos(arr, 6))
	fmt.Printf("查找0(插入位置): 位置=%d\n", binarySearchInsertPos(arr, 0))

	fmt.Println()
	fmt.Println("--- 哈希查找 ---")
	hashTable := map[int]int{
		1: 100,
		2: 200,
		3: 300,
		5: 500,
	}
	fmt.Printf("哈希表: %v\n", hashTable)
	if val, ok := hashSearch(hashTable, 3); ok {
		fmt.Printf("查找key=3: value=%d\n", val)
	}
	if _, ok := hashSearch(hashTable, 4); !ok {
		fmt.Println("查找key=4: 不存在")
	}

	fmt.Println()
	fmt.Println("--- 哈希冲突: 线性探测 ---")
	keys := []int{10, 22, 31, 4, 15, 28, 17, 88, 59}
	tableSize := 11
	fmt.Printf("关键字: %v, 表长: %d\n", keys, tableSize)
	table := hashLinearProbe(keys, tableSize)
	fmt.Println("散列表(位置->关键字):")
	for i := 0; i < tableSize; i++ {
		if v, ok := table[i]; ok {
			fmt.Printf("  [%d] = %d\n", i, v)
		} else {
			fmt.Printf("  [%d] = 空\n", i)
		}
	}

	fmt.Println()
	fmt.Println("--- BST查找 ---")
	var root *BSTNode
	bstVals := []int{50, 30, 70, 20, 40, 60, 80}
	for _, v := range bstVals {
		root = bstInsert(root, v)
	}
	fmt.Printf("BST插入 %v\n", bstVals)
	if node := bstSearch(root, 40); node != nil {
		fmt.Printf("查找40: 找到\n")
	}
	if node := bstSearch(root, 45); node != nil {
		fmt.Printf("查找45: 找到\n")
	} else {
		fmt.Printf("查找45: 未找到\n")
	}

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. 二分查找: 要求顺序存储+有序，O(logn)，仅适用于顺序表")
	fmt.Println("2. 变形二分: 查第一个/最后一个等于target，查找插入位置")
	fmt.Println("3. ASL(平均查找长度): 衡量查找效率的核心指标")
	fmt.Println("   - 顺序查找ASL = (n+1)/2")
	fmt.Println("   - 二分查找ASL ≈ log2(n+1) - 1")
	fmt.Println("4. 哈希表:")
	fmt.Println("   - 除留余数法: H(key) = key % p (p取不大于表长的最大素数)")
	fmt.Println("   - 冲突处理: 开放定址法(线性/二次探测)、拉链法")
	fmt.Println("   - 装填因子α = n/m，影响查找效率")
	fmt.Println("5. BST查找: O(logn)~O(n)，取决于树形态")
	fmt.Println("6. 平衡二叉树(AVL): 保证O(logn)，LL/RR/LR/RL四种旋转")
	fmt.Println("7. B树/B+树: 多路平衡查找树，磁盘友好，数据库索引核心")
	fmt.Println("8. 红黑树: 近似平衡，插入/删除比AVL高效")
}
