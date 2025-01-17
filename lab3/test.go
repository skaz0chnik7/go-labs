package main

import "fmt"

func main() {
	nameArray := [6]string{"D", "a", "n", "i", "i", "l"}
	nameSlice := nameArray[:]
	nameSlice = append(nameSlice, "K")
	nameSlice[0] = "K"
	fmt.Println(nameArray, nameSlice)
}
