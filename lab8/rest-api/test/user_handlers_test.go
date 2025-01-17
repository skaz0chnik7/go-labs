package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rest-api/database"
	user_handlers "rest-api/handlers"
	"rest-api/models"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestMain инициализирует тестовую БД один раз для всех тестов
func TestMain(m *testing.M) {
	// Инициализация тестовой базы данных
	database.InitTestDB()
	defer database.CloseDB()
	m.Run()
}

// setupRouter инициализирует маршруты для тестирования
func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/users", user_handlers.GetUsers)
	router.GET("/users/:id", user_handlers.GetUser)
	router.POST("/users", user_handlers.CreateUser)
	router.PUT("/users/:id", user_handlers.UpdateUser)
	router.DELETE("/users/:id", user_handlers.DeleteUser)
	return router
}

// Тест для GetUsers
func TestGetUsers(t *testing.T) {
	router := setupRouter()

	// Создаем хотя бы одного пользователя для проверки
	createTestUser(t, router)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<table") // Проверка, что вернулся HTML-код таблицы
}

// Тест для GetUser
func TestGetUser(t *testing.T) {
	router := setupRouter()

	// Создаем нового пользователя для теста
	user := createTestUser(t, router)

	req, _ := http.NewRequest("GET", "/users/"+strconv.Itoa(user.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<table") // Проверка, что вернулся HTML-код таблицы

	// Удаляем тестового пользователя после теста
	deleteTestUser(t, router, user.ID)
}

// createTestUser создает тестового пользователя и возвращает его
func createTestUser(t *testing.T, router *gin.Engine) models.User {
	user := models.User{
		Name:  "Test User",
		Age:   25,
		Email: "testuser@example.com",
	}
	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Декодируем ответ для получения ID пользователя
	var createdUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.NoError(t, err, "Error unmarshalling response")
	return createdUser
}

// deleteTestUser удаляет пользователя по ID
func deleteTestUser(t *testing.T, router *gin.Engine, userID int) {
	req, _ := http.NewRequest("DELETE", "/users/"+strconv.Itoa(userID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User deleted") // Проверка на успешное удаление
}

// Тест для CreateUser
func TestCreateUser(t *testing.T) {
	router := setupRouter()

	// Создаем и проверяем, что новый пользователь был добавлен
	user := createTestUser(t, router)

	// Удаляем созданного пользователя после теста
	deleteTestUser(t, router, user.ID)
}

// Тест для UpdateUser
func TestUpdateUser(t *testing.T) {
	router := setupRouter()

	// Создаем тестового пользователя для обновления
	user := createTestUser(t, router)

	updatedUser := models.User{
		Name:  "Updated User",
		Age:   30,
		Email: "updateduser@example.com",
	}
	jsonValue, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/users/"+strconv.Itoa(user.ID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated User")

	// Удаляем обновленного пользователя после теста
	deleteTestUser(t, router, user.ID)
}

// Тест для DeleteUser
func TestDeleteUser(t *testing.T) {
	router := setupRouter()

	// Создаем тестового пользователя для удаления
	user := createTestUser(t, router)

	// Выполняем запрос на удаление
	deleteTestUser(t, router, user.ID)
}
