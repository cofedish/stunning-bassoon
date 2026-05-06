// HTTP-клиент в Go через net/http (стандартная библиотека).
// Запуск: go run examples.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Никогда не используй http.DefaultClient в проде —
// у него нет таймаута, и зависший сервер заморозит твой код навсегда.
var client = &http.Client{Timeout: 10 * time.Second}

// ── GET с параметрами ────────────────────────────────────────────────
func exampleGet() error {
	u, _ := url.Parse("https://httpbin.org/get")
	q := u.Query()
	q.Set("lang", "ru")
	q.Set("topic", "http")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("User-Agent", "cheatsheet/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("status:", resp.Status)
	fmt.Println("content-type:", resp.Header.Get("Content-Type"))

	body, _ := io.ReadAll(resp.Body)
	var parsed struct {
		Args map[string]string `json:"args"`
	}
	_ = json.Unmarshal(body, &parsed)
	fmt.Println("args:", parsed.Args)
	return nil
}

// ── POST с JSON-телом ────────────────────────────────────────────────
func examplePost() error {
	payload, _ := json.Marshal(map[string]any{"name": "Ден", "age": 29})
	req, _ := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer demo-token")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}

	var echo struct {
		JSON map[string]any `json:"json"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&echo)
	fmt.Println("echo:", echo.JSON)
	return nil
}

// ── Простейший HTTP-сервер для понимания «другой стороны» ────────────
func exampleServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var u struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "name": u.Name})
	})

	fmt.Println("сервер запущен на :8080 (раскомментируй вызов в main)")
	_ = http.ListenAndServe(":8080", mux)
}

func main() {
	if err := exampleGet(); err != nil {
		fmt.Println("GET:", err)
	}
	fmt.Println()
	if err := examplePost(); err != nil {
		fmt.Println("POST:", err)
	}
	// exampleServer() // блокирующий — раскомментируй чтобы поднять сервер
}
