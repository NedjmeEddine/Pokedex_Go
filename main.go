package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func main() {
	fmt.Println("Gotta catch 'em all!")
}
