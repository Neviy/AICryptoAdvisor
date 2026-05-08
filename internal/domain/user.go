package domain

import (
	"context"
	"errors"
	"time"
)

type UserRepository interface {
	Save(ctx context.Context,user *User) error
	FindByName(ctx context.Context,username string)(*User,error)
	AdminFindByName(ctx context.Context,username string)(*User,error)
	FindByID(ctx context.Context,id int64)(*User,error)
	Update(ctx context.Context,user *User) error
	Delete(ctx context.Context,id int64) error
}
// User представляет пользователя платформы  с полями ID, Email, Username,CreatedAt ,Password,AdminStatus  и Active
type User struct{
	ID int64             `json:"id"`
	UserName string      `json:"username"`
	Active bool          `json:"active"`
	Email string         `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time  `json:"updated_at,omitempty"`
	Password string      `json:"-"`
	AdminStatus bool     `json:"admin_status"`
}

// Инициализация нового пользователя
func NewUser(userName,email,password string,adminStatus bool) (*User,error){
	var now time.Time= time.Now()
	var errs []error
	if userName == ""{
		errs=append(errs, errors.New(" username is required"))
	}
	if email == ""{
		errs=append(errs, errors.New(" email is required"))
	}
	if password == ""{
		errs=append(errs, errors.New(" password is required"))
	}
	if len(errs)>0{
		return nil, errors.Join(errs...)
	}
	return &User{
		UserName: userName,
		Active: true,
		Email: email,
		Password: password,
		AdminStatus: adminStatus,
		CreatedAt: now,
		UpdatedAt: nil,
	}, nil
}

// NewUserFromDB - Инициализация пользователя из данных БД
func NewUserFromDB(id int64,userName,email,password string,adminStatus bool,active bool,createdAt time.Time,updatedAt *time.Time) *User{
	return &User{
		ID: id,
		UserName: userName,
		Active: active,
		Email: email,
		Password: password,
		AdminStatus: adminStatus,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}