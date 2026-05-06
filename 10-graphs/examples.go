// Алгоритмы на графах: Дейкстра, топосорт, DSU, Kruskal.
// Запуск: go run examples.go
package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// ── Дейкстра ───────────────────────────────────────────────────────
type Edge struct {
	To, Weight int
}

type pqItem struct {
	v, d int
}
type pq []pqItem

func (p pq) Len() int            { return len(p) }
func (p pq) Less(i, j int) bool  { return p[i].d < p[j].d }
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x any)         { *p = append(*p, x.(pqItem)) }
func (p *pq) Pop() any {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[:n-1]
	return x
}

func Dijkstra(graph map[int][]Edge, source int, n int) []int {
	const INF = 1 << 30
	dist := make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[source] = 0

	q := &pq{}
	heap.Init(q)
	heap.Push(q, pqItem{source, 0})
	for q.Len() > 0 {
		it := heap.Pop(q).(pqItem)
		if it.d > dist[it.v] {
			continue
		}
		for _, e := range graph[it.v] {
			nd := it.d + e.Weight
			if nd < dist[e.To] {
				dist[e.To] = nd
				heap.Push(q, pqItem{e.To, nd})
			}
		}
	}
	return dist
}

// ── Топологическая сортировка (Кан) ────────────────────────────────
func TopoSort(graph map[int][]int, n int) ([]int, bool) {
	inDeg := make([]int, n)
	for _, neigh := range graph {
		for _, v := range neigh {
			inDeg[v]++
		}
	}
	queue := []int{}
	for v := 0; v < n; v++ {
		if inDeg[v] == 0 {
			queue = append(queue, v)
		}
	}
	order := []int{}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)
		for _, v := range graph[u] {
			inDeg[v]--
			if inDeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	return order, len(order) == n
}

// ── DSU ────────────────────────────────────────────────────────────
type DSU struct {
	parent, rank []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), rank: make([]int, n)}
	for i := range d.parent {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) Find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) Union(x, y int) bool {
	rx, ry := d.Find(x), d.Find(y)
	if rx == ry {
		return false
	}
	if d.rank[rx] < d.rank[ry] {
		rx, ry = ry, rx
	}
	d.parent[ry] = rx
	if d.rank[rx] == d.rank[ry] {
		d.rank[rx]++
	}
	return true
}

// ── Kruskal: MST ───────────────────────────────────────────────────
type WEdge struct{ U, V, W int }

func Kruskal(n int, edges []WEdge) (int, []WEdge) {
	sorted := make([]WEdge, len(edges))
	copy(sorted, edges)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].W < sorted[j].W })

	dsu := NewDSU(n)
	total := 0
	tree := []WEdge{}
	for _, e := range sorted {
		if dsu.Union(e.U, e.V) {
			total += e.W
			tree = append(tree, e)
		}
	}
	return total, tree
}

func main() {
	g := map[int][]Edge{
		0: {{1, 4}, {2, 1}},
		1: {{3, 1}},
		2: {{1, 2}, {3, 5}},
		3: {},
	}
	fmt.Println("Dijkstra от 0:", Dijkstra(g, 0, 4))

	dag := map[int][]int{
		0: {1},      // compile -> link
		1: {2},      // link    -> test
		2: {},
		3: {2},      // lint    -> test
	}
	order, ok := TopoSort(dag, 4)
	fmt.Println("Topo:", order, "valid:", ok)

	edges := []WEdge{{0, 1, 4}, {0, 2, 1}, {1, 2, 2}, {1, 3, 5}, {2, 3, 8}}
	w, tree := Kruskal(4, edges)
	fmt.Println("MST вес:", w, "рёбра:", tree)
}
