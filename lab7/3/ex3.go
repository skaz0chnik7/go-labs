package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Подключаемся к серверу
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка подключения к серверу:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Подключено к серверу. Введите сообщение.")

	// Чтение данных от пользователя и отправка на сервер
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Введите сообщение: ")
		message, _ := reader.ReadString('\n')

		// Отправляем сообщение на сервер
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Ошибка при отправке сообщения:", err)
			return
		}

		// Если пользователь ввел "exit", завершаем работу клиента
		if message == "exit\n" {
			fmt.Println("Завершение работы клиента.")
			return
		}

		// Получаем ответ от сервера
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка при получении ответа от сервера:", err)
			return
		}

		fmt.Print("Ответ сервера: " + response)
	}
}
