package main

import "fmt"

type PageRef struct {
	Page     int
	Modified bool
}

func fifo(pages []int, frameCount int) (int, float64) {
	frames := make([]int, frameCount)
	for i := range frames {
		frames[i] = -1
	}
	faults := 0
	pointer := 0
	inFrame := make(map[int]bool)

	fmt.Printf("FIFO (帧数=%d):\n", frameCount)
	fmt.Printf("访问  帧状态          缺页\n")

	for _, page := range pages {
		if inFrame[page] {
			fmt.Printf("  %2d   %v   命中\n", page, frameState(frames))
			continue
		}

		faults++
		if frames[pointer] != -1 {
			delete(inFrame, frames[pointer])
		}
		frames[pointer] = page
		inFrame[page] = true
		pointer = (pointer + 1) % frameCount

		fmt.Printf("  %2d   %v   缺页!\n", page, frameState(frames))
	}

	rate := float64(faults) / float64(len(pages)) * 100
	fmt.Printf("缺页次数: %d, 缺页率: %.1f%%\n\n", faults, rate)
	return faults, rate
}

func lru(pages []int, frameCount int) (int, float64) {
	frames := make([]int, frameCount)
	for i := range frames {
		frames[i] = -1
	}
	faults := 0
	count := 0
	lastUsed := make(map[int]int)

	fmt.Printf("LRU (帧数=%d):\n", frameCount)
	fmt.Printf("访问  帧状态          缺页\n")

	for _, page := range pages {
		count++
		lastUsed[page] = count

		inFrame := false
		for _, f := range frames {
			if f == page {
				inFrame = true
				break
			}
		}

		if inFrame {
			fmt.Printf("  %2d   %v   命中\n", page, frameState(frames))
			continue
		}

		faults++
		replaceIdx := -1
		minTime := 1 << 30
		for i, f := range frames {
			if f == -1 {
				replaceIdx = i
				break
			}
			if lastUsed[f] < minTime {
				minTime = lastUsed[f]
				replaceIdx = i
			}
		}
		frames[replaceIdx] = page
		fmt.Printf("  %2d   %v   缺页!\n", page, frameState(frames))
	}

	rate := float64(faults) / float64(len(pages)) * 100
	fmt.Printf("缺页次数: %d, 缺页率: %.1f%%\n\n", faults, rate)
	return faults, rate
}

func opt(pages []int, frameCount int) (int, float64) {
	frames := make([]int, frameCount)
	for i := range frames {
		frames[i] = -1
	}
	faults := 0

	fmt.Printf("OPT (帧数=%d):\n", frameCount)
	fmt.Printf("访问  帧状态          缺页\n")

	for idx, page := range pages {
		inFrame := false
		for _, f := range frames {
			if f == page {
				inFrame = true
				break
			}
		}

		if inFrame {
			fmt.Printf("  %2d   %v   命中\n", page, frameState(frames))
			continue
		}

		faults++
		replaceIdx := -1
		hasEmpty := false
		for i, f := range frames {
			if f == -1 {
				replaceIdx = i
				hasEmpty = true
				break
			}
		}

		if !hasEmpty {
			farthest := -1
			farthestUse := -1
			for i, f := range frames {
				nextUse := 1 << 30
				for j := idx + 1; j < len(pages); j++ {
					if pages[j] == f {
						nextUse = j
						break
					}
				}
				if nextUse > farthestUse {
					farthestUse = nextUse
					farthest = i
				}
			}
			replaceIdx = farthest
		}

		frames[replaceIdx] = page
		fmt.Printf("  %2d   %v   缺页!\n", page, frameState(frames))
	}

	rate := float64(faults) / float64(len(pages)) * 100
	fmt.Printf("缺页次数: %d, 缺页率: %.1f%%\n\n", faults, rate)
	return faults, rate
}

func clock(pages []int, frameCount int) (int, float64) {
	type Frame struct {
		Page    int
		RefBit  bool
	}
	frames := make([]Frame, frameCount)
	for i := range frames {
		frames[i].Page = -1
	}
	faults := 0
	pointer := 0

	fmt.Printf("CLOCK (帧数=%d):\n", frameCount)

	for _, page := range pages {
		hit := false
		for i := range frames {
			if frames[i].Page == page {
				frames[i].RefBit = true
				hit = true
				break
			}
		}

		if hit {
			continue
		}

		faults++
		for {
			if frames[pointer].Page == -1 {
				frames[pointer].Page = page
				frames[pointer].RefBit = true
				pointer = (pointer + 1) % frameCount
				break
			}
			if !frames[pointer].RefBit {
				frames[pointer].Page = page
				frames[pointer].RefBit = true
				pointer = (pointer + 1) % frameCount
				break
			}
			frames[pointer].RefBit = false
			pointer = (pointer + 1) % frameCount
		}
	}

	rate := float64(faults) / float64(len(pages)) * 100
	fmt.Printf("缺页次数: %d, 缺页率: %.1f%%\n\n", faults, rate)
	return faults, rate
}

func frameState(frames []int) []int {
	result := make([]int, len(frames))
	copy(result, frames)
	return result
}

func main() {
	fmt.Println("=== 页面置换算法 ===")
	fmt.Println()

	pages := []int{7, 0, 1, 2, 0, 3, 0, 4, 2, 3, 0, 3, 2, 1, 2, 0, 1, 7, 0, 1}
	frameCount := 3
	fmt.Printf("页面引用串: %v\n", pages)
	fmt.Printf("物理帧数: %d\n\n", frameCount)

	opt(pages, frameCount)
	fifo(pages, frameCount)
	lru(pages, frameCount)
	clock(pages, frameCount)

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. OPT: 未来最长时间不用的页面换出，理论最优但无法实现")
	fmt.Println("2. FIFO: 先进先出，可能Belady异常(帧数增多缺页反而增多)")
	fmt.Println("3. LRU: 最近最少使用，性能接近OPT，实现开销大")
	fmt.Println("4. CLOCK: LRU的近似，用引用位+循环指针，实现简单")
	fmt.Println("5. Belady异常: 仅FIFO存在，LRU/OPT不会出现")
	fmt.Println("6. 抖动(Thrashing): 缺页频繁，CPU利用率低")
	fmt.Println("7. 驻留集: 分配给进程的物理帧数，影响缺页率")
	fmt.Println("8. 工作集: 进程在时间窗口内访问的页面集合")
	fmt.Println("9. 缺页率 = 缺页次数 / 总访问次数")
}
