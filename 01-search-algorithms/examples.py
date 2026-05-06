"""
BFS и DFS на графе, заданном списком смежности.
Запуск: python examples.py
"""

from collections import deque


def bfs(graph: dict[int, list[int]], start: int) -> list[int]:
    """Возвращает порядок обхода в ширину от start."""
    visited = {start}
    queue = deque([start])
    order = []
    while queue:
        v = queue.popleft()
        order.append(v)
        for u in graph.get(v, []):
            if u not in visited:
                visited.add(u)
                queue.append(u)
    return order


def bfs_shortest_path(graph: dict[int, list[int]], start: int, goal: int) -> list[int] | None:
    """Кратчайший путь в невзвешенном графе. None — если пути нет."""
    if start == goal:
        return [start]
    parents = {start: None}
    queue = deque([start])
    while queue:
        v = queue.popleft()
        for u in graph.get(v, []):
            if u in parents:
                continue
            parents[u] = v
            if u == goal:
                path = [u]
                while parents[path[-1]] is not None:
                    path.append(parents[path[-1]])
                return list(reversed(path))
            queue.append(u)
    return None


def dfs_recursive(graph: dict[int, list[int]], start: int) -> list[int]:
    """Рекурсивный DFS. Осторожно: глубокая рекурсия → RecursionError."""
    visited: set[int] = set()
    order: list[int] = []

    def go(v: int) -> None:
        visited.add(v)
        order.append(v)
        for u in graph.get(v, []):
            if u not in visited:
                go(u)

    go(start)
    return order


def dfs_iterative(graph: dict[int, list[int]], start: int) -> list[int]:
    """Итеративный DFS — без риска переполнить стек."""
    visited = {start}
    stack = [start]
    order = []
    while stack:
        v = stack.pop()
        order.append(v)
        # reversed чтобы порядок совпадал с рекурсивной версией
        for u in reversed(graph.get(v, [])):
            if u not in visited:
                visited.add(u)
                stack.append(u)
    return order


def has_cycle_undirected(graph: dict[int, list[int]]) -> bool:
    """Поиск цикла в неориентированном графе через DFS."""
    visited: set[int] = set()

    def go(v: int, parent: int) -> bool:
        visited.add(v)
        for u in graph.get(v, []):
            if u not in visited:
                if go(u, v):
                    return True
            elif u != parent:
                return True
        return False

    for v in graph:
        if v not in visited and go(v, -1):
            return True
    return False


if __name__ == "__main__":
    #     1 - 2
    #     |   |
    #     3 - 4 - 5
    g = {
        1: [2, 3],
        2: [1, 4],
        3: [1, 4],
        4: [2, 3, 5],
        5: [4],
    }
    print("BFS от 1:", bfs(g, 1))
    print("DFS (рекурсия) от 1:", dfs_recursive(g, 1))
    print("DFS (итеративный) от 1:", dfs_iterative(g, 1))
    print("Кратчайший путь 1 -> 5:", bfs_shortest_path(g, 1, 5))
    print("Есть ли цикл:", has_cycle_undirected(g))
