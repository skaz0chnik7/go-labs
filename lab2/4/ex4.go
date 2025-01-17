package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var text string
	fmt.Println("Введите строку: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text = scanner.Text()
	fmt.Println("Введено:", text)
	fmt.Println("Длина введенной строки: ", length(text))
}

func length(length string) int {
	return len([]rune(length))
}
