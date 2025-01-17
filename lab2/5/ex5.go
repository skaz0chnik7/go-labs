package main

import "fmt"

type Rectangle struct {
	a, b float64
}

func main() {
	var a, b float64 = 5.8, 3
	var rec Rectangle = Rectangle{a, b}
	fmt.Println("Площадь прямоугольника со сторонами ", a, " и ", b, " = ", square(rec))
}

func square(rec Rectangle) float64 {
	return rec.a * rec.b
}
