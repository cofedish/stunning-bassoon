"""
HTTP-клиент в Python.
  - urllib (стандартная библиотека)
  - requests (де-факто стандарт, pip install requests)

Запуск: python examples.py
"""

import json
import urllib.parse
import urllib.request
import urllib.error


# ── 1. Голый urllib (без зависимостей) ────────────────────────────────
def example_urllib() -> None:
    url = "https://httpbin.org/get?lang=ru&topic=http"
    req = urllib.request.Request(
        url,
        headers={"User-Agent": "cheatsheet/1.0"},
        method="GET",
    )
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            print("status:", resp.status)
            print("content-type:", resp.headers["Content-Type"])
            body = resp.read().decode("utf-8")
            print("body:", json.loads(body)["args"])
    except urllib.error.HTTPError as e:
        print(f"HTTP {e.code}: {e.reason}")
    except urllib.error.URLError as e:
        print(f"сеть: {e.reason}")


def example_urllib_post() -> None:
    payload = json.dumps({"name": "Ден"}).encode("utf-8")
    req = urllib.request.Request(
        "https://httpbin.org/post",
        data=payload,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    with urllib.request.urlopen(req, timeout=10) as resp:
        print("posted, echo:", json.loads(resp.read())["json"])


# ── 2. requests — удобнее ─────────────────────────────────────────────
def example_requests() -> None:
    try:
        import requests
    except ImportError:
        print("requests не установлен (pip install requests)")
        return

    # GET с query-параметрами
    r = requests.get(
        "https://httpbin.org/get",
        params={"lang": "ru", "topic": "http"},
        headers={"User-Agent": "cheatsheet/1.0"},
        timeout=10,
    )
    r.raise_for_status()        # бросит HTTPError на 4xx/5xx
    print("requests GET:", r.json()["args"])

    # POST с JSON
    r = requests.post(
        "https://httpbin.org/post",
        json={"name": "Ден", "age": 29},
        timeout=10,
    )
    print("requests POST echo:", r.json()["json"])

    # Сессия (переиспользует TCP-соединение, хранит куки)
    with requests.Session() as s:
        s.headers.update({"Authorization": "Bearer demo-token"})
        r = s.get("https://httpbin.org/headers", timeout=10)
        print("session headers:", r.json()["headers"]["Authorization"])

    # Скачивание файла большими кусками
    with requests.get("https://httpbin.org/bytes/1024", stream=True, timeout=10) as r:
        total = 0
        for chunk in r.iter_content(chunk_size=256):
            total += len(chunk)
        print(f"скачано: {total} байт")


if __name__ == "__main__":
    print("=== urllib GET ===")
    example_urllib()
    print("\n=== urllib POST ===")
    example_urllib_post()
    print("\n=== requests ===")
    example_requests()
