package main

import (
	"fmt"
	"time"
)

func main() {
	currentTime := time.Now()
	fmt.Println("Дата и время: ", currentTime.Format("2006-01-02 15:04:05"))
}
