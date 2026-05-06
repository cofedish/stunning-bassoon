// Классические задачи ДП на Go.
// Запуск: go run examples.go
package main

import (
	"fmt"
	"sort"
)

// ── Фибоначчи итеративный ──────────────────────────────────────────
func Fib(n int) int {
	if n < 2 {
		return n
	}
	a, b := 0, 1
	for i := 1; i < n; i++ {
		a, b = b, a+b
	}
	return b
}

// ── Размен монет: -1 если невозможно ───────────────────────────────
func CoinChange(coins []int, amount int) int {
	const INF = 1 << 30
	dp := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
		dp[i] = INF
	}
	for s := 1; s <= amount; s++ {
		for _, c := range coins {
			if c <= s && dp[s-c]+1 < dp[s] {
				dp[s] = dp[s-c] + 1
			}
		}
	}
	if dp[amount] == INF {
		return -1
	}
	return dp[amount]
}

// ── 0/1 рюкзак ──────────────────────────────────────────────────────
func Knapsack01(weights, values []int, capacity int) int {
	dp := make([]int, capacity+1)
	for i := 0; i < len(weights); i++ {
		for w := capacity; w >= weights[i]; w-- {
			if dp[w-weights[i]]+values[i] > dp[w] {
				dp[w] = dp[w-weights[i]] + values[i]
			}
		}
	}
	return dp[capacity]
}

// ── LCS ────────────────────────────────────────────────────────────
func LCS(a, b string) int {
	n, m := len(a), len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] > dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}
	return dp[n][m]
}

// ── Расстояние Левенштейна ─────────────────────────────────────────
func EditDistance(a, b string) int {
	n, m := len(a), len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
		dp[i][0] = i
	}
	for j := 0; j <= m; j++ {
		dp[0][j] = j
	}
	min3 := func(x, y, z int) int {
		m := x
		if y < m {
			m = y
		}
		if z < m {
			m = z
		}
		return m
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min3(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}
	return dp[n][m]
}

// ── LIS за O(n log n) ──────────────────────────────────────────────
func LIS(a []int) int {
	tails := []int{}
	for _, x := range a {
		i := sort.SearchInts(tails, x)
		if i == len(tails) {
			tails = append(tails, x)
		} else {
			tails[i] = x
		}
	}
	return len(tails)
}

func main() {
	fmt.Println("Fib(20):", Fib(20))
	fmt.Println("CoinChange [1,3,4], 6:", CoinChange([]int{1, 3, 4}, 6))
	fmt.Println("Knapsack:", Knapsack01([]int{2, 3, 4, 5}, []int{3, 4, 5, 6}, 5))
	fmt.Println("LCS ABCBDAB / BDCABA:", LCS("ABCBDAB", "BDCABA"))
	fmt.Println("EditDistance kitten/sitting:", EditDistance("kitten", "sitting"))
	fmt.Println("LIS:", LIS([]int{10, 9, 2, 5, 3, 7, 101, 18}))
}
