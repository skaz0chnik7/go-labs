package main

import (
	database "rest-api/database"
	user_handlers "rest-api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	database.InitDB()

	// Создание маршрутизатора
	router := gin.Default()

	// Маршруты пользователей
	router.GET("/users", user_handlers.GetUsers)
	router.GET("/users/:id", user_handlers.GetUser)
	router.POST("/users", user_handlers.CreateUser)
	router.PUT("/users/:id", user_handlers.UpdateUser)
	router.DELETE("/users/:id", user_handlers.DeleteUser)

	// Запуск сервера
	router.Run(":8080")
}
