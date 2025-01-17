package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

// Функция для реверсирования строки
func reverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// Воркер для обработки задач
func worker(id int, tasks <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("Воркер %d начал задачу: %s\n", id, task)
		time.Sleep(500 * time.Millisecond) // Имитируем задержку
		results <- reverseString(task)
		fmt.Printf("Воркер %d завершил задачу: %s\n", id, task)
	}
}

func main() {
	var numWorkers int

	// Ввод количества воркеров от пользователя
	fmt.Print("Введите количество воркеров: ")
	fmt.Scan(&numWorkers)

	// Открытие файла для чтения строк
	inputFile, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer inputFile.Close()

	// Открытие файла для записи результатов
	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outputFile.Close()

	tasks := make(chan string, 10)   // Канал для задач
	results := make(chan string, 10) // Канал для результатов

	var wg sync.WaitGroup

	// Запуск воркеров
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	// Чтение строк из файла и отправка их в канал задач
	go func() {
		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			tasks <- scanner.Text() // Отправляем строки на обработку
		}
		close(tasks) // Закрываем канал задач после отправки всех строк
		if err := scanner.Err(); err != nil {
			fmt.Println("Ошибка при чтении файла:", err)
		}
	}()

	// Ожидание завершения всех воркеров
	go func() {
		wg.Wait()
		close(results) // Закрываем канал результатов после завершения воркеров
	}()

	// Получение результатов и запись их в файл
	for result := range results {
		fmt.Println("Результат:", result) // Можно оставить вывод в консоль
		_, err := outputFile.WriteString(result + "\n")
		if err != nil {
			fmt.Println("Ошибка при записи в файл:", err)
		}
	}

	fmt.Println("Все задачи завершены. Результаты записаны в output.txt.")
}
