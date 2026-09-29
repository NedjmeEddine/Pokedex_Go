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
func appendPokedex(client *pokeapi.Client, mon string, dex *Pokedex) error {
	pokemon, err := client.GetPokemon(&mon)
	if err != nil {
		return err
	}
	dex.pokemons[mon] = *pokemon
	return nil
}
