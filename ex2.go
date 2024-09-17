package main

import (
	"fmt"
)

func main() {
	var integer int = 5
	var boolean bool = true
	var float float64 = 37.85
	var string string = "Строка"
	fmt.Println("Типы данных",
		"\nint: ", integer,
		"\nfloat64: ", float,
		"\nstring: ", string,
		"\nbool: ", boolean)
}
