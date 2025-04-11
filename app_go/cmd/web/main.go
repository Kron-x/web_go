package main

import (
    "net/http"
    "os"
    "log"
    "os/signal"
    "syscall"
    "time"
    "context"

    "github.com/prometheus/client_golang/prometheus/promhttp"
	"app_go/internal/handlers"   //импорт кастомных хендлеров
    "app_go/pkg/config"
)

func main() {
    config := config.LoadConfig()

    // Настройка логгера
    logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
    defer logFile.Close()
    log.SetOutput(logFile)
    log.Printf("Starting server on port %s", config.Port)

    // 1. Создаём отдельный HTTP-сервер для метрик на порту 8080
    metricsMux := http.NewServeMux()
    metricsMux.Handle("/metrics", promhttp.Handler())
    
    metricsServer := &http.Server{
        Addr:    ":8080",  // Порту 8080 соответствует port-metrics из PodMonitor
        Handler: metricsMux,
    }

    // 2. Основной сервер (как у вас было)
    mainMux := http.NewServeMux()
    mainMux.Handle("/", handlers.LoggingMiddleware(http.HandlerFunc(handlers.HomeHandler)))
    mainMux.Handle("/images/", handlers.LoggingMiddleware(http.StripPrefix("/images/", http.FileServer(http.Dir(config.ImagesDir)))))
    mainMux.HandleFunc("/new-dimension", handlers.NewDimensionHandler)
    mainMux.HandleFunc("/health", handlers.HealthHandler)
    mainMux.HandleFunc("/submit-text", handlers.SubmitTextHandler)

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