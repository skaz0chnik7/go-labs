package main

import (
	"fmt"
	"sync"
)

// Структура для запроса на выполнение операции
type CalcRequest struct {
	A, B   float64
	Op     string
	Result chan float64
}

// Функция для обработки операций калькулятора
func calculator(requests chan CalcRequest) {
	for req := range requests {
		var res float64
		switch req.Op {
		case "+":
			res = req.A + req.B
		case "-":
			res = req.A - req.B
		case "*":
			res = req.A * req.B
		case "/":
			if req.B != 0 {
				res = req.A / req.B
			} else {
				fmt.Println("Ошибка: деление на ноль")
				res = 0
			}
		default:
			fmt.Println("Ошибка: неверная операция")
			res = 0
		}
		req.Result <- res // Отправляем результат через канал
	}
}

func main() {
	// Канал для отправки запросов калькулятору
	requests := make(chan CalcRequest)
	var wg sync.WaitGroup

	// Запускаем 3 горутины калькулятора
	for i := 0; i < 3; i++ {
		go calculator(requests)
	}

	// Функция для отправки запроса и получения результата
	processRequest := func(a, b float64, op string) {
		resultCh := make(chan float64)
		requests <- CalcRequest{A: a, B: b, Op: op, Result: resultCh} // Отправляем запрос
		result := <-resultCh                                          // Получаем результат
		fmt.Printf("%.2f %s %.2f = %.2f\n", a, op, b, result)
		wg.Done()
	}

	// Запросы
	operations := []struct {
		a, b float64
		op   string
	}{
		{10, 5, "+"},
		{7, 3, "*"},
		{9, 3, "-"},
		{20, 5, "/"},
		{15, 10, "*"},
	}

	// Отправляем запросы параллельно
	for _, op := range operations {
		wg.Add(1)
		go processRequest(op.a, op.b, op.op)
	}

	// Ожидаем завершения всех запросов
	wg.Wait()

	// Закрываем канал, когда все операции завершены
	close(requests)
}
