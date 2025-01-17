package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Структура для хранения данных сессии
type Session struct {
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Expires   time.Time `json:"expires"`
	CSRFToken string    `json:"csrf_token"`
}

var (
	sessions     = make(map[string]Session) // Хранилище сессий
	sessionsFile = "sessions.json"          // Файл для сохранения сессий
	mu           sync.Mutex                 // Для синхронизации доступа к сессиям
)

// Генерация уникального токена
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Загрузка сессий из файла
func loadSessions() error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Open(sessionsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&sessions)
}

// Сохранение сессий в файл
func saveSessions() error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Create(sessionsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(sessions)
}

// Аутентификация и создание сессии
func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	var role string
	if username == "admin" && password == "password" {
		role = "admin"
	} else if username == "user" && password == "password" {
		role = "user"
	} else {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	sessionID, err := generateToken()
	if err != nil {
		http.Error(w, "Ошибка генерации сессии", http.StatusInternalServerError)
		return
	}

	csrfToken, err := generateToken()
	if err != nil {
		http.Error(w, "Ошибка генерации CSRF токена", http.StatusInternalServerError)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	session := Session{
		Username:  username,
		Role:      role,
		Expires:   expirationTime,
		CSRFToken: csrfToken,
	}

	mu.Lock()
	sessions[sessionID] = session
	mu.Unlock()
	saveSessions()

	response := map[string]string{
		"session_id": sessionID,
		"csrf_token": csrfToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Middleware для проверки сессии
func authorize(next http.HandlerFunc, allowedRoles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.URL.Query().Get("session_id")
		if sessionID == "" {
			http.Error(w, "Сессия не найдена", http.StatusUnauthorized)
			return
		}

		mu.Lock()
		session, exists := sessions[sessionID]
		mu.Unlock()

		if !exists || session.Expires.Before(time.Now()) {
			http.Error(w, "Сессия истекла", http.StatusUnauthorized)
			return
		}

		roleAllowed := false
		for _, role := range allowedRoles {
			if session.Role == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			http.Error(w, "Недостаточно прав", http.StatusForbidden)
			return
		}

		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			csrfToken := r.Header.Get("X-CSRF-Token")
			if csrfToken == "" || csrfToken != session.CSRFToken {
				http.Error(w, "Неверный CSRF токен", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}

// Обработчики маршрутов
func adminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Доступ к данным только для администраторов")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Доступ к данным для пользователей")
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Данные обновлены успешно")
}

func main() {
	err := loadSessions()
	if err != nil {
		log.Fatalf("Ошибка загрузки сессий: %v", err)
	}

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/admin", authorize(adminHandler, "admin"))
	http.HandleFunc("/user", authorize(userHandler, "admin", "user"))
	http.HandleFunc("/update", authorize(updateUserHandler, "user"))

	fmt.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
