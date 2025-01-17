package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Выполняем GET-запрос на маршрут /hello
	getHello()

	// Выполняем POST-запрос на маршрут /data
	postData()
}

func getHello() {
	resp, err := http.Get("http://localhost:8080/hello")
	if err != nil {
		fmt.Println("Ошибка GET:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}

	fmt.Println("GET /hello response:", string(body))
}

func postData() {
	// Данные для отправки в формате JSON
	data := map[string]string{"message": "Цитата Достоевского: 'Красота спасет мир'"}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Ошибка направления JSON:", err)
		return
	}

	// Выполняем POST-запрос на маршрут /data
	resp, err := http.Post("http://localhost:8080/data", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка POST:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("POST /data response:", string(body))
}
