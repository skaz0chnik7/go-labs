package main

import (
	"1-3/mathutils"
	"1-3/stringutils"
	"bufio"
	"fmt"
	"os"
)

func main() {
	var number int = 5
	fmt.Println("Задание 1-2: ")
	fmt.Println("Введите число: ")
	fmt.Scanln(&number)
	fmt.Println("Факториал: ", mathutils.Factorial(number), "\n")

	fmt.Println("Задание 3: ")
	var text string
	fmt.Println("Введите строку: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text = scanner.Text()
	fmt.Println("Введено:", text)
	fmt.Println("Перевернутая строка: ", stringutils.Reverse(text))
}
