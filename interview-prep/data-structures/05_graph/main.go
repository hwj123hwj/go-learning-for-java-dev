package main

import "fmt"

type Graph struct {
	Vertices int
	AdjList  map[int][]int
	Directed bool
}

func NewGraph(vertices int, directed bool) *Graph {
	return &Graph{
		Vertices: vertices,
		AdjList:  make(map[int][]int),
		Directed: directed,
	}
}

func (g *Graph) AddEdge(from, to int) {
	g.AdjList[from] = append(g.AdjList[from], to)
	if !g.Directed {
		g.AdjList[to] = append(g.AdjList[to], from)
	}
}

func (g *Graph) AddWeightedEdge(from, to int, weight int, weights map[[2]int]int) {
	g.AddEdge(from, to)
	weights[[2]int{from, to}] = weight
	if !g.Directed {
		weights[[2]int{to, from}] = weight
	}
}

func (g *Graph) DFS(start int) []int {
	visited := make(map[int]bool)
	result := []int{}
	var dfs func(v int)
	dfs = func(v int) {
		visited[v] = true
		result = append(result, v)
		for _, neighbor := range g.AdjList[v] {
			if !visited[neighbor] {
				dfs(neighbor)
			}
		}
	}
	dfs(start)
	return result
}

func (g *Graph) BFS(start int) []int {
	visited := make(map[int]bool)
	result := []int{}
	queue := []int{start}
	visited[start] = true

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		result = append(result, v)
		for _, neighbor := range g.AdjList[v] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	return result
}

func (g *Graph) Dijkstra(start int, weights map[[2]int]int) map[int]int {
	dist := make(map[int]int)
	visited := make(map[int]bool)
	for i := 0; i < g.Vertices; i++ {
		dist[i] = 1 << 30
	}
	dist[start] = 0

	for i := 0; i < g.Vertices; i++ {
		u := -1
		minDist := 1 << 30
		for v := 0; v < g.Vertices; v++ {
			if !visited[v] && dist[v] < minDist {
				minDist = dist[v]
				u = v
			}
		}
		if u == -1 {
			break
		}
		visited[u] = true
		for _, v := range g.AdjList[u] {
			w, ok := weights[[2]int{u, v}]
			if ok && dist[u]+w < dist[v] {
				dist[v] = dist[u] + w
			}
		}
	}
	return dist
}

func (g *Graph) TopologicalSort() []int {
	inDegree := make(map[int]int)
	for i := 0; i < g.Vertices; i++ {
		inDegree[i] = 0
	}
	for _, neighbors := range g.AdjList {
		for _, v := range neighbors {
			inDegree[v]++
		}
	}

	queue := []int{}
	for i := 0; i < g.Vertices; i++ {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	result := []int{}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		result = append(result, u)
		for _, v := range g.AdjList[u] {
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	return result
}

type Edge struct {
	From   int
	To     int
	Weight int
}

func (g *Graph) Kruskal(edges []Edge) []Edge {
	parent := make([]int, g.Vertices)
	rank := make([]int, g.Vertices)
	for i := range parent {
		parent[i] = i
	}

	var find func(x int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	union := func(x, y int) bool {
		px, py := find(x), find(y)
		if px == py {
			return false
		}
		if rank[px] < rank[py] {
			px, py = py, px
		}
		parent[py] = px
		if rank[px] == rank[py] {
			rank[px]++
		}
		return true
	}

	for i := 0; i < len(edges)-1; i++ {
		for j := i + 1; j < len(edges); j++ {
			if edges[i].Weight > edges[j].Weight {
				edges[i], edges[j] = edges[j], edges[i]
			}
		}
	}

	mst := []Edge{}
	for _, e := range edges {
		if union(e.From, e.To) {
			mst = append(mst, e)
			if len(mst) == g.Vertices-1 {
				break
			}
		}
	}
	return mst
}

func main() {
	fmt.Println("=== 图 (Graph) ===")
	fmt.Println()

	fmt.Println("--- 无向图 ---")
	g := NewGraph(7, false)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 3)
	g.AddEdge(1, 4)
	g.AddEdge(2, 5)
	g.AddEdge(2, 6)
	fmt.Println("    0")
	fmt.Println("   / \\")
	fmt.Println("  1   2")
	fmt.Println(" /\\  /\\")
	fmt.Println("3 4 5  6")
	fmt.Println()

	fmt.Printf("DFS(从0): %v\n", g.DFS(0))
	fmt.Printf("BFS(从0): %v\n", g.BFS(0))

	fmt.Println()
	fmt.Println("--- Dijkstra最短路径 ---")
	dg := NewGraph(5, true)
	weights := make(map[[2]int]int)
	dg.AddWeightedEdge(0, 1, 10, weights)
	dg.AddWeightedEdge(0, 3, 30, weights)
	dg.AddWeightedEdge(0, 4, 100, weights)
	dg.AddWeightedEdge(1, 2, 50, weights)
	dg.AddWeightedEdge(2, 4, 10, weights)
	dg.AddWeightedEdge(3, 2, 20, weights)
	dg.AddWeightedEdge(3, 4, 60, weights)

	dist := dg.Dijkstra(0, weights)
	fmt.Printf("从顶点0到各顶点最短距离:\n")
	for v := 0; v < 5; v++ {
		fmt.Printf("  0 -> %d: %d\n", v, dist[v])
	}

	fmt.Println()
	fmt.Println("--- 拓扑排序 ---")
	tg := NewGraph(6, true)
	tg.AddEdge(5, 2)
	tg.AddEdge(5, 0)
	tg.AddEdge(4, 0)
	tg.AddEdge(4, 1)
	tg.AddEdge(2, 3)
	tg.AddEdge(3, 1)
	fmt.Printf("拓扑排序结果: %v\n", tg.TopologicalSort())

	fmt.Println()
	fmt.Println("--- Kruskal最小生成树 ---")
	mstGraph := NewGraph(4, false)
	mstEdges := []Edge{
		{0, 1, 1},
		{0, 2, 2},
		{0, 3, 3},
		{1, 2, 4},
		{2, 3, 5},
	}
	mst := mstGraph.Kruskal(mstEdges)
	totalWeight := 0
	fmt.Println("最小生成树边:")
	for _, e := range mst {
		fmt.Printf("  %d - %d (权重%d)\n", e.From, e.To, e.Weight)
		totalWeight += e.Weight
	}
	fmt.Printf("总权重: %d\n", totalWeight)

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. 存储方式: 邻接矩阵(适合稠密图) vs 邻接表(适合稀疏图)")
	fmt.Println("2. DFS: 用栈/递归，O(V+E)，用于连通性判断、环检测")
	fmt.Println("3. BFS: 用队列，O(V+E)，用于最短路径(无权图)、层序遍历")
	fmt.Println("4. Dijkstra: 单源最短路径，不适用负权边，O(V^2)或O(ElogV)")
	fmt.Println("5. 拓扑排序: 仅DAG，Kahn算法(BFS)或DFS逆后序")
	fmt.Println("6. Kruskal: 贪心+并查集，O(ElogE)")
	fmt.Println("7. Prim: 从顶点扩展，O(V^2)或O(ElogV)")
	fmt.Println("8. 关键路径: AOE网，最早/最晚发生时间")
}
