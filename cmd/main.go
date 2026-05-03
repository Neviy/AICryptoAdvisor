// main — точка входа приложения, где собираются все зависимости и запускается сервер.
package main

import (
	"AICryptoAdvisor/internal/handler"
	"AICryptoAdvisor/internal/repository"
	"AICryptoAdvisor/internal/service"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Создаём базовый контекст для работы приложения.
	ctx := context.Background()

	// Подключаемся к PostgreSQL через пул соединений.
	conn, err := pgxpool.New(ctx, "postgres://user:password@localhost:5432/dbname")
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем repository слой для работы с БД.
	repo := repository.NewCoinRepository(conn)

	// Инициализируем service слой с бизнес-логикой.
	svc := service.NewService(repo)

	// Инициализируем handler для обработки HTTP-запросов.
	h := handler.NewCoinHandler(svc)

	// Создаём HTTP сервер с помощью Gin.
	r := gin.Default()

	// Регистрируем endpoint для анализа монеты.
	r.GET("/analyze", h.Analyze)

	// Запускаем сервер на порту 8080.
	r.Run(":8080")
}