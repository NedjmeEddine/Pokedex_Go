package main

import (
	"fmt"
	"os"

	"github.com/NedjmeEddine/Pokedex_Go/internal/pokeapi"
)

type config struct {
	pokeapiClient   pokeapi.Client
	nextPageURL     *string
	previousPageURL *string
	currentLocation string
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Display this help message",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Display the next page of the map of the Pokedex",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Display the previous page of the map of the Pokedex",
		callback:    commandMapB,
	},
	"explore": {
		name:        "explore",
		description: "fetch the pokemons in a given location area",
		callback:    commandExplore,
	},
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex! Here are the available commands:")
	println("  help - Display this help message")
	println("  exit - Exit the Pokedex")
	return nil
}
func commandMap(cfg *config, args []string) error {
	res, err := cfg.pokeapiClient.GetLocationAreas(cfg.nextPageURL)
	if err != nil {
		fmt.Printf("Error fetching location areas: %v\n", err)
		return err
	}
	for _, locationArea := range res.Results {
		fmt.Printf("%s\n", locationArea.Name)
	}
	cfg.nextPageURL = res.Next
	cfg.previousPageURL = res.Previous
	return nil
}
func commandMapB(cfg *config, args []string) error {
	if cfg.previousPageURL == nil {
		fmt.Println("No previous page available.")
		return nil
	}
	res, err := cfg.pokeapiClient.GetLocationAreas(cfg.previousPageURL)
	if err != nil {
		fmt.Printf("Error fetching location areas: %v\n", err)
		return err
	}
	for _, locationArea := range res.Results {
		fmt.Printf("%s\n", locationArea.Name)
	}
	cfg.nextPageURL = res.Next
	cfg.previousPageURL = res.Previous
	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Please provide a location area name.")
		return nil
	}
	locationName := args[0]
	fmt.Printf("Exploring location area: %s\n", locationName)
	res, err := cfg.pokeapiClient.GetArea(&locationName)
	if err != nil {
		fmt.Printf("Error fetching area: %v\n", err)
		return err
	}
	cfg.currentLocation = locationName
	fmt.Println("Found Pokemon:")
	for _, encounter := range res.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
}
