package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Функция для генерации случайных чисел и отправки их в канал numbers
func generateNumbers(count int, numbers chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		num := rand.Intn(100)
		numbers <- num
		time.Sleep(time.Millisecond * 500) // Имитация задержки
	}
	close(numbers) // Закрываем канал после генерации всех чисел
}

// Функция для определения чётности числа и отправки сообщения в канал parity
func determineParity(numbers <-chan int, parity chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range numbers {
		if num%2 == 0 {
			parity <- fmt.Sprintf("Число %d чётное", num)
		} else {
			parity <- fmt.Sprintf("Число %d нечётное", num)
		}
	}
	close(parity) // Закрываем канал после обработки всех чисел
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел

	var wg sync.WaitGroup
	numbers := make(chan int)
	parity := make(chan string)

	wg.Add(2) // Две горутины

	go generateNumbers(10, numbers, &wg)
	go determineParity(numbers, parity, &wg)

	// Горутина для вывода результатов с использованием select
	go func() {
		for {
			select {
			case msg, ok := <-parity:
				if !ok {
					return // Канал закрыт, выходим из горутины
				}
				fmt.Println(msg)
			}
		}
	}()

	wg.Wait() // Ожидание завершения генерации и определения чётности
	// Немного подождем, чтобы все сообщения были выведены
	time.Sleep(time.Second)
	fmt.Println("Все операции завершены")
}
