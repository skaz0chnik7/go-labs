package main

import "fmt"

func main() {
	var a, b int
	fmt.Scan(&a, &b)
	fmt.Println("Среднее значение двух целых чисел ", a, " и ", b, " = ", average(a, b))
}

func average(first, second int) float64 {
	return float64(first+second) / 2
}
