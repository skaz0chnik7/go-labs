package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	nameMap := map[string]int{
		"Олег":   22,
		"Никита": 27,
		"Сергей": 15,
	}
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Задание 1:")
	fmt.Println("Содержимое карты до изменений: ", nameMap)
	nameMap["Иван"] = 33
	fmt.Println("Содержимое карты после изменений: ", nameMap)

	fmt.Println("\nЗадание 2: ")
	averageAge := getAverageAge(nameMap)
	fmt.Printf("Средний возраст: %.2f лет\n", averageAge)

	fmt.Println("\nЗадание 3: ")
	fmt.Println("Введите имя для удаления:")
	name, _ := reader.ReadString('\n')

	// Убираем символ новой строки
	name = strings.TrimSpace(name)

	// Проверка существования ключа и удаление
	if _, exists := nameMap[name]; exists {
		delete(nameMap, name)
		fmt.Println("Элемент был удален.")
	} else {
		fmt.Println("Элемент с таким именем не найден.")
	}

	// Печать содержимого карты после удаления
	fmt.Println("Содержимое карты после удаления: ", nameMap)

	fmt.Println("\nЗадание 4: ")
	fmt.Println("Введите строку:")
	input, _ := reader.ReadString('\n')
	upper := strings.ToUpper(input)

	fmt.Println("Строка в верхнем регистре:", upper)

	fmt.Println("\nЗадание 5: ")
	var n, sum int
	fmt.Println("Введите количество чисел:")
	fmt.Scan(&n)
	fmt.Println("В4
	ведите числа:")
	for i := 0; i < n; i++ {
		var num int
		fmt.Scan(&num)
		sum += num
	}
	fmt.Println("Сумма введённых чисел:", sum)

	fmt.Println("\nЗадание 6: ")
	fmt.Println("Введите количество элементов массива:")
	fmt.Scan(&n)
	arr := make([]int, n)
	fmt.Println("Введите элементы массива:")
	for i := 0; i < n; i++ {
		fmt.Scan(&arr[i])
	}
	fmt.Println("Массив в обратном порядке:")
	for i := n - 1; i >= 0; i-- {
		fmt.Println(arr[i])
	}
}

func getAverageAge(people map[string]int) float64 {
	var totalAge int
	for _, age := range people {
		totalAge += age
	}
	return float64(totalAge) / float64(len(people))
}
