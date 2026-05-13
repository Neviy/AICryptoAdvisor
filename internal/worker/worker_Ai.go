// worker/worker_Ai.go — воркер для AI анализа монет. Запускается при старте приложения и выполняет анализ каждые 30 минут.
package worker

import (
	"context"
	"log"
	"time"
)

// UpdateRecommendation — интерфейс для обновления рекомендаций.
type UpdateRecommendation interface{
	UpdateRecommendation(ctx context.Context) error
}
// AiWorker — воркер для AI анализа монет. 
type AiWorker struct {
	update UpdateRecommendation
}
// NewAiWorker — создание нового AiWorker.
func NewAiWorker(update UpdateRecommendation) *AiWorker {
	return &AiWorker{
		update: update,
	}
}

// Start — запуск воркера. Выполняет анализ каждые 30 минут.
func (w *AiWorker) Start(ctx context.Context){
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		if err:=w.update.UpdateRecommendation(ctx);err != nil{
			log.Printf("error update recommendation on start")
		}
	}()
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Printf("AiWorker: start update recommendation")
				if err:=w.update.UpdateRecommendation(ctx);err != nil{
					log.Printf("error update recommendation")
				}
			case <-ctx.Done():
				log.Printf("AiWorker stopped")
				return
			}
		}
	}()
	return nil
}