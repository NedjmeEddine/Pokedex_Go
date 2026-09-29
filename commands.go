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
	pokedex         Pokedex
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

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex! Here are the available commands:")
	println("  help - Display this help message")
	println("  exit - Exit the Pokedex")
	return nil
}
func commandMap(cfg *config) error {
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
func commandMapB(cfg *config) error {
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

func commandExplore(cfg *config) error {
	locationName := os.Args[2]
	fmt.Printf("Exploring location area: %s\n", locationName)
	res, err := cfg.pokeapiClient.GetArea(&locationName)
	if err != nil {
		fmt.Printf("Error fetching area: %v\n", err)
		return err
	}
	fmt.Println("Found Pokemon: \n")
	for _, encounter := range res.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	appendPokedex(&cfg.pokeapiClient, locationName, &cfg.pokedex)
	return nil
}
