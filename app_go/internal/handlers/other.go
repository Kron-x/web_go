package handlers

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

// логирования запросов
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func SubmitTextHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Чтение JSON-тела запроса
	var request struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Утилизация текста (например, просто возвращаем его обратно)
	response := map[string]string{
		"message": "Вы ввели: " + request.Text,
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func NewDimensionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, `
		<h1>Вы попали в новое измерение, пристегнитесь</h1>
		<p>Текущая дата и время: %s</p>
	    `, currentTime)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
