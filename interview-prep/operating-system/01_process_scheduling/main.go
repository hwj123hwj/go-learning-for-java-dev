package main

import "fmt"

type Process struct {
	Name       string
	Arrival    int
	Burst      int
	Priority   int
	Waiting    int
	Turnaround int
	Finish     int
}

func fcfs(processes []Process) []Process {
	result := make([]Process, len(processes))
	copy(result, processes)
	currentTime := 0
	for i := range result {
		if currentTime < result[i].Arrival {
			currentTime = result[i].Arrival
		}
		result[i].Waiting = currentTime - result[i].Arrival
		result[i].Finish = currentTime + result[i].Burst
		result[i].Turnaround = result[i].Finish - result[i].Arrival
		currentTime = result[i].Finish
	}
	return result
}

func sjf(processes []Process) []Process {
	result := make([]Process, len(processes))
	copy(result, processes)
	n := len(result)
	completed := make([]bool, n)
	currentTime := 0
	done := 0

	for done < n {
		shortest := -1
		for i := 0; i < n; i++ {
			if !completed[i] && result[i].Arrival <= currentTime {
				if shortest == -1 || result[i].Burst < result[shortest].Burst {
					shortest = i
				}
			}
		}
		if shortest == -1 {
			currentTime++
			continue
		}
		result[shortest].Waiting = currentTime - result[shortest].Arrival
		result[shortest].Finish = currentTime + result[shortest].Burst
		result[shortest].Turnaround = result[shortest].Finish - result[shortest].Arrival
		currentTime = result[shortest].Finish
		completed[shortest] = true
		done++
	}
	return result
}

func roundRobin(processes []Process, quantum int) []Process {
	result := make([]Process, len(processes))
	copy(result, processes)
	n := len(result)
	remaining := make([]int, n)
	for i := range result {
		remaining[i] = result[i].Burst
	}
	queue := []int{}
	inQueue := make([]bool, n)
	currentTime := 0
	completed := 0

	for i := 0; i < n; i++ {
		if result[i].Arrival == 0 {
			queue = append(queue, i)
			inQueue[i] = true
		}
	}

	for completed < n {
		if len(queue) == 0 {
			currentTime++
			for i := 0; i < n; i++ {
				if !inQueue[i] && remaining[i] > 0 && result[i].Arrival <= currentTime {
					queue = append(queue, i)
					inQueue[i] = true
				}
			}
			continue
		}

		idx := queue[0]
		queue = queue[1:]

		execTime := quantum
		if remaining[idx] < quantum {
			execTime = remaining[idx]
		}
		remaining[idx] -= execTime
		currentTime += execTime

		for i := 0; i < n; i++ {
			if !inQueue[i] && remaining[i] > 0 && result[i].Arrival <= currentTime {
				queue = append(queue, i)
				inQueue[i] = true
			}
		}

		if remaining[idx] == 0 {
			result[idx].Finish = currentTime
			result[idx].Turnaround = result[idx].Finish - result[idx].Arrival
			result[idx].Waiting = result[idx].Turnaround - result[idx].Burst
			completed++
		} else {
			queue = append(queue, idx)
		}
	}
	return result
}

func priorityScheduling(processes []Process) []Process {
	result := make([]Process, len(processes))
	copy(result, processes)
	n := len(result)
	completed := make([]bool, n)
	currentTime := 0
	done := 0

	for done < n {
		highest := -1
		for i := 0; i < n; i++ {
			if !completed[i] && result[i].Arrival <= currentTime {
				if highest == -1 || result[i].Priority < result[highest].Priority {
					highest = i
				}
			}
		}
		if highest == -1 {
			currentTime++
			continue
		}
		result[highest].Waiting = currentTime - result[highest].Arrival
		result[highest].Finish = currentTime + result[highest].Burst
		result[highest].Turnaround = result[highest].Finish - result[highest].Arrival
		currentTime = result[highest].Finish
		completed[highest] = true
		done++
	}
	return result
}

func printSchedule(name string, processes []Process) {
	fmt.Printf("--- %s ---\n", name)
	fmt.Printf("%-8s %-8s %-8s %-8s %-10s %-10s\n", "进程", "到达", "服务", "优先级", "等待", "周转")
	totalWait := 0
	totalTurn := 0
	for _, p := range processes {
		fmt.Printf("%-8s %-8d %-8d %-8d %-10d %-10d\n",
			p.Name, p.Arrival, p.Burst, p.Priority, p.Waiting, p.Turnaround)
		totalWait += p.Waiting
		totalTurn += p.Turnaround
	}
	n := len(processes)
	fmt.Printf("平均等待时间: %.2f, 平均周转时间: %.2f\n",
		float64(totalWait)/float64(n), float64(totalTurn)/float64(n))
	fmt.Println()
}

func main() {
	fmt.Println("=== 进程调度算法 ===")
	fmt.Println()

	processes := []Process{
		{"P1", 0, 6, 3, 0, 0, 0},
		{"P2", 1, 4, 1, 0, 0, 0},
		{"P3", 2, 8, 4, 0, 0, 0},
		{"P4", 3, 2, 2, 0, 0, 0},
	}

	printSchedule("FCFS(先来先服务)", fcfs(processes))
	printSchedule("SJF(短作业优先)", sjf(processes))
	printSchedule("优先级调度(数字越小优先级越高)", priorityScheduling(processes))
	printSchedule("时间片轮转(量子=2)", roundRobin(processes, 2))

	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. FCFS: 简单但对短作业不利，平均等待时间长")
	fmt.Println("2. SJF: 平均等待时间最短，但可能饥饿长作业")
	fmt.Println("3. 优先级调度: 可抢占/非抢占，低优先级可能饥饿")
	fmt.Println("4. 时间片轮转: 公平，时间片太大退化为FCFS，太小切换开销大")
	fmt.Println("5. 多级反馈队列: 综合以上优点，考研最爱考")
	fmt.Println("6. 周转时间 = 完成时间 - 到达时间")
	fmt.Println("7. 等待时间 = 周转时间 - 服务时间")
	fmt.Println("8. 带权周转时间 = 周转时间 / 服务时间")
}
