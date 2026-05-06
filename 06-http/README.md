# 06. Структура HTTP-запросов

HTTP — текстовый протокол поверх TCP. Клиент отправляет запрос, сервер возвращает ответ. Stateless: каждый запрос самодостаточен (если нужно состояние — куки/токены).

## Структура запроса

```
GET /api/users/42?fields=name,email HTTP/1.1     ← стартовая строка
Host: api.example.com                            ← заголовки
Authorization: Bearer eyJhbGciOi...
Accept: application/json
User-Agent: curl/8.4.0
                                                 ← пустая строка
                                                 ← тело (для GET обычно пусто)
```

Стартовая строка: `МЕТОД ПУТЬ ВЕРСИЯ`.

## Структура ответа

```
HTTP/1.1 200 OK                                  ← код + текст
Content-Type: application/json; charset=utf-8    ← заголовки
Content-Length: 89
Date: Tue, 06 May 2026 12:00:00 GMT
                                                 ← пустая строка
{"id": 42, "name": "Den", "email": "d@x.ru"}     ← тело
```

## Методы

| Метод | Назначение | Идемпотентен | Тело в запросе |
|---|---|---|---|
| GET | получить ресурс | да | нет |
| POST | создать новый ресурс / действие | нет | да |
| PUT | заменить ресурс целиком | да | да |
| PATCH | частичное обновление | нет | да |
| DELETE | удалить | да | редко |
| HEAD | как GET, но только заголовки | да | нет |
| OPTIONS | какие методы поддерживает endpoint | да | нет |

«Идемпотентный» = N одинаковых запросов дают тот же эффект, что и один.

## Коды состояния

```
1xx  Информация    редко используется
2xx  Успех
3xx  Перенаправление
4xx  Ошибка клиента
5xx  Ошибка сервера
```

### Запоминать наизусть

| Код | Значение |
|---|---|
| **200** OK | Всё хорошо, тело в ответе |
| **201** Created | Ресурс создан (обычно после POST) |
| **204** No Content | Успех, но тела нет (часто после DELETE) |
| **301** Moved Permanently | Ресурс переехал навсегда |
| **302** Found / 307 / 308 | Временный редирект |
| **304** Not Modified | Кеш ещё актуален |
| **400** Bad Request | Кривое тело/параметры |
| **401** Unauthorized | Не авторизован (нет токена) |
| **403** Forbidden | Авторизован, но нельзя |
| **404** Not Found | Ресурса нет |
| **405** Method Not Allowed | Endpoint есть, но не для этого метода |
| **409** Conflict | Конфликт состояния (дубль, race condition) |
| **422** Unprocessable | Валидация бизнес-правил не прошла |
| **429** Too Many Requests | Rate limit |
| **500** Internal Server Error | Упал сервер |
| **502** Bad Gateway | Прокси не достучалось до апстрима |
| **503** Service Unavailable | Сервис временно недоступен |
| **504** Gateway Timeout | Прокси не дождалось апстрима |

## Важные заголовки

| Заголовок | В запросе/ответе | Что значит |
|---|---|---|
| `Host` | запрос | Доменное имя сервера (обязателен в HTTP/1.1) |
| `User-Agent` | запрос | Кто клиент (браузер, curl, бот) |
| `Accept` | запрос | Какие форматы клиент готов принять |
| `Content-Type` | оба | Формат тела (`application/json`, `text/html`...) |
| `Content-Length` | оба | Размер тела в байтах |
| `Authorization` | запрос | Токен / basic auth |
| `Cookie` / `Set-Cookie` | запрос / ответ | Сессии |
| `Cache-Control` | оба | Правила кеширования |
| `Location` | ответ | Куда редиректить (для 3xx и 201) |
| `ETag` / `If-None-Match` | ответ / запрос | Условный кеш |

## URL — структура

```
https://user:pass@api.example.com:443/v1/users?id=42&active=true#section
└─┬─┘   └────┬───┘ └─────┬───────┘ └┬┘└──┬──┘ └──────┬───────┘ └──┬──┘
схема   userinfo        host      port  path      query        fragment
```

- **Path** — что именно запрашиваем.
- **Query** (после `?`) — параметры через `&`.
- **Fragment** (после `#`) — обрабатывает только клиент, на сервер не отправляется.

## Примеры curl

```bash
# GET
curl https://api.example.com/users/42

# POST с JSON
curl -X POST https://api.example.com/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Den"}'

# с подробным выводом (заголовки запроса/ответа)
curl -v https://api.example.com/

# только заголовки ответа
curl -I https://api.example.com/

# отправить форму
curl -X POST https://api.example.com/login \
  -d "username=den&password=secret"

# скачать файл
curl -O https://example.com/file.zip
```

## HTTPS = HTTP + TLS

То же самое, что HTTP, но в зашифрованном TLS-туннеле. Для клиента отличия почти нет, кроме порта (443 вместо 80) и того, что сервер должен предъявить сертификат.

## Полезные ссылки

- [Гайд по протоколу HTTP (Cloud.ru)](https://cloud.ru/blog/protokol-http)
- [Обзор протокола HTTP (MDN)](https://developer.mozilla.org/ru/docs/Web/HTTP/Overview)
- [HTTP-запросы: структура, методы, коды (Хабр)](https://habr.com/ru/companies/timeweb/articles/697670/)

См. рабочие примеры: [examples.py](./examples.py), [examples.go](./examples.go).
