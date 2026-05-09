// Покет service хранит бизнес логику для пользователей, он использует интерфейс UserRepository для взаимодействия с базой данных.
package service

import (
	"AICryptoAdvisor/internal/domain"
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

// Интерфейс для того , чтобы нормально работать с любой БД
type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int64) error
}

// Структура чтобы просто было проще работать с интерфейсом 
type UserService struct {
	userRepo UserRepository
}

// Инициализация структуры
func NewUserService(ur UserRepository) *UserService {
	return &UserService{
		userRepo: ur,
	}
}

// CreateUser - Создание нового пользователя
func (s *UserService) CreateUser(ctx context.Context, username, email, passwordHash string) (*domain.User, error) {
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrInvalidEmail
	}
	user := &domain.User{
		UserName:     username,
		Email:          email,
		Password: passwordHash,
	}
	err = s.userRepo.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByEmail - Поиск пользователя по email
func (s *UserService) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser - Обновление данных пользователя
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

// DeleteUser - Удаление пользователя по email
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