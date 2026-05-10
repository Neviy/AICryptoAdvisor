// Покет service хранит бизнес логику для монет, он использует интерфейс CoinRepository для взаимодействия с базой данных.
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
	Delete(ctx context.Context, name string)error
}
// Service - структура чтобы просто было проще работать с интерфейсом и не перегружать функции 
type Service struct{
	coinRepo CoinRepository
}
// NewService - инициализация структуры.
func NewService(or CoinRepository)*Service{
	return &Service{
		coinRepo: or,
	}
}

// CreateCoin-проверяет, есть ли монета, и если её нет, то создаёт монеты. Также устанавливает начальную цену и процент изменения. Возвращает созданную монету или ошибку.
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

// FindCoin-поиск монеты по имени. 
func (s *Service)FindCoin(ctx context.Context,name string)(*domain.Coin,error){
	coin,err:=s.coinRepo.FindByName(ctx,name)
	if err != nil {
		return  nil,err
	}
	return coin,nil
}

//AllCoin- псе монеты.
func (s *Service)AllCoin(ctx context.Context)([]*domain.Coin, error){
	coins,err:=s.coinRepo.GetAll(ctx)
	if err !=nil{
		return nil, err
	}
	return coins, nil
}


// AnalysisCoin-анализ монеты. 
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

// UpdatePrice-обновление цены монеты.
func (s *Service)UpdatePrice(ctx context.Context,name string,price float64)(*domain.Coin,error){
	coin,err:=s.coinRepo.FindByName(ctx,name)
	if err !=nil{
		return nil,err
	}
	if err:=coin.UpdatePrice(price);err !=nil{
		return nil,err
	}
	if err:=s.coinRepo.Update(ctx,coin);err !=nil{
		return nil,err
	}
	return coin,nil
}

// DeleteCoin- удаление монеты.
func (s *Service)DeleteCoin(ctx context.Context,name string)error{
	 coin,err:=s.FindCoin(ctx,name)
	 if err !=nil{
		return err
	}
	if coin== nil{
		return errors.New("there is no such coin")
	}
	if err:=s.coinRepo.Delete(ctx,name);err != nil{
		return err
	}
	return nil
}