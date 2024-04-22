package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var secretKey = []byte("secret_key")

// CreateToken Создать JWT токен с использованием секретного ключа
func CreateToken(username string, password string, expirationTime time.Duration) string {
	// Создать структуру для хранения данных в токене
	claims := jwt.MapClaims{
		"username": username,
		"password": password,
		"exp":      time.Now().Add(expirationTime).Unix(),
	}

	// Создать токен с указанными данными и алгоритмом подписи
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписать токен с использованием секретного ключа
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return ""
	}

	return tokenString
}

func IsValidateToken(tokenString string) bool {
	// Парсинг токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return false // Ошибка парсинга токена
	}

	// Проверка валидности токена
	if !token.Valid {
		return false // Токен недействителен
	}

	// Проверка времени истечения токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expirationTime) {
		return false // Время истекло
	}

	return true // Токен действителен
}
