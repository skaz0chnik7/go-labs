package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Задание 4:")
	var numbers [5]int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		numbers[i] = rand.Intn(50)
		fmt.Println("Массив: arr[", i+1, "] = ", numbers[i], "\n")
	}
	fmt.Println("Задание 5: ")
	srez := numbers[0:2]
	fmt.Println("Срез до: ", srez)
	srez = append(srez, 5)
	fmt.Println("Срез после добавления: ", srez)
	srez = srez[1:]
	fmt.Println("Срез после удаления первого элемента: ", srez)
}
