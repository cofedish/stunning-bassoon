"""
Алгоритмы на графах: Дейкстра, топосорт, MST, DSU.
Запуск: python examples.py
"""

import heapq
from collections import deque


# ── Дейкстра: кратчайшие пути от source ──────────────────────────────
def dijkstra(graph: dict[int, list[tuple[int, int]]], source: int) -> dict[int, float]:
    """graph[u] = [(v, w), ...]. Возвращает dist[v] от source."""
    dist = {v: float("inf") for v in graph}
    dist[source] = 0
    pq: list[tuple[float, int]] = [(0, source)]
    while pq:
        d, u = heapq.heappop(pq)
        if d > dist[u]:
            continue
        for v, w in graph.get(u, []):
            nd = d + w
            if nd < dist[v]:
                dist[v] = nd
                heapq.heappush(pq, (nd, v))
    return dist


# ── Топологическая сортировка (Кан) ──────────────────────────────────
def topo_sort(graph: dict[int, list[int]]) -> list[int] | None:
    """None — если есть цикл (топосорт невозможен)."""
    in_deg = {v: 0 for v in graph}
    for u in graph:
        for v in graph[u]:
            in_deg[v] = in_deg.get(v, 0) + 1
            if v not in graph:
                graph[v] = []
    queue = deque([v for v, d in in_deg.items() if d == 0])
    order = []
    while queue:
        u = queue.popleft()
        order.append(u)
        for v in graph[u]:
            in_deg[v] -= 1
            if in_deg[v] == 0:
                queue.append(v)
    return order if len(order) == len(in_deg) else None


# ── DSU ──────────────────────────────────────────────────────────────
class DSU:
    def __init__(self, n: int) -> None:
        self.parent = list(range(n))
        self.rank = [0] * n

    def find(self, x: int) -> int:
        while self.parent[x] != x:
            self.parent[x] = self.parent[self.parent[x]]  # сжатие пути
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rx, ry = self.find(x), self.find(y)
        if rx == ry:
            return False
        if self.rank[rx] < self.rank[ry]:
            rx, ry = ry, rx
        self.parent[ry] = rx
        if self.rank[rx] == self.rank[ry]:
            self.rank[rx] += 1
        return True


# ── Kruskal (MST через DSU) ──────────────────────────────────────────
def kruskal(n: int, edges: list[tuple[int, int, int]]) -> tuple[int, list[tuple[int, int, int]]]:
    """edges = [(u, v, w), ...]. Возвращает (вес MST, рёбра дерева)."""
    dsu = DSU(n)
    total = 0
    tree: list[tuple[int, int, int]] = []
    for u, v, w in sorted(edges, key=lambda e: e[2]):
        if dsu.union(u, v):
            total += w
            tree.append((u, v, w))
    return total, tree


# ── Поиск цикла в ориентированном графе (3-цвета) ───────────────────
def has_cycle_directed(graph: dict[int, list[int]]) -> bool:
    WHITE, GRAY, BLACK = 0, 1, 2
    color = {v: WHITE for v in graph}

    def dfs(u: int) -> bool:
        color[u] = GRAY
        for v in graph.get(u, []):
            if color.get(v, WHITE) == GRAY:
                return True
            if color.get(v, WHITE) == WHITE and dfs(v):
                return True
        color[u] = BLACK
        return False

    return any(color[v] == WHITE and dfs(v) for v in list(graph))


if __name__ == "__main__":
    # взвешенный граф
    g = {
        0: [(1, 4), (2, 1)],
        1: [(3, 1)],
        2: [(1, 2), (3, 5)],
        3: [],
    }
    print("dijkstra от 0:", dijkstra(g, 0))

    # DAG для топосорта
    dag = {
        "compile": ["link"],
        "link": ["test"],
        "test": [],
        "lint": ["test"],
    }
    print("topo:", topo_sort(dag))

    # MST
    edges = [
        (0, 1, 4), (0, 2, 1), (1, 2, 2), (1, 3, 5), (2, 3, 8),
    ]
    print("kruskal:", kruskal(4, edges))

    cyclic = {0: [1], 1: [2], 2: [0]}
    print("has_cycle:", has_cycle_directed(cyclic))
