package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	result := strings.ToLower(text)
	return strings.Fields(result)
}

func isCommand(prompt []string) bool {
	found := false
	for _, value := range prompt {
		for key := range commands {
			if key == value {
				found = true
				break
			}
		}
	}

	return found
}

func commandExit(conf *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(conf *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")

	for key, value := range commands {
		fmt.Printf("%s: %s\n", key, value.description)
	}

	return nil
}

func fetchData(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d", res.StatusCode)
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading data: %v", err)
	}

	return raw, nil
}

func printLocations(locations []location) {
	for _, loc := range locations {
		fmt.Printf("%v\n", loc.Name)
	}
}

func getLocationsAndPrint(conf *config, url string) error {
	raw, err := fetchData(url)
	if err != nil {
		return err
	}
	conf.Cache.Add(url, raw)
	return decodeAndPrint(conf, raw)
}

func decodeAndPrint(conf *config, raw []byte) error {
	var data LocationAreaResponse
	if err := json.Unmarshal(raw, &data); err != nil {
		return fmt.Errorf("error decoding data: %v", err)
	}
	conf.Next = data.Next
	conf.Previous = data.Previous
	printLocations(data.Results)
	return nil
}

func commandMap(conf *config, args ...string) error {
	url := conf.Next
	if url == "" {
		url = LocationAreasURL
	}

	raw, ok := conf.Cache.Get(url)
	if !ok {
		return getLocationsAndPrint(conf, url)
	}

	return decodeAndPrint(conf, raw)
}

func commandMapB(conf *config, args ...string) error {
	if len(conf.Previous) == 0 {
		return fmt.Errorf("you're on the first page")
	}

	raw, ok := conf.Cache.Get(conf.Previous)
	if !ok {
		return getLocationsAndPrint(conf, conf.Previous)
	}

	return decodeAndPrint(conf, raw)
}

func decodeLocationAreaAndPrint(raw []byte) error {
	var data LocationArea
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	for _, pokemon := range data.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}

func exploreLocationArea(conf *config, url string) error {
	raw, err := fetchData(url)
	if err != nil {
		return err
	}

	conf.Cache.Add(url, raw)
	return decodeLocationAreaAndPrint(raw)
}

func commandExplore(conf *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must pass one argument")
	}

	url := LocationAreasURL + args[0]
	raw, ok := conf.Cache.Get(url)
	if !ok {
		return exploreLocationArea(conf, url)
	}

	return decodeLocationAreaAndPrint(raw)
}

func decodePokemonAndPrint(raw []byte) (Pokemon, error) {
	var data Pokemon
	if err := json.Unmarshal(raw, &data); err != nil {
		return Pokemon{}, fmt.Errorf("error decoding data: %v", err)
	}

	return data, nil
}

func getPokemon(conf *config, url string) (Pokemon, error) {
	raw, ok := conf.Cache.Get(url)
	if !ok {
		fetchedRaw, err := fetchData(url)
		if err != nil {
			return Pokemon{}, err
		}
		conf.Cache.Add(url, fetchedRaw)
		raw = fetchedRaw
	}

	return decodePokemonAndPrint(raw)
}

func getRandomCatchChance(conf *config, data Pokemon) bool {
	chance := rand.Intn(600) > data.BaseExperience
	if chance {
		conf.Pokedex[data.Name] = data
		return true
	} else {
		return false
	}
}

func commandCatch(conf *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must pass one argument")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])
	url := PokemonURL + args[0]
	data, err := getPokemon(conf, url)
	if err != nil {
		return err
	}
	if getRandomCatchChance(conf, data) {
		fmt.Printf("%s was caught!\n", args[0])
	} else {
		fmt.Printf("%s escaped!\n", args[0])
	}

	return nil
}

func commandInspect(conf *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must pass one argument")
	}
	data, ok := conf.Pokedex[args[0]]
	if !ok {
		return fmt.Errorf("%s is not catched before", args[0])
	}

	fmt.Println("Name: ", data.Name)
	fmt.Println("Height: ", data.Height)
	fmt.Println("Weight: ", data.Weight)
	fmt.Println("Stats:")
	for _, stat := range data.Stats {
		fmt.Printf("\t-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokType := range data.Types {
		fmt.Printf("\t-%s\n", pokType.Type.Name)
	}

	return nil
}

func commandPokedex(conf *config, args ...string) error {
	if len(conf.Pokedex) <= 0 {
		fmt.Println("You didn't caught any Pokemon before")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, value := range conf.Pokedex {
		fmt.Printf("\t- %s\n", value.Name)
	}

	return nil
}
