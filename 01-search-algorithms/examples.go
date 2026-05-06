// BFS и DFS на графе, заданном списком смежности.
// Запуск: go run examples.go
package main

import "fmt"

type Graph map[int][]int

func BFS(g Graph, start int) []int {
	visited := map[int]bool{start: true}
	queue := []int{start}
	order := []int{}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		order = append(order, v)
		for _, u := range g[v] {
			if !visited[u] {
				visited[u] = true
				queue = append(queue, u)
			}
		}
	}
	return order
}

// BFSShortestPath возвращает кратчайший путь в невзвешенном графе.
// nil — если пути нет.
func BFSShortestPath(g Graph, start, goal int) []int {
	if start == goal {
		return []int{start}
	}
	parents := map[int]int{start: -1}
	queue := []int{start}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, u := range g[v] {
			if _, seen := parents[u]; seen {
				continue
			}
			parents[u] = v
			if u == goal {
				path := []int{u}
				for parents[path[len(path)-1]] != -1 {
					path = append(path, parents[path[len(path)-1]])
				}
				// reverse
				for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
					path[i], path[j] = path[j], path[i]
				}
				return path
			}
			queue = append(queue, u)
		}
	}
	return nil
}

func DFSRecursive(g Graph, start int) []int {
	visited := map[int]bool{}
	order := []int{}
	var go_ func(v int)
	go_ = func(v int) {
		visited[v] = true
		order = append(order, v)
		for _, u := range g[v] {
			if !visited[u] {
				go_(u)
			}
		}
	}
	go_(start)
	return order
}

func DFSIterative(g Graph, start int) []int {
	visited := map[int]bool{start: true}
	stack := []int{start}
	order := []int{}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		neigh := g[v]
		for i := len(neigh) - 1; i >= 0; i-- {
			u := neigh[i]
			if !visited[u] {
				visited[u] = true
				stack = append(stack, u)
			}
		}
	}
	return order
}

func main() {
	g := Graph{
		1: {2, 3},
		2: {1, 4},
		3: {1, 4},
		4: {2, 3, 5},
		5: {4},
	}
	fmt.Println("BFS от 1:", BFS(g, 1))
	fmt.Println("DFS (рекурсия) от 1:", DFSRecursive(g, 1))
	fmt.Println("DFS (итеративный) от 1:", DFSIterative(g, 1))
	fmt.Println("Кратчайший путь 1 -> 5:", BFSShortestPath(g, 1, 5))
}
