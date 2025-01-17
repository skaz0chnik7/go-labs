package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Разрешаем соединения от всех источников
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) // Подключенные клиенты
var broadcast = make(chan Message)           // Канал для рассылки сообщений

// Структура сообщения
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	// Настраиваем маршруты
	http.HandleFunc("/ws", handleConnections)

	// Запускаем горутину для обработки сообщений
	go handleMessages()

	// Запускаем сервер
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Функция для обработки подключений
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Обновляем HTTP-соединение до WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Регистрируем нового клиента
	clients[ws] = true

	for {
		var msg Message
		// Читаем новое сообщение как JSON
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error: %v", err)
			delete(clients, ws)
			break
		}
		// Отправляем сообщение в канал для рассылки
		broadcast <- msg
	}
}

// Функция для рассылки сообщений всем клиентам
func handleMessages() {
	for {
		// Получаем сообщение из канала
		msg := <-broadcast
		// Рассылаем его всем клиентам
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
