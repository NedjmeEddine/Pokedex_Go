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
	"catch": {
		name:        "catch",
		description: "catch a pokemon in the current location area",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "inspect a pokemon in the pokedex",
		callback:    commandInspect,
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
		fmt.Println("Please provide a location area name to explore.")
		return nil
	}
	locationName := args[0]
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
	cfg.currentLocation = locationName
	return nil
}
func commandCatch(cfg *config, args []string) error {
	if cfg.currentLocation == "" {
		fmt.Println("You need to explore a location area first using the 'explore' command.")
		return nil
	}
	if len(args) < 1 {
		fmt.Println("Please provide a Pokemon name to catch.")
		return nil
	}
	pokemonName := args[0]
	if _, caught := cfg.pokedex.pokemons[pokemonName]; caught {
		fmt.Printf("%s is already in your Pokedex.\n", pokemonName)
		return nil
	}
	if inArea, err := pokemonInLocation(&cfg.pokeapiClient, pokemonName, cfg.currentLocation); err != nil {
		fmt.Printf("Error checking if Pokemon is in location: %v\n", err)
		return err
	} else if !inArea {
		fmt.Printf("%s is not found in %s.\n", pokemonName, cfg.currentLocation)
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s in %s...\n", pokemonName, cfg.currentLocation)
	if err := appendPokedex(&cfg.pokeapiClient, pokemonName, &cfg.pokedex); err != nil {
		fmt.Printf("Error catching %s: %v\n", pokemonName, err)
		return err
	}
	fmt.Printf("Congratulations! You caught %s!\n", pokemonName)
	return nil
}
func commandInspect(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Please provide a Pokemon name to inspect.")
		return nil
	}
	pokemonName := args[0]
	pokemon, caught := cfg.pokedex.pokemons[pokemonName]
	if !caught {
		fmt.Printf("%s is not in your Pokedex. Catch it first!\n", pokemonName)
		return nil
	}
	fmt.Print(pokemonCard(pokemon))
	return nil
}
