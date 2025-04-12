package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"log"
)

type Message struct {
	ID        int64
	Text      string
	CreatedAt time.Time
}

var db *sql.DB

// Init инициализирует подключение к базе данных
func Init(connStr string) error {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("не удалось проверить подключение к базе данных: %v", err)
	}

	// Создаем таблицу, если она еще не существует
	if err := createTableIfNotExists(); err != nil {
		return fmt.Errorf("не удалось создать таблицу: %v", err)
	}
	log.Println("PostgreSQL connection established successfully")  // Добавляем эту строку
	return nil
}

// Close закрывает подключение к базе данных
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func createTableIfNotExists() error {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`
	_, err := db.Exec(query)
	return err
}

// SaveMessage сохраняет сообщение в базе данных
func SaveMessage(text string) (*Message, error) {
	query := `INSERT INTO messages (text) VALUES ($1) RETURNING id, created_at`
	
	var msg Message
	msg.Text = text
	
	err := db.QueryRow(query, text).Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("не удалось сохранить сообщение: %v", err)
	}
	
	return &msg, nil
}

// CheckConnection проверяет доступность PostgreSQL
func CheckConnection(connStr string) error {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("connection error: %v", err)
    }
    defer db.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return fmt.Errorf("ping error: %v", err)
    }

    return nil
}