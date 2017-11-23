package main

// PokemonAbility is a ability certain pokemon can have
type PokemonAbility struct {
	Name string `json:"name"`
}

// PokemonAbilityEntry is a single entry about an ability
type PokemonAbilityEntry struct {
	Ability PokemonAbility `json:"ability"`
}

// PokemonMove is a move certain pokemon can have
type PokemonMove struct {
	Name string `json:"name"`
}

// PokemonMoveEntry an entry about a move a pokemon has
type PokemonMoveEntry struct {
	Move PokemonMove `json:"move"`
}

// PokemonType is a single type
type PokemonType struct {
	Name string `json:"name"`
}

// PokemonTypeEntry Single entry for a type
type PokemonTypeEntry struct {
	Type PokemonType `json:"type"`
}

// PokemonSprites are images associated with the pokemon
type PokemonSprites struct {
	BackFemale       string `json:"back_female"`
	BackShinyFemale  string `json:"back_shiny_female"`
	BackDefault      string `json:"back_default"`
	FrontFemale      string `json:"front_female"`
	FrontShinyFemale string `json:"front_shiny_female"`
	BackShiny        string `json:"back_shiny"`
	FrontDefault     string `json:"front_default"`
	FrontShiny       string `json:"front_shiny"`
}

// PokemonSearchResponse is a response we get from searching a certain pokemon
type PokemonSearchResponse struct {
	Abilities      []PokemonAbilityEntry `json:"abilities"`
	Moves          []PokemonMoveEntry    `json:"moves"`
	Types          []PokemonTypeEntry    `json:"types"`
	Weight         int                   `json:"weight"`
	Name           string                `json:"name"`
	Height         int                   `json:"height"`
	BaseExperience int                   `json:"base_experience"`
	Sprites        PokemonSprites        `json:"sprites"`
}
