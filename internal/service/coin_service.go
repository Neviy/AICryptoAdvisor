// Покет service хранит бизнес логику
package service

import (
	"AICryptoAdvisor/internal/domain"
	"context"
	"errors"
	"fmt"
)

// Интерфейс для того , чтобы нормально работать с любой БД
type CoinRepository interface {
	Save(ctx context.Context, coin *domain.Coin) error
	FindByName(ctx context.Context,name string)(*domain.Coin,error)
	GetAll(ctx context.Context) ([]*domain.Coin, error)
	Update(ctx context.Context, coin *domain.Coin)error
}
// Структура чтобы просто было проще работать с интерфейсом и не перегружать функции 
type Service struct{
	coinRepo CoinRepository
}
// инициализация структуры 
func NewService(or CoinRepository)*Service{
	return &Service{
		coinRepo: or,
	}
}

// Проверяет, есть ли монета, и если её нет, то создаёт монеты
func (s *Service)CreateCoin(ctx context.Context,name string,price float64)(*domain.Coin,error){
	existing,err:=s.coinRepo.FindByName(ctx,name)
	if err !=nil {
		return nil,  err
	}
	if existing !=nil{
		return nil,errors.New("coin already exists")
	}
	coin,err:=domain.NewCoin(name,price)
	if err !=nil{
		return nil, fmt.Errorf("constructor error: %w", err)
	}
	if err:=s.coinRepo.Save(ctx,coin);err != nil{
		return nil,  err
	}
	return coin,nil
}

// Поиск монеты по имени 
func (s *Service)FindCoin(ctx context.Context,name string)(*domain.Coin,error){
	coin,err:=s.coinRepo.FindByName(ctx,name)
	if err != nil {
		return  nil,err
	}
	return coin,nil
}
// Все монеты
func (s *Service)AllCoin(ctx context.Context)error{
	coins,err:=s.coinRepo.GetAll(ctx)
	if err !=nil{
		return err
	}
	for i:=0;i<len(coins);i++{
		fmt.Println(coins[i])
	}
	return nil
}

// Анализ монеты 
func (s *Service)AnalysisCoin(ctx context.Context,name string)(*domain.Coin,error){
	coin,err:=s.coinRepo.FindByName(ctx,name)
	if err !=nil{
		return nil,err
	}
	//передает монету GPT для рекомендаций её
	coin.Recommendation="WTF" 
	if err:=s.coinRepo.Update(ctx,coin);err !=nil{
		return nil,err
	}
	return coin,nil 
}

