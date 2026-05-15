// main - Точка входа в приложение AICryptoAdvisor
package main

import (
	"AICryptoAdvisor/cmd/app"
	"log"
)

func main() {
	if err:=app.Run();err!=nil{
		log.Fatalf("Failed to run application: %v", err)
	}
}