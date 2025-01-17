package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Запрашиваем у пользователя имя и роль
	fmt.Println("Введите имя пользователя (например, user или admin):")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Println("Введите пароль (по умолчанию: password):")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Выполняем вход
	token, err := login(username, password)
	if err != nil {
		fmt.Println("Ошибка входа:", err)
		os.Exit(1)
	}
	fmt.Printf("Успешный вход. Токен для %s: %s\n", username, token)

	// Пользователь выбирает маршрут для тестирования
	for {
		fmt.Println("\nВведите маршрут для тестирования (/user или /admin, или 'exit' для выхода):")
		route, _ := reader.ReadString('\n')
		route = strings.TrimSpace(route)

		if route == "exit" {
			fmt.Println("Выход из программы.")
			break
		}

		// Отправляем запрос с токеном
		testRoute(route, token)
	}
}

func login(username, password string) (string, error) {
	url := fmt.Sprintf("http://localhost:9090/login?username=%s&password=%s", username, password)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("Ошибка аутентификации: %s", string(body))
	}

	var result LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Token, nil
}

func testRoute(route, token string) {
	url := "http://localhost:9090" + route
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка выполнения запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return
	}

	fmt.Printf("Статус-код: %d\nОтвет: %s\n", resp.StatusCode, string(body))
}
