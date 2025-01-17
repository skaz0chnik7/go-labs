package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"
)

// setupTestDB создает временную базу данных в памяти для тестирования и создает в ней таблицу пользователей
func setupTestDB() (*sql.DB, func()) {
	// Подключаем временную базу данных SQLite в памяти
	testDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	// Создаем таблицу пользователей для тестирования
	createUserTable := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		age INTEGER,
		password TEXT
	);
	`
	if _, err := testDB.Exec(createUserTable); err != nil {
		log.Fatal(err)
	}

	// Функция очистки базы данных после каждого теста
	cleanup := func() {
		testDB.Exec("DELETE FROM users") // Очистка всех пользователей
		testDB.Close()
	}

	return testDB, cleanup
}

// Тест регистрации нового пользователя
func TestRegister(t *testing.T) {
	testDB, cleanup := setupTestDB()
	defer cleanup()

	// Используем testDB вместо глобальной db для вызовов SQL
	db = testDB

	// Подготовка запроса с данными пользователя
	payload := []byte(`{"name":"testuser","age":30,"password":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(register)

	handler.ServeHTTP(rr, req)

	// Проверка успешного ответа
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("ошибка кода состояния: получили %v, ожидали %v", status, http.StatusCreated)
	}

	var resp map[string]string
	err := json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("не удалось декодировать ответ: %v", err)
	}

	// Проверка наличия токена в ответе
	if _, exists := resp["token"]; !exists {
		t.Errorf("ожидался токен в ответе, получили %v", resp)
	}
}

// Тест входа пользователя в систему
func TestLogin(t *testing.T) {
	testDB, cleanup := setupTestDB()
	defer cleanup()

	// Используем testDB вместо глобальной db для вызовов SQL
	db = testDB

	// Создаем пользователя в базе данных для теста
	_, err := db.Exec("INSERT INTO users (name, age, password) VALUES (?, ?, ?)", "testuser", 30, "password123")
	if err != nil {
		t.Fatalf("не удалось создать пользователя: %v", err)
	}

	// Подготовка запроса с именем пользователя и паролем
	payload := []byte(`{"name":"testuser","password":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(login)

	handler.ServeHTTP(rr, req)

	// Проверка успешного ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ошибка кода состояния: получили %v, ожидали %v", status, http.StatusOK)
	}

	var resp map[string]string
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("не удалось декодировать ответ: %v", err)
	}

	// Проверка наличия токена в ответе
	if _, exists := resp["token"]; !exists {
		t.Errorf("ожидался токен в ответе, получили %v", resp)
	}
}
