package main

import (
	"time"

	"pokedex/internal/pokecache"
)

var (
	LocationAreasURL = "https://pokeapi.co/api/v2/location-area/"
	PokemonURL       = "https://pokeapi.co/api/v2/pokemon/"
)

var (
	commands    commandRegistery
	userPokedex map[string]Pokemon = make(map[string]Pokemon, 0)
	conf        config
)

func init() {
	commands = commandRegistery{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Gets the next 20 Pokemon location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Gets the previous 20 Pokemon location areas",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "List of all the Pokemon located at a certain location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Gives you a chance to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Gives you all the basic information about the given Pokemon if it's catched before",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all the pokemon you have caught before",
			callback:    commandPokedex,
		},
	}

	conf.commands = commands
	conf.Next = ""
	conf.Previous = ""
	conf.Cache = pokecache.NewCache(time.Second * 10)
	conf.Pokedex = userPokedex
}
