// База данных для хранения информации о монетах
package repository

import (
	"AICryptoAdvisor/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CoinRepository struct {
	conn *pgxpool.Pool
}

func NewCoinRepository(conn *pgxpool.Pool)*CoinRepository{
	return &CoinRepository{
		conn: conn,
	}
}

// Сохранение данных
func (cr *CoinRepository)Save(ctx context.Context,coin *domain.Coin)error{
	var IdCoin int64
	query:=`INSERT INTO coins (name, price, percent, recommendation)
					VALUES ($1, $2, $3, $4)
					RETURNING coin_id`
	err:=cr.conn.QueryRow(ctx,query,coin.Name,coin.Price,coin.Percent,coin.Recommendation).Scan(&IdCoin)
	if err !=nil{
		return  fmt.Errorf("failed to save coin: %w", err)
	}
	coin.ID=IdCoin
	return nil
}

//Поиск по имени 
func (cr *CoinRepository)FindByName(ctx context.Context,name string)(*domain.Coin,error){
	query:=`SELECT coin_id,name,price,percent,recommendation
	        FROM coins
					WHERE name=$1`
	var id int64
	var nameBD string
	var price float64
	var percent float64
	var recommendation string
	err:=cr.conn.QueryRow(ctx,query,name).Scan(&id,&nameBD,&price,&percent,&recommendation)
	if err != nil {
		if errors.Is(err,pgx.ErrNoRows){
			return  nil,nil
		}
		return nil,fmt.Errorf("failed find coin:%w",err)
	}
	coin:=domain.NewCoinFromDB(id,nameBD,price,percent,recommendation)
	return coin,nil
}

//Вывод всех монет 
func (cr *CoinRepository) GetAll(ctx context.Context) ([]*domain.Coin, error) {
	query := `SELECT coin_id, name, price, percent, recommendation
		        FROM coins`
	rows, err := cr.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get coins: %w", err)
	}
	defer rows.Close()
	var coins []*domain.Coin
	for rows.Next() {
		var id int64
		var name string
		var price float64
		var percent float64
		var recommendation string
		if err := rows.Scan(&id, &name, &price, &percent, &recommendation); err != nil {
			return nil, fmt.Errorf("failed to scan coin: %w", err)
		}
		coin := domain.NewCoinFromDB(id, name, price, percent, recommendation)
		coins = append(coins, coin)
	}
	return coins, nil
}

//Обновление данных в БД
func (cr *CoinRepository)Update(ctx context.Context, coin *domain.Coin)error{
	if coin.ID==0{
		return errors.New("coin ID is required")
	}
	query:=`UPDATE coins
	        SET price=$1,percent=$2,recommendation=$3
					WHERE coin_id=$4`
	result,err:=cr.conn.Exec(ctx,query,coin.Price,coin.Percent,coin.Recommendation,coin.ID)
	if err != nil {
		return fmt.Errorf("failed to update coin: %w",err)
	}
	if result.RowsAffected()==0{
		return  fmt.Errorf("coins not found")
	}
	return nil
}