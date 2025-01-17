// database.go
package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

// Подключение основной БД
func InitDB() {
	var err error
	connStr := "host=localhost port=5432 user=youruser password=pradamgoraj dbname=users sslmode=disable"
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln("Ошибка подключения к БД:", err)
	}
}

// Подключение тестовой БД
func InitTestDB() {
	var err error
	testConnStr := "host=localhost port=5432 user=youruser password=pradamgoraj dbname=test_users sslmode=disable"
	DB, err = sqlx.Connect("postgres", testConnStr)
	if err != nil {
		log.Fatalln("Ошибка подключения к тестовой БД:", err)
	}
}

// Закрытие подключения к БД
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
