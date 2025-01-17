package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	// Подключаемся к серверу
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	// Запрашиваем имя пользователя один раз
	fmt.Print("Enter your username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username) // Убираем пробелы и символы новой строки

	// Запускаем горутину для получения сообщений
	go receiveMessages(conn)

	// Ввод сообщений пользователем
	for {
		fmt.Print("Enter your message: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message) // Убираем символы новой строки

		// Отправляем сообщение на сервер
		msg := map[string]string{
			"username": username,
			"message":  message,
		}
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Println("Error writing to WebSocket:", err)
			return
		}
	}
}

// Функция для приёма сообщений от сервера
func receiveMessages(conn *websocket.Conn) {
	for {
		// Читаем сообщение с сервера
		var msg map[string]string
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading from WebSocket:", err)
			return
		}
		fmt.Printf("[%s]: %s\n", msg["username"], msg["message"])
	}
}
