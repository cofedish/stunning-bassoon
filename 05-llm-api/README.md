# 05. Внешние API у БЯМ (LLM)

Большие языковые модели (БЯМ) обычно дёргают через REST API: HTTP POST с JSON-телом, в ответе тоже JSON. Под капотом всё одинаковое — меняется домен и формат тела.

## Общая схема работы

1. Регистрируешься в кабинете провайдера → получаешь **API-ключ** (или client_id + client_secret).
2. Если нужен OAuth2 — отправляешь client_id/secret → получаешь **access_token** (живёт N минут).
3. Делаешь POST-запрос к endpoint `/chat/completions` с массивом сообщений.
4. Парсишь ответ — забираешь `content` ассистента.

## Формат сообщений (везде один и тот же)

```json
{
  "model": "GigaChat",
  "messages": [
    {"role": "system",    "content": "Ты помощник по Python."},
    {"role": "user",      "content": "Объясни декораторы."},
    {"role": "assistant", "content": "Декоратор — это функция..."},
    {"role": "user",      "content": "А с аргументами?"}
  ],
  "temperature": 0.7,
  "max_tokens": 1024
}
```

Поля:
- `model` — какая модель отвечает.
- `messages` — список реплик в хронологическом порядке. Сначала `system` (опц.), потом чередуются `user` / `assistant`.
- `temperature` — креативность. 0 = детерминированно, 1 = разнообразно.
- `max_tokens` — потолок длины ответа.

## GigaChat — российский LLM от Сбера

### Авторизация (OAuth2)

```
POST https://ngw.devices.sberbank.ru:9443/api/v2/oauth
Headers:
  Authorization: Basic <base64(client_id:client_secret)>
  RqUID: <случайный UUID>
  Content-Type: application/x-www-form-urlencoded
Body:
  scope=GIGACHAT_API_PERS    # или GIGACHAT_API_CORP для юрлиц
```

Ответ:
```json
{"access_token": "eyJ...", "expires_at": 1714998000000}
```

Токен живёт **30 минут** — кешируй и обновляй по истечении.

### Запрос к модели

```
POST https://gigachat.devices.sberbank.ru/api/v1/chat/completions
Headers:
  Authorization: Bearer <access_token>
  Content-Type: application/json
Body:
  {модель + messages + параметры}
```

### Сертификат

GigaChat требует доверять корневому сертификату Минцифры. На Linux/macOS его кладут в системный CA-bundle, в Python — указывают через `verify=...`. На время разработки можно `verify=False`, но в проде — никогда.

## Универсальные параметры

| Параметр | Что делает |
|---|---|
| `temperature` (0..1) | Случайность. 0 — повторяемые ответы, 1 — творческие |
| `top_p` (0..1) | Nucleus sampling. Альтернатива temperature |
| `max_tokens` | Максимум токенов в ответе |
| `stream` (bool) | Если true — ответ приходит частями (SSE) |
| `n` | Сколько вариантов ответа сгенерировать |

## Стриминг (Server-Sent Events)

Чтобы текст «печатался» по кусочкам, отправляй `"stream": true`. Ответ приходит как поток строк вида:

```
data: {"choices":[{"delta":{"content":"Дек"}}]}
data: {"choices":[{"delta":{"content":"оратор"}}]}
data: [DONE]
```

## Учёт токенов

Биллинг — за **токены**. 1 токен ≈ 0.75 слова в английском, ~0.5 в русском. В ответе обычно есть `usage`:

```json
"usage": {"prompt_tokens": 42, "completion_tokens": 100, "total_tokens": 142}
```

## Типичные ошибки

| Код | Причина | Что делать |
|---|---|---|
| 401 | Невалидный/просроченный токен | Получить новый |
| 429 | Rate limit | Бэкофф + ретрай |
| 400 | Кривое тело запроса | Проверить JSON и роли в messages |
| 5xx | Проблема на стороне API | Ретрай с экспоненциальным бэкоффом |

## Полезные ссылки

- [Документация GigaChat API (SberDev)](https://developers.sber.ru/docs/ru/gigachat/api/overview)
- [Как получить ключ GigaChat](https://developers.sber.ru/docs/ru/gigachat/individuals-quickstart)

См. рабочие примеры: [examples.py](./examples.py), [examples.go](./examples.go).
