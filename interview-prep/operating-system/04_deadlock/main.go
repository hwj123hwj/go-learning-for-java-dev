package main

import (
	"fmt"
)

type Resource string

type Allocation struct {
	Process   string
	Allocated map[Resource]int
	Need      map[Resource]int
}

func bankerSafety(processes []string, available map[Resource]int, allocated []map[Resource]int, need []map[Resource]int) ([]string, bool) {
	n := len(processes)
	work := make(map[Resource]int)
	for k, v := range available {
		work[k] = v
	}
	finish := make([]bool, n)
	safeSeq := []string{}

	for {
		found := false
		for i := 0; i < n; i++ {
			if finish[i] {
				continue
			}
			canAllocate := true
			for r, v := range need[i] {
				if v > work[r] {
					canAllocate = false
					break
				}
			}
			if canAllocate {
				for r, v := range allocated[i] {
					work[r] += v
				}
				finish[i] = true
				safeSeq = append(safeSeq, processes[i])
				found = true
			}
		}
		if !found {
			break
		}
	}

	allSafe := true
	for _, f := range finish {
		if !f {
			allSafe = false
			break
		}
	}
	return safeSeq, allSafe
}

func bankerRequest(processes []string, available map[Resource]int, allocated []map[Resource]int, need []map[Resource]int, processIdx int, request map[Resource]int) bool {
	fmt.Printf("进程 %s 请求资源: ", processes[processIdx])
	for r, v := range request {
		fmt.Printf("%s=%d ", r, v)
	}
	fmt.Println()

	for r, v := range request {
		if v > need[processIdx][r] {
			fmt.Println("  错误: 请求超过最大需求量!")
			return false
		}
		if v > available[r] {
			fmt.Println("  等待: 请求超过可用资源!")
			return false
		}
	}

	newAvailable := make(map[Resource]int)
	for r, v := range available {
		newAvailable[r] = v
	}
	newAllocated := make([]map[Resource]int, len(allocated))
	newNeed := make([]map[Resource]int, len(need))
	for i := range processes {
		newAllocated[i] = make(map[Resource]int)
		newNeed[i] = make(map[Resource]int)
		for r, v := range allocated[i] {
			newAllocated[i][r] = v
		}
		for r, v := range need[i] {
			newNeed[i][r] = v
		}
	}

	for r, v := range request {
		newAvailable[r] -= v
		newAllocated[processIdx][r] += v
		newNeed[processIdx][r] -= v
	}

	safeSeq, isSafe := bankerSafety(processes, newAvailable, newAllocated, newNeed)
	if isSafe {
		fmt.Printf("  分配成功! 安全序列: %v\n", safeSeq)
		return true
	}
	fmt.Println("  分配失败: 会导致不安全状态!")
	return false
}

type Edge struct {
	From string
	To   string
}

func detectDeadlock(processes []string, resources []Resource, requestEdges []Edge, assignmentEdges []Edge) bool {
	nodes := []string{}
	for _, p := range processes {
		nodes = append(nodes, p)
	}
	for _, r := range resources {
		nodes = append(nodes, string(r))
	}

	adjList := make(map[string][]string)
	for _, e := range requestEdges {
		adjList[e.From] = append(adjList[e.From], e.To)
	}
	for _, e := range assignmentEdges {
		adjList[e.From] = append(adjList[e.From], e.To)
	}

	visited := make(map[string]bool)
	inStack := make(map[string]bool)

	var dfs func(node string) bool
	dfs = func(node string) bool {
		visited[node] = true
		inStack[node] = true
		for _, neighbor := range adjList[node] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if inStack[neighbor] {
				return true
			}
		}
		inStack[node] = false
		return false
	}

	for _, node := range nodes {
		if !visited[node] {
			if dfs(node) {
				return true
			}
		}
	}
	return false
}

func main() {
	fmt.Println("=== 死锁与银行家算法 ===")
	fmt.Println()

	fmt.Println("--- 银行家算法: 安全性检查 ---")
	processes := []string{"P0", "P1", "P2", "P3", "P4"}
	resources := []Resource{"A", "B", "C"}

	available := map[Resource]int{"A": 3, "B": 3, "C": 2}
	allocated := []map[Resource]int{
		{"A": 0, "B": 1, "C": 0},
		{"A": 2, "B": 0, "C": 0},
		{"A": 3, "B": 0, "C": 2},
		{"A": 2, "B": 1, "C": 1},
		{"A": 0, "B": 0, "C": 2},
	}
	maxNeed := []map[Resource]int{
		{"A": 7, "B": 5, "C": 3},
		{"A": 3, "B": 2, "C": 2},
		{"A": 9, "B": 0, "C": 2},
		{"A": 2, "B": 2, "C": 2},
		{"A": 4, "B": 3, "C": 3},
	}
	need := make([]map[Resource]int, len(processes))
	for i := range need {
		need[i] = make(map[Resource]int)
		for _, r := range resources {
			need[i][r] = maxNeed[i][r] - allocated[i][r]
		}
	}

	fmt.Println("Available:", available)
	fmt.Println("Allocated + Max -> Need:")
	for i, p := range processes {
		fmt.Printf("  %s: Allocated=%v, Max=%v, Need=%v\n",
			p, allocated[i], maxNeed[i], need[i])
	}
	fmt.Println()

	safeSeq, isSafe := bankerSafety(processes, available, allocated, need)
	fmt.Printf("初始状态安全: %v, 安全序列: %v\n\n", isSafe, safeSeq)

	fmt.Println("--- 银行家算法: 资源请求 ---")
	bankerRequest(processes, available, allocated, need, 1, map[Resource]int{"A": 1, "B": 0, "C": 2})
	fmt.Println()

	fmt.Println("--- 资源分配图检测死锁 ---")
	procNodes := []string{"P1", "P2"}
	resNodes := []Resource{"R1", "R2"}
	requestEdges := []Edge{
		{"P1", "R1"},
		{"P2", "R2"},
	}
	assignmentEdges := []Edge{
		{"R1", "P2"},
		{"R2", "P1"},
	}
	fmt.Println("P1请求R1, R1分配给P2, P2请求R2, R2分配给P1")
	hasDeadlock := detectDeadlock(procNodes, resNodes, requestEdges, assignmentEdges)
	fmt.Printf("存在死锁: %v\n\n", hasDeadlock)

	noCycleRequest := []Edge{{"P1", "R1"}}
	noCycleAssign := []Edge{{"R1", "P1"}}
	fmt.Println("P1请求R1, R1分配给P1 (无环)")
	hasDeadlock2 := detectDeadlock([]string{"P1"}, []Resource{"R1"}, noCycleRequest, noCycleAssign)
	fmt.Printf("存在死锁: %v\n\n", hasDeadlock2)

	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. 死锁四个必要条件: 互斥、占有并等待、非抢占、循环等待")
	fmt.Println("2. 死锁预防: 破坏四个条件之一(如: 一次性申请所有资源)")
	fmt.Println("3. 死锁避免: 银行家算法，检查安全性")
	fmt.Println("4. 死锁检测: 资源分配图找环")
	fmt.Println("5. 死锁解除: 剥夺资源、撤销进程")
	fmt.Println("6. 安全状态: 存在安全序列，所有进程可按序完成")
	fmt.Println("7. 银行家算法: Request <= Need, Request <= Available, 试分配后检查安全性")
	fmt.Println("8. 资源分配图: 有环且每种资源只有一个实例 => 死锁")
}
