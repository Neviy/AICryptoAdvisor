// package auth- Пакет для реализации логики аутентификации с использованием JWT
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig - Структура для хранения конфигурации JWT
type JWTConfig struct {
	SecretKey string
}

//Claims - Структура для хранения данных, которые будут включены в JWT токен
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// NewJWTConfig - Функция для создания новой конфигурации JWT
func NewJWTConfig(secretKey string) *JWTConfig {
	return &JWTConfig{
		SecretKey: secretKey,
	}
}


// GenerateToken - Метод для генерации JWT токена
func (config *JWTConfig) GenerateToken(email string) (string, error){
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
