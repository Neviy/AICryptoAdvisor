// package service хранит бизнес логику для монет.
package service

import (
	"AICryptoAdvisor/internal/domain"
	"context"
	"errors"
	"fmt"
)

// CoinRepository — интерфейс для работы с БД монет.
type CoinRepository interface {
	Save(ctx context.Context, coin *domain.Coin) error
	FindByName(ctx context.Context, name string) (*domain.Coin, error)
	GetAll(ctx context.Context) ([]*domain.Coin, error)
	Update(ctx context.Context, coin *domain.Coin) error
	Delete(ctx context.Context, name string) error
}

// AIClient — интерфейс AI клиента.
type AIClient interface {
	GetAIRecommendation(ctx context.Context,	prompt string) (string, error)
}

// Service — сервис для работы с монетами.
type Service struct {
	CoinRepo CoinRepository
	AiClient AIClient
}

// NewService — создание сервиса.
func NewService(repo CoinRepository,ai AIClient,) *Service {
	return &Service{
		CoinRepo: repo,
		AiClient: ai,
	}
}

// CreateCoin — создание монеты.
func (s *Service) CreateCoin(ctx context.Context,name string,price float64) (*domain.Coin, error) {
	existing, err := s.CoinRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("coin already exists")
	}
	coin, err := domain.NewCoin(name, price)
	if err != nil {
		return nil, fmt.Errorf(
			"constructor error: %w",
			err,
		)
	}
	if err := s.CoinRepo.Save(ctx, coin); err != nil {
		return nil, err
	}
	return coin, nil
}

// FindCoin — поиск монеты.
func (s *Service) FindCoin(	ctx context.Context,	name string) (*domain.Coin, error) {
	return s.CoinRepo.FindByName(ctx, name)
}
// AllCoin — получение всех монет.
func (s *Service) AllCoin(ctx context.Context,) ([]*domain.Coin, error) {
	return s.CoinRepo.GetAll(ctx)
}

// AnalysisCoin — AI анализ монеты.
func (s *Service) AnalysisCoin(ctx context.Context,name string) (*domain.Coin, error) {
	coin, err := s.CoinRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	prompt := fmt.Sprintf("Coin %s price is %.2f USDT. Give short BUY, SELL or HOLD recommendation.",coin.Name,	coin.Price,)
	recommendation, err := s.AiClient.GetAIRecommendation(
		ctx,
		prompt,
	)
	if err != nil {
		return nil, err
	}
	coin.Recommendation = recommendation
	if err := s.CoinRepo.Update(ctx, coin); err != nil {
		return nil, err
	}
	return coin, nil
}

// UpdatePrice — обновление цены.
func (s *Service) UpdatePrice(ctx context.Context,name string,price float64) (*domain.Coin, error) {
	coin, err := s.CoinRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if err := coin.UpdatePrice(price); err != nil {
		return nil, err
	}
	if err := s.CoinRepo.Update(ctx, coin); err != nil {
		return nil, err
	}
	return coin, nil
}

// DeleteCoin — удаление монеты.
func (s *Service) DeleteCoin(ctx context.Context,name string,) error {
	coin, err := s.FindCoin(ctx, name)
	if err != nil {
		return err
	}
	if coin == nil {
		return errors.New("coin not found")
	}
	return s.CoinRepo.Delete(ctx, name)
}