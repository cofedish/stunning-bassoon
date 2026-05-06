// Жадные алгоритмы на Go.
// Запуск: go run examples.go
package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// ── Activity selection ─────────────────────────────────────────────
type Interval struct{ Start, Finish int }

func ActivitySelection(intervals []Interval) []Interval {
	sorted := make([]Interval, len(intervals))
	copy(sorted, intervals)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Finish < sorted[j].Finish })

	chosen := []Interval{}
	lastEnd := -1 << 30
	for _, iv := range sorted {
		if iv.Start >= lastEnd {
			chosen = append(chosen, iv)
			lastEnd = iv.Finish
		}
	}
	return chosen
}

// ── Размен монет (канонические номиналы) ───────────────────────────
func CoinChangeGreedy(coins []int, amount int) []int {
	desc := make([]int, len(coins))
	copy(desc, coins)
	sort.Sort(sort.Reverse(sort.IntSlice(desc)))

	res := []int{}
	for _, c := range desc {
		for amount >= c {
			amount -= c
			res = append(res, c)
		}
	}
	if amount != 0 {
		return nil
	}
	return res
}

// ── Хаффман: построение кода через min-heap ────────────────────────
type huffNode struct {
	freq        int
	symbol      byte // 0 если внутренний узел
	left, right *huffNode
}

type huffHeap []*huffNode

func (h huffHeap) Len() int            { return len(h) }
func (h huffHeap) Less(i, j int) bool  { return h[i].freq < h[j].freq }
func (h huffHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *huffHeap) Push(x any)         { *h = append(*h, x.(*huffNode)) }
func (h *huffHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func HuffmanCodes(freq map[byte]int) map[byte]string {
	h := &huffHeap{}
	heap.Init(h)
	for sym, f := range freq {
		heap.Push(h, &huffNode{freq: f, symbol: sym})
	}
	for h.Len() > 1 {
		a := heap.Pop(h).(*huffNode)
		b := heap.Pop(h).(*huffNode)
		heap.Push(h, &huffNode{freq: a.freq + b.freq, left: a, right: b})
	}

	codes := map[byte]string{}
	if h.Len() == 0 {
		return codes
	}
	var walk func(n *huffNode, prefix string)
	walk = func(n *huffNode, prefix string) {
		if n.left == nil && n.right == nil {
			if prefix == "" {
				prefix = "0"
			}
			codes[n.symbol] = prefix
			return
		}
		walk(n.left, prefix+"0")
		walk(n.right, prefix+"1")
	}
	walk(heap.Pop(h).(*huffNode), "")
	return codes
}

func main() {
	iv := []Interval{
		{1, 4}, {3, 5}, {0, 6}, {5, 7}, {3, 9},
		{5, 9}, {6, 10}, {8, 11}, {8, 12}, {2, 14}, {12, 16},
	}
	fmt.Println("ActivitySelection:", ActivitySelection(iv))
	fmt.Println("CoinChangeGreedy [1,5,10,25,100], 187:", CoinChangeGreedy([]int{1, 5, 10, 25, 100}, 187))

	codes := HuffmanCodes(map[byte]int{'a': 5, 'b': 9, 'c': 12, 'd': 13, 'e': 16, 'f': 45})
	for sym, code := range codes {
		fmt.Printf("  %c -> %s\n", sym, code)
	}
}
