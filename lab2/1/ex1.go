package main

import "fmt"

func main() {
	var number int = 0
	var check int8 = 0
	for {
		fmt.Println("Введите число: ")
		fmt.Scan(&number)
		if number%2 == 0 {
			fmt.Println("Введенное число - четное")
		} else {
			fmt.Println("Введенное число - нечетное")
		}
		fmt.Println("Хотите продолжить?\nДа - 1, нет - 0")
		fmt.Scan(&check)
		if check == 0 {
			break
		}
	}
}
