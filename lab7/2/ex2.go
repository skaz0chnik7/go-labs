package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Подключение к серверу
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка подключения к серверу:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Ввод сообщения от пользователя
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("введите сообщение: ")
	message, _ := reader.ReadString('\n')

	// Отправка сообщения на сервер
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Ошибка отправки сообщения:", err)
		return
	}

	// Получение ответа от сервера
	response, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Ответ сервера:", response)
}
