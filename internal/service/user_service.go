// Покет service хранит бизнес логику для пользователей, он использует интерфейс UserRepository для взаимодействия с базой данных.
package service

import (
	"AICryptoAdvisor/internal/domain"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")



// UserService - сервис для работы с БД пользователей, который использует интерфейс UserRepository для взаимодействия с базой данных.
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService - инициализация структуры UserService
func NewUserService(ur domain.UserRepository) *UserService {
	return &UserService{
		userRepo: ur,
	}
}

// createUser - Создание нового пользователя.
func (s *UserService) createUser(ctx context.Context, username, email, passwordHash string) (*domain.User, error) {
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrInvalidEmail
	}
	user, err := domain.NewUser(username, email, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("constructor error: %w", err)
	}
	err = s.userRepo.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByEmail - Поиск пользователя по email.
func (s *UserService) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser - Обновление данных пользователя.
func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	existing, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}
	return nil
}


// DeleteUser - Удаление пользователя по email.
func (s *UserService) DeleteUser(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	if err := s.userRepo.Delete(ctx, user.ID); err != nil {
		return err
	}
	return nil
}


// RegisterUser - Проверяет, есть ли пользователь, и если его нет, то создаёт пользователя. Также хеширует пароль перед сохранением в базу данных. Возвращает созданного пользователя или ошибку.
func (s *UserService) RegisterUser(ctx context.Context, username, email, passwordHash string) (*domain.User, error) {
	hash,err:=bcrypt.GenerateFromPassword([]byte(passwordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user, err := s.createUser(ctx, username, email, string(hash))
	if err != nil {
		return nil, err
	}
	return user, nil
}

// AuthenticateUser - Аутентификация пользователя.
func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

//ChangePassword - Изменение пароля пользователя.
func (s *UserService) ChangePassword(ctx context.Context, email, newPassword string) (*domain.User, error) {
	user, err := s.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	if err := user.ChangePassword(string(hash)); err != nil {
		return nil, err
	}
	if err := s.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	 return user, nil
}


// DeactivateUser - Деактивация пользователя.
func (s *UserService) DeactivateUser(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if err:=user.Deactivate(); err != nil {
		return nil, err
	}
	if err := s.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// ActivateUser - Активация пользователя.
func (s *UserService) ActivateUser(ctx context.Context, email string) (*domain.User, error) {
		user, err := s.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	user.Activate()
	if err := s.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

//RemoveAdmin- Удаление прав администратора у пользователя.
func (s *UserService) RemoveAdmin(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.FindUserByEmail(ctx, email)	
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if user.Role != "admin" {
		return nil, errors.New("user is not an admin")
	}
	user.PromoteToUser()
	if err := s.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// MakeAdmin- Назначение прав администратора пользователю.
func (s *UserService) MakeAdmin(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if user.Role == "admin" {
		return nil, errors.New("user is already an admin")
	}
	user.PromoteToAdmin()
	if err := s.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
