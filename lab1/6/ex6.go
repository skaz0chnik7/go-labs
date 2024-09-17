package main

import (
	"fmt"
)

func main() {
	var first, second, third, result float64 //объявление переменных для слагаемых
	fmt.Println("Введите первое слагаемое: ")
	fmt.Scan(&first)
	fmt.Println("Введите второе слагаемое: ")
	fmt.Scan(&second)
	fmt.Println("Введите третье слагаемое: ")
	fmt.Scan(&third)
	result = (first + second + third) / 3
	fmt.Printf("Среднее трех чисел: %.2f", result) //вывод среднего арифметического с форматированием
}
