// app - Пакет для реализации основной логики приложения AICryptoAdvisor
package app

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Run() {
	err:=godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	API_KEY := os.Getenv("API_KEY")
	if API_KEY == "" {
		log.Fatal("API_KEY is not set in .env file")
	}
}