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

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите идентификатор сессии (нажмите Enter, чтобы выполнить вход): ")
	sessionID, _ := reader.ReadString('\n')
	sessionID = strings.TrimSpace(sessionID)

	var csrfToken string

	if sessionID == "" {
		// Выполняем вход
		fmt.Print("Введите имя пользователя (admin или user): ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("Введите пароль: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		var err error
		sessionID, csrfToken, err = login(username, password)
		if err != nil {
			fmt.Println("Ошибка входа:", err)
			return
		}

		fmt.Println("Вход выполнен успешно.")
		fmt.Println("Ваш идентификатор сессии:", sessionID)
		fmt.Println("Ваш CSRF-токен:", csrfToken)
	} else {
		fmt.Print("Введите CSRF-токен (нажмите Enter, если неизвестен): ")
		csrfToken, _ = reader.ReadString('\n')
		csrfToken = strings.TrimSpace(csrfToken)
	}

	testRoutes(sessionID, csrfToken)
}

func login(username, password string) (string, string, error) {
	url := fmt.Sprintf("http://localhost:8080/login?username=%s&password=%s", username, password)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", "", fmt.Errorf("Ошибка аутентификации: %s", string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var response map[string]string
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", "", err
	}

	return response["session_id"], response["csrf_token"], nil
}

func accessProtectedRoute(method, route, sessionID, csrfToken string) error {
	url := fmt.Sprintf("http://localhost:8080%s?session_id=%s", route, sessionID)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	if method != http.MethodGet && method != http.MethodHead {
		req.Header.Set("X-CSRF-Token", csrfToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Маршрут: %s | Статус: %d | Ответ: %s\n", route, resp.StatusCode, string(body))
	return nil
}

func testRoutes(sessionID, csrfToken string) {
	routes := []struct {
		method, path string
	}{
		{"GET", "/user"},
		{"GET", "/admin"},
	}

	for _, route := range routes {
		fmt.Printf("\nТестируем маршрут %s %s:\n", route.method, route.path)
		_ = accessProtectedRoute(route.method, route.path, sessionID, csrfToken)
	}
}
