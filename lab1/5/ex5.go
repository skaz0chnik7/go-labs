package main

import (
	"fmt"
)

func main() {
	var check int8 = 0
	var first, second float64 = 0, 0
	for {
		fmt.Println("Введите первое число: ")
		fmt.Scan(&first)
		fmt.Println("Введите второе число: ")
		fmt.Scan(&second)
		sumAndSub(first, second)
		fmt.Println("Хотите продолжить?\nДа - 1, нет - 0")
		fmt.Scan(&check)
		if check == 0 {
			break
		}
	}
}

func sumAndSub(first float64, second float64) {
	fmt.Println("Сумма: ", first+second, "\nРазность: ", first-second)
}
