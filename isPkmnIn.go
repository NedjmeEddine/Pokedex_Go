package main

import (
	"github.com/NedjmeEddine/Pokedex_Go/internal/pokeapi"
)

func pokemonInLocation(client *pokeapi.Client, pokemonName string, locationAreaName string) (bool, error) {
	pokemon, err := client.GetPokemonArea(&pokemonName)
	if err != nil {
		return false, err
	}

	for _, encounter := range *pokemon {
		if encounter.LocationArea.Name == locationAreaName {
			return true, nil
		}
	}
	return false, nil
}
