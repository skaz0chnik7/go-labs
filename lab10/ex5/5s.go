package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Конфигурация секретного ключа для JWT
var jwtKey = []byte("my_secret_key")

// Claims структура для JWT
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Генерация JWT-токена
func generateJWT(username, role string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // Токен действует 15 минут
	claims := &Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func authorize(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"message": "Токен не предоставлен"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"message": "Невалидный токен"}`, http.StatusUnauthorized)
			return
		}

		if requiredRole != "" {
			// Доступ администратора к маршрутам пользователей
			if claims.Role != requiredRole && !(claims.Role == "admin" && requiredRole == "user") {
				http.Error(w, `{"message": "Недостаточно прав"}`, http.StatusForbidden)
				return
			}
		}

		next(w, r)
	}
}

// Обработчик для входа
func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	var token string
	var err error

	if username == "admin" && password == "password" {
		token, err = generateJWT(username, "admin")
	} else if username == "user" && password == "password" {
		token, err = generateJWT(username, "user")
	} else {
		http.Error(w, `{"message": "Неверные учетные данные"}`, http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, `{"message": "Ошибка создания токена"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Обработчик для данных пользователей
func userHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Доступ к данным для пользователей"}
	json.NewEncoder(w).Encode(response)
}

// Обработчик для данных администраторов
func adminHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Доступ к данным только для админов"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/user", authorize("user", userHandler))
	http.HandleFunc("/admin", authorize("admin", adminHandler))

	fmt.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
