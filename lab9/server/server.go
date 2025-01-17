package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

var db *sql.DB
var sessions = make(map[string]int)

func handleError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		age INTEGER,
		password TEXT
	);
	`
	_, err = db.Exec(createUserTable)
	if err != nil {
		log.Fatal(err)
	}
}

func generateToken(userID int) string {
	token := uuid.New().String()
	sessions[token] = userID
	return token
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		userID, exists := sessions[token]
		if !exists || userID == 0 {
			handleError(w, errors.New("не авторизован"), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func getCurrentUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	userID := sessions[token]

	var user User
	err := db.QueryRow("SELECT id, name, age FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		handleError(w, errors.New("пользователь не найден"), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func register(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		handleError(w, errors.New("неверный формат данных"), http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Age <= 0 || user.Password == "" {
		handleError(w, errors.New("неверные данные пользователя"), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO users (name, age, password) VALUES (?, ?, ?)", user.Name, user.Age, user.Password)
	if err != nil {
		handleError(w, errors.New("не удалось создать пользователя"), http.StatusInternalServerError)
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		handleError(w, errors.New("ошибка при получении ID пользователя"), http.StatusInternalServerError)
		return
	}

	token := generateToken(int(userID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "пользователь успешно зарегистрирован", "token": token})
}

func login(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil || data.Name == "" || data.Password == "" {
		handleError(w, errors.New("неверный формат запроса"), http.StatusBadRequest)
		return
	}

	var userID int
	var storedPassword string

	err = db.QueryRow("SELECT id, password FROM users WHERE name = ?", data.Name).Scan(&userID, &storedPassword)
	if err != nil || storedPassword != data.Password {
		handleError(w, errors.New("неверное имя пользователя или пароль"), http.StatusUnauthorized)
		return
	}

	token := generateToken(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	userID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, errors.New("некорректный ID пользователя"), http.StatusBadRequest)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Name == "" || user.Age <= 0 {
		handleError(w, errors.New("некорректные данные"), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET name = ?, age = ? WHERE id = ?", user.Name, user.Age, userID)
	if err != nil {
		handleError(w, errors.New("ошибка обновления пользователя"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Пользователь обновлен")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	userID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, errors.New("некорректный ID пользователя"), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		handleError(w, errors.New("не удалось удалить пользователя"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Пользователь удален")
}

func main() {
	initDB()
	defer db.Close()
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/current_user", authMiddleware(getCurrentUser))
	http.HandleFunc("/users", authMiddleware(listUsers))

	// Используем разные HTTP-методы для путей с "/users/{id}"
	http.Handle("/users/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			updateUser(w, r)
		case "DELETE":
			deleteUser(w, r)
		default:
			handleError(w, errors.New("метод не поддерживается"), http.StatusMethodNotAllowed)
		}
	})))

	fmt.Println("Сервер запущен на порту 9999...")
	log.Fatal(http.ListenAndServe(":9999", nil))
}
