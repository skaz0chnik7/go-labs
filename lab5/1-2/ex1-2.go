package main

import "fmt"

type Person struct {
	name string
	age  int
}

func main() {
	person := Person{name: "Алексей", age: 30}
	person.birthday()
	person.Info()
}

func (p Person) Info() {
	fmt.Printf("Имя: %s, Возраст: %d\n", p.name, p.age)
}

func (p *Person) birthday() {
	p.age++
}
