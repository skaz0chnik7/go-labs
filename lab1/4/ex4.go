package main

import (
	"fmt"
)

func main() {
	int1 := 15
	int2 := 7
	fmt.Println("Арифметические операции с ", int1, " и ", int2,
		"\nСложение: ", int1+int2,
		"\nВычитание: ", int1-int2,
		"\nУмножение: ", int1*int2,
		"\nДеление: ", int1/int2,
		"\nОстаток от деления: ", int1%int2)
}
