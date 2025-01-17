package main

import (
	"fmt"
	"sync"
)

// Функция для генерации первых n чисел Фибоначчи и отправки их в канал
func generateFibonacci(n int, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch) // Закрываем канал после отправки всех чисел

	a, b := 0, 1
	for i := 0; i < n; i++ {
		ch <- a
		a, b = b, a+b
	}
	fmt.Println("Горутина генерации Фибоначчи завершена")
}

// Функция для чтения чисел из канала и вывода их на экран
func printNumbers(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch { // Цикл завершается, когда канал закрыт и все данные прочитаны
		fmt.Printf("Получено число: %d\n", num)
	}
	fmt.Println("Горутина печати чисел завершена")
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)

	wg.Add(2) // Две горутины

	go generateFibonacci(10, ch, &wg)
	go printNumbers(ch, &wg)

	wg.Wait() // Ожидание завершения горутин
	fmt.Println("Все горутины завершены")
}
