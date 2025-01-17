package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

// Мьютекс для синхронизации доступа к соединениям
var mu sync.Mutex

// Мапа для хранения активных соединений
var connections = make(map[net.Conn]struct{})

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Сервер слушает на порту 8080")

	// Канал для завершения работы сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем горутину, которая будет ожидать сигнала завершения
	go func() {
		<-quit
		fmt.Println("\nОстановка сервера...")

		// Закрываем listener, чтобы прекратить прием новых подключений
		listener.Close()

		// Закрываем все активные соединения
		mu.Lock()
		for conn := range connections {
			conn.Close()
		}
		mu.Unlock()

		// Ожидаем завершения всех горутин
		wg.Wait()
		fmt.Println("Сервер корректно завершил работу.")
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			// Проверяем, если listener был закрыт
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				fmt.Println("Слушатель закрыт, завершение работы сервера...")
				return
			}
			fmt.Println("Ошибка при принятии подключения:", err)
			continue
		}

		// Добавляем соедsинение в мапу активных соединений
		mu.Lock()
		connections[conn] = struct{}{}
		mu.Unlock()

		// Добавляем горутину в WaitGroup
		wg.Add(1)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer wg.Done()

	// Убираем соединение из мапы при завершении работы горутины
	defer func() {
		mu.Lock()
		delete(connections, conn)
		mu.Unlock()
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Клиент закрыл соединение.")
				return
			}
			fmt.Println("Ошибка при чтении сообщения:", err)
			return
		}

		fmt.Println("Получено сообщение:", message)

		_, err = conn.Write([]byte("Сообщение получено\n"))
		if err != nil {
			fmt.Println("Ошибка при отправке ответа:", err)
			return
		}
	}
}
