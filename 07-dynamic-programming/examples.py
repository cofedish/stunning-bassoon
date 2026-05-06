"""
Классические задачи ДП.
Запуск: python examples.py
"""

from functools import lru_cache


# ── Фибоначчи: 3 версии ──────────────────────────────────────────────
def fib_naive(n: int) -> int:
    """Экспонента, для иллюстрации."""
    if n < 2:
        return n
    return fib_naive(n - 1) + fib_naive(n - 2)


@lru_cache(maxsize=None)
def fib_memo(n: int) -> int:
    if n < 2:
        return n
    return fib_memo(n - 1) + fib_memo(n - 2)


def fib_iter(n: int) -> int:
    if n < 2:
        return n
    a, b = 0, 1
    for _ in range(n - 1):
        a, b = b, a + b
    return b


# ── Размен монет: min(c) на сумму S ──────────────────────────────────
def coin_change(coins: list[int], amount: int) -> int:
    INF = float("inf")
    dp = [0] + [INF] * amount
    for s in range(1, amount + 1):
        for c in coins:
            if c <= s and dp[s - c] + 1 < dp[s]:
                dp[s] = dp[s - c] + 1
    return -1 if dp[amount] == INF else dp[amount]


# ── 0/1 рюкзак ───────────────────────────────────────────────────────
def knapsack_01(weights: list[int], values: list[int], capacity: int) -> int:
    """O(N * W), память сжата до O(W)."""
    n = len(weights)
    dp = [0] * (capacity + 1)
    for i in range(n):
        # справа налево, чтобы не переиспользовать предмет
        for w in range(capacity, weights[i] - 1, -1):
            dp[w] = max(dp[w], dp[w - weights[i]] + values[i])
    return dp[capacity]


# ── LCS: длина общей подпоследовательности ───────────────────────────
def lcs(a: str, b: str) -> int:
    n, m = len(a), len(b)
    dp = [[0] * (m + 1) for _ in range(n + 1)]
    for i in range(1, n + 1):
        for j in range(1, m + 1):
            if a[i - 1] == b[j - 1]:
                dp[i][j] = dp[i - 1][j - 1] + 1
            else:
                dp[i][j] = max(dp[i - 1][j], dp[i][j - 1])
    return dp[n][m]


# ── Левенштейн: расстояние редактирования ────────────────────────────
def edit_distance(a: str, b: str) -> int:
    n, m = len(a), len(b)
    dp = [[0] * (m + 1) for _ in range(n + 1)]
    for i in range(n + 1):
        dp[i][0] = i
    for j in range(m + 1):
        dp[0][j] = j
    for i in range(1, n + 1):
        for j in range(1, m + 1):
            if a[i - 1] == b[j - 1]:
                dp[i][j] = dp[i - 1][j - 1]
            else:
                dp[i][j] = 1 + min(
                    dp[i - 1][j],      # удалить
                    dp[i][j - 1],      # вставить
                    dp[i - 1][j - 1],  # заменить
                )
    return dp[n][m]


# ── LIS: длина наибольшей возрастающей подпоследовательности ─────────
def lis(arr: list[int]) -> int:
    """O(n log n) через терпеливую раскладку."""
    from bisect import bisect_left

    tails: list[int] = []
    for x in arr:
        i = bisect_left(tails, x)
        if i == len(tails):
            tails.append(x)
        else:
            tails[i] = x
    return len(tails)


if __name__ == "__main__":
    print("fib(20) iter:", fib_iter(20))
    print("fib(30) memo:", fib_memo(30))
    print("coin_change [1,3,4], 6:", coin_change([1, 3, 4], 6))
    print("knapsack:", knapsack_01([2, 3, 4, 5], [3, 4, 5, 6], capacity=5))
    print("lcs ABCBDAB / BDCABA:", lcs("ABCBDAB", "BDCABA"))
    print("edit_distance kitten / sitting:", edit_distance("kitten", "sitting"))
    print("lis [10,9,2,5,3,7,101,18]:", lis([10, 9, 2, 5, 3, 7, 101, 18]))
