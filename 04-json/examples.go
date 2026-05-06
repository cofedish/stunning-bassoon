// Базовая работа с JSON в Go через encoding/json.
// Запуск: go run examples.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Тег `json:"name"` определяет имя поля в JSON.
// `omitempty` пропускает поле, если значение нулевое.
type User struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Skills  []string `json:"skills"`
	Address *string  `json:"address"` // указатель → null когда nil
}

func main() {
	// ── Сериализация ───────────────────────────────────────────────
	u := User{
		Name:    "Ден",
		Age:     29,
		Skills:  []string{"python", "go"},
		Address: nil,
	}

	b, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	fmt.Println("compact:", string(b))

	// С отступами
	pretty, _ := json.MarshalIndent(u, "", "  ")
	fmt.Println("pretty:")
	fmt.Println(string(pretty))

	// ── Десериализация в struct ────────────────────────────────────
	raw := []byte(`{"name":"Ден","age":29,"skills":["python","go"],"address":null}`)
	var parsed User
	if err := json.Unmarshal(raw, &parsed); err != nil {
		panic(err)
	}
	fmt.Printf("parsed struct: %+v\n", parsed)

	// ── Десериализация в map (когда схема неизвестна) ──────────────
	var m map[string]any
	_ = json.Unmarshal(raw, &m)
	fmt.Printf("as map: %v\n", m)
	// внимание: числа в map[string]any парсятся как float64
	if age, ok := m["age"].(float64); ok {
		fmt.Println("age как float64:", age)
	}

	// ── Запись / чтение файла ──────────────────────────────────────
	f, _ := os.Create("user.json")
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	_ = enc.Encode(u)
	_ = f.Close()

	f2, _ := os.Open("user.json")
	defer f2.Close()
	var loaded User
	_ = json.NewDecoder(f2).Decode(&loaded)
	fmt.Printf("from file: %+v\n", loaded)

	// ── Обработка ошибок ───────────────────────────────────────────
	if err := json.Unmarshal([]byte(`{"name": broken}`), &parsed); err != nil {
		fmt.Println("ошибка парсинга:", err)
	}
}
