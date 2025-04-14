package main

import (
    "net/http"
    "os"
    "io"
    "fmt"
    "log"
    "os/signal"
    "syscall"
    "time"
    "context"

    _ "github.com/lib/pq"
    "github.com/prometheus/client_golang/prometheus/promhttp"
	"app_go/internal/handlers"   //импорт кастомных хендлеров
    "app_go/pkg/config"
    "app_go/pkg/postgres"
    "app_go/pkg/metrics"  
)

func main() {

    config := config.LoadConfig() // берем данные из конфига

    metrics.Init() // добавление кастомных метрик в /metrics

    // Настройка логгера
    logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
    defer logFile.Close()
    log.SetOutput(io.MultiWriter(os.Stdout, logFile))

    log.Printf("Starting server on port %s", config.Port)
    
    
    // Инициализация базы данных
    connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
                            os.Getenv("POSTGRES_USER"),
                            os.Getenv("POSTGRES_DB"),
                            os.Getenv("POSTGRES_PASSWORD"),
                            os.Getenv("POSTGRES_HOST"),
                            os.Getenv("POSTGRES_PORT"))
    
    // Ждём готовности PostgreSQL (макс. 30 секунд)
    log.Println("Waiting for PostgreSQL...")
    for i := 0; i < 6; i++ {
        if err := postgres.CheckConnection(connStr); err == nil {
            break
        }
        log.Printf("Attempt %d: PostgreSQL not ready, retrying...", i+1)
        time.Sleep(5 * time.Second)
    }
    
    
    if err := postgres.Init(connStr); err != nil {
        log.Fatalf("Ошибка инициализации базы данных: %v", err)
    }
    defer postgres.Close()
    

    // 1. Создаём отдельный HTTP-сервер для метрик на порту 8080
    metricsMux := http.NewServeMux()
    metricsMux.Handle("/metrics", promhttp.Handler())
    
    metricsServer := &http.Server{
        Addr:    ":8080",  // Порту 8080 соответствует port-metrics из PodMonitor
        Handler: metricsMux,
    }

    // 2. Основной сервер (как было)
    mainMux := http.NewServeMux()

    mainMux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
    mainMux.Handle("/new-dimension", handlers.ActivityMiddleware(handlers.LoggingMiddleware(http.HandlerFunc(handlers.NewDimensionHandler))))
    mainMux.Handle("/submit-text", handlers.ActivityMiddleware(handlers.LoggingMiddleware(http.HandlerFunc(handlers.SubmitTextHandler))))
    mainMux.Handle("/images/", handlers.LoggingMiddleware(http.StripPrefix("/images/", http.FileServer(http.Dir(config.ImagesDir)))))
    mainMux.HandleFunc("/health", handlers.HealthHandler)  
    mainMux.Handle("/", handlers.ActivityMiddleware(handlers.LoggingMiddleware(
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.URL.Path != "/" { // Если путь не корневой
                http.NotFound(w, r)
                return
            }
            handlers.HomeHandler(w, r)
        }),
    )))



    mainServer := &http.Server{
        Addr: ":" + config.Port,
        Handler: mainMux,
    }

    // Запускаем оба сервера
    go func() {
        log.Println("Starting metrics server on :8080")
        if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Metrics server error: %v", err)
        }
    }()

    go func() {
        log.Printf("Starting main server on :%s", config.Port)
        if err := mainServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Main server error: %v", err)
        }
    }()

    // Graceful shutdown для обоих серверов
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    <-stop
    log.Println("Shutting down servers...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Останавливаем серверы
    if err := metricsServer.Shutdown(ctx); err != nil {
        log.Printf("Metrics server shutdown error: %v", err)
    }

    if err := mainServer.Shutdown(ctx); err != nil {
        log.Printf("Main server shutdown error: %v", err)
    }

    log.Println("Servers stopped.")
}