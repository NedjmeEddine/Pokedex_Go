package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/NedjmeEddine/Pokedex_Go/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func REPL(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		command := cleanInput(scanner.Text())
		if len(command) == 0 {
			continue
		}
		commandName := command[0]
		if cmd, ok := commands[commandName]; ok {
			err := cmd.callback(cfg)
			if err != nil {
				fmt.Printf("Error executing command '%s': %v\n", commandName, err)
			}
		} else {
			fmt.Println("Unkown command")
		}
	}
}

func main() {
	cfg := &config{pokeapiClient: pokeapi.NewClient()}
	REPL(cfg)
}
