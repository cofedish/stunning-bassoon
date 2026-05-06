"""
Жадные алгоритмы: классические задачи.
Запуск: python examples.py
"""

import heapq


# ── Activity selection ───────────────────────────────────────────────
def activity_selection(intervals: list[tuple[int, int]]) -> list[tuple[int, int]]:
    """Выбираем максимум непересекающихся интервалов.

    Критерий: брать тот, что заканчивается раньше.
    """
    # сортируем по времени окончания
    sorted_iv = sorted(intervals, key=lambda x: x[1])
    chosen: list[tuple[int, int]] = []
    last_end = float("-inf")
    for s, f in sorted_iv:
        if s >= last_end:
            chosen.append((s, f))
            last_end = f
    return chosen


# ── Размен монет (канонические номиналы) ─────────────────────────────
def coin_change_greedy(coins: list[int], amount: int) -> list[int] | None:
    """Работает только для канонической системы (1, 5, 10, 50, 100, ...).

    Для произвольных монет нужно ДП — см. 07-dynamic-programming.
    """
    coins_desc = sorted(coins, reverse=True)
    result = []
    for c in coins_desc:
        while amount >= c:
            amount -= c
            result.append(c)
    return result if amount == 0 else None


# ── Алгоритм Хаффмана (генерация префиксного кода) ───────────────────
def huffman_codes(freq: dict[str, int]) -> dict[str, str]:
    """Возвращает кодировку: символ → битовая строка."""
    # heap из (частота, уникальный_id, поддерево)
    counter = 0
    heap: list[tuple[int, int, object]] = []
    for ch, f in freq.items():
        heapq.heappush(heap, (f, counter, ch))
        counter += 1

    while len(heap) > 1:
        f1, _, left = heapq.heappop(heap)
        f2, _, right = heapq.heappop(heap)
        heapq.heappush(heap, (f1 + f2, counter, (left, right)))
        counter += 1

    codes: dict[str, str] = {}

    def assign(node, prefix: str) -> None:
        if isinstance(node, str):
            codes[node] = prefix or "0"  # на случай 1 символа
            return
        left, right = node
        assign(left, prefix + "0")
        assign(right, prefix + "1")

    if heap:
        assign(heap[0][2], "")
    return codes


# ── Заправки на трассе ───────────────────────────────────────────────
def min_refuels(stations: list[int], distance: int, tank: int) -> int:
    """Сколько раз минимум придётся заправиться, чтобы проехать distance км.

    stations отсортированы по возрастанию, бак заполнен в начале.
    """
    points = stations + [distance]
    refuels = 0
    current_pos = 0
    last_refuel = 0
    i = 0
    while current_pos + tank < distance:
        # доезжаем до самой дальней досягаемой
        farthest = -1
        while i < len(points) and points[i] <= current_pos + tank:
            farthest = points[i]
            i += 1
        if farthest == -1 or farthest == last_refuel:
            return -1  # не дотягиваем
        current_pos = farthest
        last_refuel = farthest
        refuels += 1
        if current_pos == distance:
            return refuels
    return refuels


if __name__ == "__main__":
    iv = [(1, 4), (3, 5), (0, 6), (5, 7), (3, 9), (5, 9), (6, 10), (8, 11), (8, 12), (2, 14), (12, 16)]
    print("activity_selection:", activity_selection(iv))

    print("coin_change_greedy [1,5,10,25,100], 187:", coin_change_greedy([1, 5, 10, 25, 100], 187))
    print("coin_change_greedy [1,3,4], 6 (тут жадность даёт 3, а оптимум 2!):",
          coin_change_greedy([1, 3, 4], 6))

    print("huffman:", huffman_codes({"a": 5, "b": 9, "c": 12, "d": 13, "e": 16, "f": 45}))

    print("min_refuels:", min_refuels([200, 375, 550, 750, 950], distance=1000, tank=400))
