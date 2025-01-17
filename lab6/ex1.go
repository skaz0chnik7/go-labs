package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Функция для расчета факториала числа n
func factorial(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(1 * time.Second) // Имитация задержки
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
		time.Sleep(100 * time.Millisecond) // Имитация задержки
	}
	fmt.Printf("Факториал %d равен %d\n", n, result)
}

// Функция для генерации случайных чисел
func generateRandomNumbers(count int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second) // Имитация задержки
		num := rand.Intn(100)
		fmt.Printf("Случайное число: %d\n", num)
	}
}

// Функция для вычисления суммы числового ряда
func sumSeries(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(1 * time.Second) // Имитация задержки
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
		time.Sleep(100 * time.Millisecond) // Имитация задержки
	}
	fmt.Printf("Сумма ряда до %d равна %d\n", n, sum)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3) // Три горутины

	go factorial(5, &wg)
	go generateRandomNumbers(5, &wg)
	go sumSeries(10, &wg)

	wg.Wait() // Ожидание завершения всех горутин
	fmt.Println("Все горутины завершены")
}
