package models

import (
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID    int    `db:"id"`    // Поле для id из таблицы
	Name  string `db:"name"`  // Поле для имени
	Age   int    `db:"age"`   // Поле для возраста
	Email string `db:"email"` // Поле для email
}

var validate = validator.New()

func ValidateUser(user *User) error {
	return validate.Struct(user)
}
