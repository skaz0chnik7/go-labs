package user_handlers

import (
	"fmt"
	"net/http"
	"os"
	database "rest-api/database"
	"rest-api/models"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olekukonko/tablewriter"
)

// PrintUsersTable - выводит таблицу с данными пользователей в консоль
func PrintUsersTable(users []models.User) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Age", "Email"})

	for _, user := range users {
		row := []string{
			strconv.Itoa(user.ID),
			user.Name,
			strconv.Itoa(user.Age),
			user.Email,
		}
		table.Append(row)
	}

	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

// generateHTMLTable - создает HTML-код для таблицы с данными пользователей
func generateHTMLTable(users []models.User) string {
	html := "<html><body><table border='1'><tr><th>ID</th><th>Name</th><th>Age</th><th>Email</th></tr>"
	for _, user := range users {
		html += "<tr><td>" + strconv.Itoa(user.ID) + "</td><td>" + user.Name + "</td><td>" + strconv.Itoa(user.Age) + "</td><td>" + user.Email + "</td></tr>"
	}
	html += "</table></body></html>"
	return html
}

// GetUsers - обработчик для получения списка пользователей
func GetUsers(c *gin.Context) {
	name := c.Query("name")
	ageStr := c.Query("age")
	sortBy := c.Query("sort") // параметр сортировки ("name", "age" или "id")
	order := c.Query("order") // порядок сортировки ("asc" или "desc")

	users := []models.User{}
	query := "SELECT * FROM users WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if name != "" {
		query += " AND name ILIKE $" + strconv.Itoa(argIdx)
		args = append(args, "%"+name+"%")
		argIdx++
	}

	if ageStr != "" {
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age format"})
			return
		}
		query += " AND age = $" + strconv.Itoa(argIdx)
		args = append(args, age)
		argIdx++
	}

	err := database.DB.Select(&users, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Сортировка данных
	if sortBy == "name" {
		if order == "desc" {
			sort.Slice(users, func(i, j int) bool { return users[i].Name > users[j].Name })
		} else {
			sort.Slice(users, func(i, j int) bool { return users[i].Name < users[j].Name })
		}
	} else if sortBy == "age" {
		if order == "desc" {
			sort.Slice(users, func(i, j int) bool { return users[i].Age > users[j].Age })
		} else {
			sort.Slice(users, func(i, j int) bool { return users[i].Age < users[j].Age })
		}
	} else if sortBy == "id" {
		if order == "desc" {
			sort.Slice(users, func(i, j int) bool { return users[i].ID > users[j].ID })
		} else {
			sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })
		}
	}

	PrintUsersTable(users) // Табличный вывод в консоль

	// Возвращаем HTML-таблицу на страницу
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(generateHTMLTable(users)))
}

// GetUser - обработчик для получения пользователя по ID
func GetUser(c *gin.Context) {
	id := c.Param("id")
	user := models.User{}
	err := database.DB.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Табличный вывод одного пользователя в консоль
	PrintUsersTable([]models.User{user})

	// Возвращаем HTML-таблицу на страницу
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(generateHTMLTable([]models.User{user})))
}

// CreateUser - обработчик для создания нового пользователя
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec("INSERT INTO users (name, age, email) VALUES ($1, $2, $3)", user.Name, user.Age, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Вывод добавленного пользователя в табличном формате
	fmt.Println("New user added:")
	PrintUsersTable([]models.User{user})

	c.JSON(http.StatusCreated, user)
}

// UpdateUser - обработчик для обновления данных пользователя по ID
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec("UPDATE users SET name=$1, age=$2, email=$3 WHERE id=$4", user.Name, user.Age, user.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Вывод обновленного пользователя
	fmt.Printf("User with ID %s updated:\n", id)
	PrintUsersTable([]models.User{user})

	c.JSON(http.StatusOK, user)
}

// DeleteUser - обработчик для удаления пользователя по ID
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	_, err := database.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Подтверждение удаления в консоли
	fmt.Printf("User with ID %s deleted\n", id)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
