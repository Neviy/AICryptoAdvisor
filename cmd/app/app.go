// app - Пакет для реализации основной логики приложения AICryptoAdvisor
package app

import (
	"AICryptoAdvisor/internal/client"
	"AICryptoAdvisor/internal/handler"
	"AICryptoAdvisor/internal/repository"
	"AICryptoAdvisor/internal/service"
	"AICryptoAdvisor/internal/worker"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func Run() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	JWTSecret := os.Getenv("JWT_SECRET")
	ApiKey := os.Getenv("API_KEY")
	if ApiKey == "" {
		return fmt.Errorf("API_KEY is not set in .env file")
	}
	host := os.Getenv("HOST")
	if host == "" {
		return fmt.Errorf("HOST is not set in .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		return fmt.Errorf("PORT is not set in .env file")
	}
	userBase := os.Getenv("USER_BASE")
	if userBase == "" {
		return fmt.Errorf("USER_BASE is not set in .env file")
	}
	password := os.Getenv("PASSWORD")
	if password == "" {
		return fmt.Errorf("PASSWORD is not set in .env file")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return fmt.Errorf("DB_NAME is not set in .env file")
	}
	connStr := "host=" + host + " port=" + port + " user=" + userBase + " password=" + password + " dbname=" + dbName
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Ping(context.Background())
	if err != nil {
		return err
	}
	log.Println("Successfully connected to database")
	userRepo := repository.NewUserRepository(db)
	coinRepo := repository.NewCoinRepository(db)
	aiClient := client.NewGeminiClient(ApiKey)
	coin := service.NewService(coinRepo, aiClient)
	user := service.NewUserService(userRepo)
	handlers := handler.NewCoinHandler(coin)
	userHandler := handler.NewUserHandler(user)
	router := handler.NewRouter(handlers, userHandler)
	aiWorker := worker.NewAiWorker(coin)
	ctx := context.Background()
	aiWorker.Start(ctx)
	return router.Run(":8080")
}