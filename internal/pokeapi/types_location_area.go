package pokeapi

type LocationAreasRes struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Area struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel int `json:"min_level"`
				MaxLevel int `json:"max_level"`
				Chance   int `json:"chance"`
				Method   struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				ConditionValues []any `json:"condition_values"`
				PokemonDetails  any   `json:"pokemon_details"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type PokemonLocation []struct {
	LocationArea struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location_area"`
	VersionDetails []struct {
		Version struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
		MaxChance        int `json:"max_chance"`
		EncounterDetails []struct {
			MinLevel int `json:"min_level"`
			MaxLevel int `json:"max_level"`
			Chance   int `json:"chance"`
			Method   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"method"`
			ConditionValues []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"condition_values"`
			PokemonDetails any `json:"pokemon_details"`
		} `json:"encounter_details"`
	} `json:"version_details"`
}
type Pokemon struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"` // decimetres
	Weight int    `json:"weight"` // hectograms

	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`

	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`

	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`

	Cries struct {
		Latest string `json:"latest"`
	} `json:"cries"`

	// URL of a separate endpoint listing the areas where this Pokémon appears
	LocationAreaEncounters string `json:"location_area_encounters"`
}
