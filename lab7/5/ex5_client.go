package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Выполняем GET-запрос на маршрут /hello
	getMessage()

	// Выполняем POST-запрос на маршрут /data
	postData()
}

func getMessage() {
	// GET-запрос на /message
	resp, err := http.Get("http://localhost:8080/message")
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ сервера
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Выводим ответ с цитатой из "Тёмной Башни"
	fmt.Println("GET /hello response:", string(body))
}

func postData() {
	// Данные для отправки
	data := []byte("Some data for the /data endpoint")

	// POST-запрос на /data
	resp, err := http.Post("http://localhost:8080/data", "text/plain", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ сервера
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("POST /data response:", string(body))
}
