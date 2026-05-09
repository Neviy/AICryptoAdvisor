package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserInactive    = errors.New("user is inactive")
)

// User — доменная сущность пользователя
type User struct {
	ID          int64      `json:"id"`
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