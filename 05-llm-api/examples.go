// Минимальный клиент GigaChat на голом net/http.
// Запуск:
//   export GIGACHAT_AUTH_KEY="<base64(client_id:client_secret)>"
//   go run examples.go
package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	authURL = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	chatURL = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"` // миллисекунды
}

type GigaChatClient struct {
	authKey   string
	scope     string
	http      *http.Client
	token     string
	expiresAt time.Time
}

func NewGigaChatClient(authKey string) *GigaChatClient {
	// Для прода: указать корневой сертификат Минцифры через RootCAs.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &GigaChatClient{
		authKey: authKey,
		scope:   "GIGACHAT_API_PERS",
		http:    &http.Client{Timeout: 60 * time.Second, Transport: tr},
	}
}

func (c *GigaChatClient) getToken() (string, error) {
	if c.token != "" && time.Now().Before(c.expiresAt.Add(-time.Minute)) {
		return c.token, nil
	}

	form := url.Values{"scope": {c.scope}}
	req, _ := http.NewRequest("POST", authURL, strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "Basic "+c.authKey)
	req.Header.Set("RqUID", uuid.NewString())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("auth %d: %s", resp.StatusCode, body)
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}
	c.token = tr.AccessToken
	c.expiresAt = time.UnixMilli(tr.ExpiresAt)
	return c.token, nil
}

func (c *GigaChatClient) Chat(messages []Message) (string, error) {
	token, err := c.getToken()
	if err != nil {
		return "", err
	}

	body, _ := json.Marshal(chatRequest{
		Model:       "GigaChat",
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   1024,
	})
	req, _ := http.NewRequest("POST", chatURL, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("chat %d: %s", resp.StatusCode, b)
	}

	var cr chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return "", err
	}
	if len(cr.Choices) == 0 {
		return "", fmt.Errorf("пустой ответ")
	}
	return cr.Choices[0].Message.Content, nil
}

func main() {
	authKey := os.Getenv("GIGACHAT_AUTH_KEY")
	if authKey == "" {
		fmt.Println("Установи GIGACHAT_AUTH_KEY")
		return
	}

	client := NewGigaChatClient(authKey)
	answer, err := client.Chat([]Message{
		{Role: "system", Content: "Ты лаконичный помощник, отвечай в 1-2 предложения."},
		{Role: "user", Content: "Что такое JSON?"},
	})
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("Ответ:", answer)
}
