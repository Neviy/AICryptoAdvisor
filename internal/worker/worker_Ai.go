// worker/worker_Ai.go — воркер для AI анализа монет. Запускается при старте приложения и выполняет анализ каждые 30 минут.
package worker

import (
	"AICryptoAdvisor/internal/service"
	"context"
	"log"
	"time"
)

// AiWorker — воркер для AI анализа монет.
type AiWorker struct {
	service *service.Service
}
// NewAiWorker — создание нового AiWorker.
func NewAiWorker(service *service.Service) *AiWorker {
	return &AiWorker{
		service: service,
	}
}

// Start — запуск воркера. Выполняет анализ каждые 30 минут.
func (w *AiWorker) Start(ctx context.Context){
	timer := time.NewTicker(30 * time.Minute)
	go func() {
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				log.Printf("AiWorker started analysis")
				coin, err := w.service.AllCoin(ctx)
				if err != nil {
					log.Printf("Error fetching coins: %v", err)
					continue
				}
				for _, c := range coin {
					result, err := w.service.AnalysisCoin(ctx, c.Name)
					if err != nil {
						log.Printf("Error analyzing coin %s: %v", c.Name, err)
						continue
					}
					log.Printf("Analysis result for %s: %s", c.Name, result)
				}
			case <-ctx.Done():
				log.Printf("AiWorker stopped")
				return
			}
		}
	}()
}