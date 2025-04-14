package handlers

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
	"app_go/pkg/postgres"
	"app_go/pkg/metrics"
)

// логирования запросов
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // игнорим запрос иконки
		if r.URL.Path == "/favicon.ico" {
            next.ServeHTTP(w, r)
            return
        }
        log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
func SubmitTextHandler(w http.ResponseWriter, r *http.Request) {
	
	metrics.PageVisits.WithLabelValues("/submit-text").Inc()
    metrics.ActiveUsers.Inc()
    defer metrics.ActiveUsers.Dec()

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
	
	
	// Сохраняем сообщение в базу данных

	log.Printf("Saving message to DB: %s", request.Text)

	msg , err := postgres.SaveMessage(request.Text)
	if err != nil {
		log.Printf("DB Save error: %v", err)
		http.Error(w, "Ошибка при сохранении сообщения", http.StatusInternalServerError)
		return
	}
	log.Printf("Message saved with ID: %d", msg.ID)
	
	// Утилизация текста (например, просто возвращаем его обратно)
	response := map[string]string{
		"message": "Вы ввели: " + request.Text,
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func NewDimensionHandler(w http.ResponseWriter, r *http.Request) {

	metrics.PageVisits.WithLabelValues("/new-dimension").Inc()
    metrics.ActiveUsers.Inc()
    defer metrics.ActiveUsers.Dec()

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

// middleware для отслеживания активности
func ActivityMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        metrics.ActiveUsers.Inc()
        defer metrics.ActiveUsers.Dec()
        next.ServeHTTP(w, r)
    })
}