package main

import "fmt"

func main() {
	var number float32 = 0
	var check int8 = 0
	for {
		fmt.Println("Введите число: ")
		fmt.Scan(&number)
		fmt.Println(num(number))
		fmt.Println("Хотите закончить? Введите 0, а иначе любое число: ")
		fmt.Scan(&check)
		if check == 0 {
			break
		}
	}
}

func num(number float32) string {
	if number > 0 {
		return "Positive"
	} else if number == 0 {
		return "Zero"
	} else {
		return "Negative"
	}
}
