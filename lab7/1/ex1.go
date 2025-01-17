package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Создание слушателя на порту 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Сервер слушает на порту 8080")

	for {
		// Принимаем входящее соединение
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка принятия подключения:", err)
			continue
		}

		// Обрабатываем соединение в отдельной горутине
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Читаем сообщение от клиента
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Полученное сообщение:", message)

	// Отправляем ответ
	_, err := conn.Write([]byte("Сообщение отправлено\n"))
	if err != nil {
		fmt.Println("Ошибка при отправке ответа:", err)
	}
}
