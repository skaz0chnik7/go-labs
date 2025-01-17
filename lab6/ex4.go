package main

import (
	"fmt"
	"sync"
)

var (
	counterWithMutex    int
	counterWithoutMutex int
	mutex               sync.Mutex
)

// Вариант с использованием мьютекса
func incrementWithMutex(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10000; i++ {
		mutex.Lock()       // Закрываем доступ к общему ресурсу
		counterWithMutex++ // Увеличиваем счётчик
		mutex.Unlock()     // Открываем доступ к ресурсу
	}
}

// Вариант без использования мьютекса
func incrementWithoutMutex(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10000; i++ {
		counterWithoutMutex++ // Увеличиваем счётчик
	}
}

func main() {
	var wg sync.WaitGroup

	// Запуск варианта с мьютексом
	for i := 0; i < 5; i++ { // Запускаем 5 горутин
		wg.Add(1)
		go incrementWithMutex(&wg)
	}

	wg.Wait() // Ожидаем завершения всех горутин
	fmt.Println("Итоговое значение счётчика (с мьютексом):", counterWithMutex)

	// Сбрасываем счетчик и ждем новые горутины для варианта без мьютекса
	wg = sync.WaitGroup{} // Сбрасываем состояние wait group

	// Запуск варианта без мьютекса
	for i := 0; i < 5; i++ { // Запускаем 5 горутин
		wg.Add(1)
		go incrementWithoutMutex(&wg)
	}

	wg.Wait() // Ожидаем завершения всех горутин
	fmt.Println("Итоговое значение счётчика (без мьютекса):", counterWithoutMutex)
}
