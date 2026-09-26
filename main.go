package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func main() {
	scaner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scaner.Scan()
		input := scaner.Text()
		if input == "exit" {
			break
		}
		cleanedInput := cleanInput(input)
		fmt.Println("Cleaned input:", cleanedInput[0])
	}
}
