// repository/user_memory- Реализация репозитория для пользователей в памяти.
package repository

import (
	"AICryptoAdvisor/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository - реализация репозитория для пользователей в памяти.
type UserRepository struct{
	conn *pgxpool.Pool
}

// NewUserRepository - инициализация структуры UserRepository с переданным соединением к базе данных.
func NewUserRepository(conn *pgxpool.Pool)*UserRepository{
	return &UserRepository{
		conn: conn,
	}
}

// Save-Сохранение данных.
func (ur *UserRepository)Save(ctx context.Context,user *domain.User)error{
	var IdUser int64
	query:=`INSERT INTO users (username, email, password, role, active, coin_id, created_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7)
					RETURNING user_id`
	err:=ur.conn.QueryRow(ctx,query,user.UserName,user.Email,user.Password,user.Role,user.Active,user.IDCoin,user.CreatedAt).Scan(&IdUser)
	if err !=nil{
		return  fmt.Errorf("failed to save user: %w", err)
	}
	user.ID=IdUser
	return nil
}

// FindByEmail-Поиск по email.
func (ur *UserRepository)FindByEmail(ctx context.Context,email string)(*domain.User,error){
	query:=`SELECT user_id, username, email, password, role, active, coin_id, created_at, updated_at
	        FROM users
					WHERE email=$1`
	var id int64
	var username string
	var emailBD string
	var password string
	var role string
	var active bool
	var coinID int64
	var createdAt  time.Time
	var updatedAt *time.Time
	err:=ur.conn.QueryRow(ctx,query,email).Scan(&id,&username,&emailBD,&password,&role,&active,&coinID,&createdAt, &updatedAt)
	if err != nil {
		return  nil,err
	}
	user:=domain.NewUserFromDB(id,username,emailBD,password,role,active,createdAt,coinID, updatedAt)
	return user,nil
}

// Update-Обновление данных в БД.
func (ur *UserRepository)Update(ctx context.Context, user *domain.User)error{
	if user.ID==0{
		return fmt.Errorf("user ID is required")
	}
	query:=`UPDATE users
					SET username=$1, email=$2, password=$3, role=$4, active=$5, coin_id=$6, updated_at=$7
					WHERE user_id=$8`
	_,err:=ur.conn.Exec(ctx,query,user.UserName,user.Email,user.Password,user.Role,user.Active,user.IDCoin,time.Now(),user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete - удаление пользователя из БД.
func (ur *UserRepository)Delete(ctx context.Context, id int64)error{
	if id==0{
		return fmt.Errorf("user ID is required")
	}
	query:=`DELETE FROM users
					WHERE user_id=$1`
	result,err:=ur.conn.Exec(ctx,query,id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if result.RowsAffected()==0{
		return  fmt.Errorf("user not found")
	}
	return nil
}

