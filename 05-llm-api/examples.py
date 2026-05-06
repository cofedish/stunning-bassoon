"""
Минимальный клиент GigaChat без сторонних SDK — на голом requests.

Нужно:
  pip install requests
  export GIGACHAT_AUTH_KEY="<base64(client_id:client_secret)>"
  # ключ выдают в личном кабинете SberDev

Запуск: python examples.py
"""

import os
import time
import uuid
import requests


AUTH_URL = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
CHAT_URL = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"


class GigaChatClient:
    def __init__(self, auth_key: str, scope: str = "GIGACHAT_API_PERS") -> None:
        self.auth_key = auth_key
        self.scope = scope
        self._token: str | None = None
        self._token_expires_at: float = 0.0

    def _get_token(self) -> str:
        # Освежаем за минуту до истечения, чтобы не нарваться на 401
        if self._token and time.time() < self._token_expires_at - 60:
            return self._token

        resp = requests.post(
            AUTH_URL,
            headers={
                "Authorization": f"Basic {self.auth_key}",
                "RqUID": str(uuid.uuid4()),
                "Content-Type": "application/x-www-form-urlencoded",
            },
            data={"scope": self.scope},
            verify=False,  # для прода — путь к сертификату Минцифры
            timeout=10,
        )
        resp.raise_for_status()
        data = resp.json()
        self._token = data["access_token"]
        # expires_at в миллисекундах
        self._token_expires_at = data["expires_at"] / 1000
        return self._token

    def chat(
        self,
        messages: list[dict],
        model: str = "GigaChat",
        temperature: float = 0.7,
        max_tokens: int = 1024,
    ) -> str:
        token = self._get_token()
        resp = requests.post(
            CHAT_URL,
            headers={
                "Authorization": f"Bearer {token}",
                "Content-Type": "application/json",
            },
            json={
                "model": model,
                "messages": messages,
                "temperature": temperature,
                "max_tokens": max_tokens,
            },
            verify=False,
            timeout=60,
        )
        resp.raise_for_status()
        data = resp.json()
        return data["choices"][0]["message"]["content"]


def main() -> None:
    auth_key = os.environ.get("GIGACHAT_AUTH_KEY")
    if not auth_key:
        print("Установи переменную окружения GIGACHAT_AUTH_KEY")
        return

    client = GigaChatClient(auth_key)
    messages = [
        {"role": "system", "content": "Ты лаконичный помощник, отвечай в 1-2 предложения."},
        {"role": "user", "content": "Что такое JSON?"},
    ]
    answer = client.chat(messages)
    print("Ответ:", answer)


if __name__ == "__main__":
    # отключим warning о verify=False для примера
    import urllib3

    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
    main()
