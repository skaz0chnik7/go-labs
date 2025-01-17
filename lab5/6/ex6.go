package main

import (
	"fmt"
)

type Stringer interface {
	String() string
}

type Book struct {
	title  string
	author string
	year   int
}

// Реализация интерфейса Stringer для Book
func (b Book) String() string {
	return fmt.Sprintf("Книга: \"%s\", Автор: %s, Год: %d", b.title, b.author, b.year)
}

func main() {
	book := Book{title: "Темная Башня", author: "Стивен Кинг", year: 1966}
	fmt.Println(book.String())
}
