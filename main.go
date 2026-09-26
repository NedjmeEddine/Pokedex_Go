package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func main() {
	scaner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scaner.Scan()
		input := scaner.Text()
		if input == "exit" {
			commandExit()
		} else if input == "help" {
			fmt.Println("Welcome to the Pokedex! Here are the available commands:")
			fmt.Println("exit - Exit the Pokedex")
			fmt.Println("help - Show a help message")
		} else {
			fmt.Println("Unkown command")
		}
	}
}
