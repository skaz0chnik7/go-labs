package main

import (
	"fmt"
)

func main() {
	strings := []string{"да", "нет", "наверное", "двадцатипятимиллиметровый", "гранат"}
	longest := findLongestString(strings)
	fmt.Println("Самая длинная строка:", longest)
}

func findLongestString(strings []string) string {
	if len(strings) == 0 {
		return ""
	}
	longest := strings[0]
	for _, str := range strings {
		if len(str) > len(longest) {
			longest = str
		}
	}
	return longest
}
