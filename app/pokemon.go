package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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
	Detail         string                `json:"detail"`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func pokedexSerach(term string, url string, update Update) {

	log.Println("searching pokedex: " + term)
	searchURL := "https://pokeapi.co/api/v2/pokemon/" + term
	resp, err := http.Get(searchURL)

	if err != nil {
		log.Println("Error Searching Pokedex")
		sendMessage("Error Searching Pokedex", url, update)
		return
	}

	defer resp.Body.Close()

	r := PokemonSearchResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &r)
	if err != nil {
		log.Println("Error Parsing")
		sendMessage("Error Reading Response From Pokedex", url, update)
		return
	}

	if r.Detail == "Not found." {
		sendMessage(term+" was not found in the Pokedex", url, update)
		return
	}

	returnMessage := "<b>" + strings.ToUpper(r.Name) + "</b>\n<i>"

	// Get the types
	for i := 0; i < len(r.Types); i++ {
		returnMessage += r.Types[i].Type.Name
		if i < len(r.Types)-1 {
			returnMessage += " - "
		}
	}

	// basic info
	returnMessage += " type\n</i>Weight: " + strconv.Itoa(r.Weight) + "\n"
	returnMessage += "Height: " + strconv.Itoa(r.Height) + "\n"
	returnMessage += "Base Exp: " + strconv.Itoa(r.BaseExperience) + "\n"

	// Get the moves
	returnMessage += "\nMoves: <i>"
	numberMovesToList := min(len(r.Moves), 4)
	for i := 0; i < numberMovesToList; i++ {
		returnMessage += r.Moves[i].Move.Name
		if i < numberMovesToList-1 {
			returnMessage += ", "
		}
	}

	if len(r.Moves) > 4 {
		returnMessage += ", and " + strconv.Itoa(len(r.Moves)-4) + " others"
	}

	// Get the moves
	returnMessage += "</i>\n\nAbilities: <i>"
	numberMovesToList = min(len(r.Abilities), 4)
	for i := 0; i < numberMovesToList; i++ {
		returnMessage += r.Abilities[i].Ability.Name
		if i < numberMovesToList-1 {
			returnMessage += ", "
		}
	}

	if len(r.Abilities) > 4 {
		returnMessage += ", and " + strconv.Itoa(len(r.Abilities)-4) + " others"
	}

	returnMessage += "</i>\n\n" + r.Sprites.FrontDefault

	sendMessage(returnMessage, url, update)
}
