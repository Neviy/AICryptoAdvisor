// domain/User- хранится структура User и методы связанные с ним.
package domain

import (
	"context"
	"errors"
	"strings"
	"time"
)

// UserRepository - интерфейс для работы с пользователями, который должен реализовать любой репозиторий, работающий с пользователями.
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
}

// Ошибки для валидации данных пользователя.
var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserInactive    = errors.New("user is inactive")
)

// User — доменная сущность пользователя
type User struct {
	ID          int64      `json:"id"`
	IDCoin      int64      `json:"coin_id"`
	UserName    string     `json:"username"`
	Email       string     `json:"email"`
	Password    string     `json:"-"`
	Role        string     `json:"role"`
	Active      bool       `json:"active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// NewUser — создание нового пользователя
func NewUser(
	username string,
	email string,
	passwordHash string,
) (*User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	if len(username) < 3 {
		return nil, ErrInvalidUsername
	}
	if !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}
	if passwordHash == "" {
		return nil, ErrInvalidPassword
	}
	now := time.Now()
	return &User{
		UserName:  username,
		Email:     email,
		Password:  passwordHash,
		Role:      "user",
		Active:    true,
		CreatedAt: now,
	}, nil
}

// NewUserFromDB — восстановление пользователя из БД
func NewUserFromDB(
	id int64,
	username string,
	email string,
	password string,
	role string,
	active bool,
	createdAt time.Time,
	coinID int64,
	updatedAt *time.Time,
) *User {

	return &User{
		ID:        id,
		UserName:  username,
		Email:     email,
		Password:  password,
		Role:      role,
		Active:    active,
		CreatedAt: createdAt,
		IDCoin:    coinID,
		UpdatedAt: updatedAt,
	}
}

// Activate — активация пользователя
func (u *User) Activate() {
	now := time.Now()
	u.Active = true
	u.UpdatedAt = &now
}

// Deactivate — деактивация пользователя
func (u *User) Deactivate() error {
	if u.Role == "admin" {
		return errors.New("cannot deactivate admin")
	}
	now := time.Now()
	u.Active = false
	u.UpdatedAt = &now

	return nil
}

// ChangePassword — смена пароля
func (u *User) ChangePassword(newHash string) error {
	if newHash == "" {
		return ErrInvalidPassword
	}
	now := time.Now()
	u.Password = newHash
	u.UpdatedAt = &now
	return nil
}

// PromoteToAdmin — выдача роли admin
func (u *User) PromoteToAdmin() {
	now := time.Now()
	u.Role = "admin"
	u.UpdatedAt = &now
}

// IsAdmin — проверка роли
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsActive — проверка активности
func (u *User) IsActive() bool {
	return u.Active
}

//PromoteToAdmin — выдача роли user.
func (u *User) PromoteToUser() {
	now := time.Now()
	u.Role = "user"
	u.UpdatedAt = &now
}