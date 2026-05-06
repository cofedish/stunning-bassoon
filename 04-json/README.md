# 04. Работа с JSON

JSON — текстовый формат обмена данными. Просто структура из объектов, массивов и примитивов.

## Синтаксис

```json
{
  "name": "Den",
  "age": 29,
  "active": true,
  "skills": ["python", "go", "rust"],
  "address": null,
  "meta": {
    "city": "Moscow",
    "tags": ["backend", "ml"]
  }
}
```

### Допустимые типы

| Тип | Пример |
|---|---|
| string | `"hello"` (только двойные кавычки!) |
| number | `42`, `3.14`, `-1e5` |
| boolean | `true`, `false` |
| null | `null` |
| array | `[1, 2, "three"]` |
| object | `{"key": "value"}` |

### Что **не** допускается

- Одинарные кавычки: `'hello'` → ошибка.
- Запятая после последнего элемента: `[1, 2,]` → ошибка.
- Комментарии: `// ...` или `/* */` → ошибка (есть нестандартные диалекты типа JSON5/JSONC, но это не JSON).
- Ключи без кавычек: `{name: "Den"}` — это JS, а не JSON.

## Когда применяется

- REST API: тело запросов и ответов.
- Конфиги: `package.json`, `tsconfig.json`, `composer.json`.
- Логи в structured-формате.
- Хранение данных в NoSQL (MongoDB, etc.).

## Сериализация vs десериализация

- **Сериализация (encode / dump)** — объект в памяти → строка JSON.
- **Десериализация (decode / load)** — строка JSON → объект в памяти.

## Соответствие типов

| JSON | Python | Go |
|---|---|---|
| object | `dict` | `map[string]any` или `struct` |
| array | `list` | `[]any` или `[]T` |
| string | `str` | `string` |
| number (int) | `int` | `int64` / `float64` |
| number (float) | `float` | `float64` |
| true / false | `bool` | `bool` |
| null | `None` | `nil` |

## Типичные ловушки

- **Числа без указания формата в Go** парсятся как `float64`. Целое 1 станет `1.0` после round-trip через `map[string]any`.
- **`json.dumps` с русскими буквами** в Python: по умолчанию экранирует в `\uXXXX`. Передавай `ensure_ascii=False`.
- **Большие числа в JS** теряют точность (число > 2^53). Если в API приходят `int64` — забирай как строку.
- **Безопасность**: никогда не делай `eval` на JSON — используй парсер.

## Полезные ссылки

- [Что такое JSON (Хабр)](https://habr.com/ru/companies/ruvds/articles/513026/)
- [Работа с JSON (MDN Web Docs)](https://developer.mozilla.org/ru/docs/Learn/JavaScript/Objects/JSON)

См. рабочие примеры: [examples.py](./examples.py), [examples.go](./examples.go).
