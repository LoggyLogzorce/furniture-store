package api

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var secretKey = []byte("secret_key")

// CreateToken Создать JWT токен с использованием секретного ключа
func CreateToken(login string, role string, expirationTime time.Duration) string {
	// Создать структуру для хранения данных в токене
	claims := jwt.MapClaims{
		"login": login,
		"role":  role,
		"exp":   time.Now().Add(expirationTime).Unix(),
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
	token, err := ParseToken(tokenString)
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

func GetRoleFromToken(tokenString string) (string, error) {
	// Парсинг JWT токена
	token, err := ParseToken(tokenString)

	// Проверка на ошибку при парсинге токена
	if err != nil {
		return "", err
	}

	// Проверка наличия утверждений (claims) в токене
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Получение роли пользователя из утверждений
		if role, exists := claims["role"].(string); exists {
			return role, nil
		}
	}

	// Если роль не найдена в токене, возвращаем ошибку
	return "", errors.New("role not found in token")
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func GetLoginFromToken(tokenString string) (string, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if login, exists := claims["login"].(string); exists {
			return login, nil
		}
	}
	return "", errors.New("login not found in token")
}
