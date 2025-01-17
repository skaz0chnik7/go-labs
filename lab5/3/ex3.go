package main

import (
	"fmt"
	"math"
)

type Circle struct {
	radius float64
}

func main() {
	circle := Circle{radius: 5}
	fmt.Printf("Площадь круга: %.2f\n", circle.Area())
}

// Метод для вычисления площади круга
func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}
