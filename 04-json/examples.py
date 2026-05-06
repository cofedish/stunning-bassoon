"""
Базовая работа с JSON в Python через стандартный модуль json.
Запуск: python examples.py
"""

import json
from dataclasses import dataclass, asdict


# ── Сериализация: объект → строка ────────────────────────────────────
data = {
    "name": "Ден",
    "age": 29,
    "skills": ["python", "go"],
    "active": True,
    "address": None,
}

# Минимально (одна строка)
s = json.dumps(data)
print("compact:", s)

# Красиво и без \uXXXX
pretty = json.dumps(data, ensure_ascii=False, indent=2)
print("pretty:\n" + pretty)


# ── Десериализация: строка → объект ──────────────────────────────────
raw = '{"name": "Ден", "age": 29, "skills": ["python", "go"]}'
obj = json.loads(raw)
print("parsed:", obj, type(obj))


# ── Чтение / запись файла ────────────────────────────────────────────
with open("user.json", "w", encoding="utf-8") as f:
    json.dump(data, f, ensure_ascii=False, indent=2)

with open("user.json", "r", encoding="utf-8") as f:
    loaded = json.load(f)
print("from file:", loaded)


# ── Кастомные классы через dataclass ─────────────────────────────────
@dataclass
class User:
    name: str
    age: int
    skills: list[str]


u = User(name="Ден", age=29, skills=["python", "go"])
encoded = json.dumps(asdict(u), ensure_ascii=False)
print("dataclass:", encoded)

decoded_dict = json.loads(encoded)
u2 = User(**decoded_dict)
print("back to obj:", u2)


# ── Не сериализуемые объекты ─────────────────────────────────────────
import datetime


def default(o: object) -> str:
    if isinstance(o, datetime.datetime):
        return o.isoformat()
    raise TypeError(f"{type(o)} is not JSON serializable")


now = {"ts": datetime.datetime(2026, 5, 6, 12, 0)}
print("with custom default:", json.dumps(now, default=default))


# ── Обработка ошибок ─────────────────────────────────────────────────
try:
    json.loads("{name: 'broken'}")
except json.JSONDecodeError as e:
    print(f"невалидный JSON: {e.msg} в позиции {e.pos}")
