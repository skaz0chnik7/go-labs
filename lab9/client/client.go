package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

var token string

const baseURL = "http://localhost:9999"

func registerUser() {
	var name string
	var age int
	var password string

	fmt.Print("Введите имя: ")
	fmt.Scanln(&name)
	fmt.Print("Введите возраст: ")
	fmt.Scanln(&age)
	fmt.Print("Введите пароль: ")
	fmt.Scanln(&password)

	data := fmt.Sprintf(`{"name":"%s","age":%d,"password":"%s"}`, name, age, password)
	resp, err := http.Post(baseURL+"/register", "application/json", strings.NewReader(data))
	if err != nil {
		fmt.Println("Ошибка при регистрации пользователя:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Ответ:", string(body))
}

func login() {
	var name string
	var password string

	fmt.Print("Введите имя пользователя для входа: ")
	fmt.Scanln(&name)
	fmt.Print("Введите пароль: ")
	fmt.Scanln(&password)

	data := fmt.Sprintf(`{"name":"%s","password":"%s"}`, name, password)
	resp, err := http.Post(baseURL+"/login", "application/json", strings.NewReader(data))
	if err != nil {
		fmt.Println("Ошибка авторизации:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		token = result["token"]
		fmt.Println("Авторизация успешна. Токен:", token)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка авторизации:", string(body))
	}
}

func getCurrentUser() {
	req, _ := http.NewRequest("GET", baseURL+"/current_user", nil)
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при получении данных текущего пользователя:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var user User
		json.NewDecoder(resp.Body).Decode(&user)
		fmt.Printf("Текущий пользователь: ID: %d, Имя: %s, Возраст: %d\n", user.ID, user.Name, user.Age)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка:", string(body))
	}
}

func listAllUsers() {
	req, _ := http.NewRequest("GET", baseURL+"/users", nil)
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при получении списка пользователей:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var users []User
		json.NewDecoder(resp.Body).Decode(&users)
		for _, user := range users {
			fmt.Printf("ID: %d, Имя: %s, Возраст: %d\n", user.ID, user.Name, user.Age)
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка:", string(body))
	}
}

func updateUser() {
	var userID int
	var newName string
	var newAge int

	fmt.Print("Введите ID пользователя для обновления: ")
	fmt.Scanln(&userID)
	fmt.Print("Введите новое имя: ")
	fmt.Scanln(&newName)
	fmt.Print("Введите новый возраст: ")
	fmt.Scanln(&newAge)

	data := fmt.Sprintf(`{"name":"%s","age":%d}`, newName, newAge)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/users/%d", baseURL, userID), strings.NewReader(data))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при обновлении пользователя:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Ответ:", string(body))
}

func deleteUser() {
	var userID int

	fmt.Print("Введите ID пользователя для удаления: ")
	fmt.Scanln(&userID)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%d", baseURL, userID), nil)
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при удалении пользователя:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Ответ:", string(body))
}

func main() {
	for {
		fmt.Println("╔═════════════════════════════╗")
		fmt.Println("║            МЕНЮ             ║")
		fmt.Println("╠═════════════════════════════╣")
		fmt.Println("║ 1. Зарегистрировать         ║")
		fmt.Println("║ 2. Авторизоваться           ║")
		fmt.Println("║ 3. Текущий пользователь     ║")
		fmt.Println("║ 4. Список всех пользователей║")
		fmt.Println("║ 5. Обновить пользователя    ║")
		fmt.Println("║ 6. Удалить пользователя     ║")
		fmt.Println("║ 7. Выйти                    ║")
		fmt.Println("╚═════════════════════════════╝")

		var choice int
		fmt.Print("Ваш выбор: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			registerUser()
		case 2:
			login()
		case 3:
			getCurrentUser()
		case 4:
			listAllUsers()
		case 5:
			updateUser()
		case 6:
			deleteUser()
		case 7:
			fmt.Println("Завершение работы...")
			return
		default:
			fmt.Println("Некорректный выбор, попробуйте снова.")
		}
	}
}
