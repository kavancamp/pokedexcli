package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"io"
	"math/rand"
	"github.com/kavancamp/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args []string) error
}

var commandMap map[string]cliCommand

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range commandMap {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

type locationAreaResponse struct {
	Results  []struct{ Name string } `json:"results"`
	Next     string                  `json:"next"`
	Previous string                  `json:"previous"`
}

func commandMapExplore(cfg *config, args []string) error {
	url := cfg.next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	}

	data, err := fetchWithCache(url, cfg.cache)
	if err != nil {
		return err
	}

	var result locationAreaResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err)
	}

	for _, loc := range result.Results {
		fmt.Println(loc.Name)
	}

	cfg.next = result.Next
	cfg.previous = result.Previous
	return nil
}

func commandMapBack(cfg *config, args []string) error {
	if cfg.previous == "" {
		fmt.Println("You're already at the beginning")
		return nil
	}

	data, err := fetchWithCache(cfg.previous, cfg.cache)
	if err != nil {
		return err
	}

	var result locationAreaResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err)
	}

	for _, loc := range result.Results {
		fmt.Println(loc.Name)
	}

	cfg.next = result.Next
	cfg.previous = result.Previous
	return nil
}

func fetchWithCache(url string, cache *pokecache.Cache) ([]byte, error) {
	if val, ok := cache.Get(url); ok {
		fmt.Println("🔁 Cache hit:", url)
		return val, nil
	}

	//fmt.Println("🌐 Fetching from API:", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	cache.Add(url, data)
	return data, nil
}

type config struct {
	next     string
	previous string
	cache    *pokecache.Cache
	pokedex map[string]pokemonEntry
}
type pokemonEntry struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type locationExploreResponse struct{
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		}	`json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func commandMapPokemon(cfg *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: explore <location-area>")
	}
	area := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", area)
	
	data, err := fetchWithCache(url, cfg.cache)
	if err != nil {
		return fmt.Errorf("failed to fetch location-area data: %v", err)
	}

	var parsed locationExploreResponse
	if err := json.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if len(parsed.PokemonEncounters) == 0 {
		fmt.Println("No Pokémon have been found in this area.")
		return nil
	}
	//fmt.Println("Exploring pastoria-city-area...")
	fmt.Printf("Found %d Pokémon in %s:\n", len(parsed.PokemonEncounters), area)
	for _, encounter := range parsed.PokemonEncounters {
		fmt.Println(" -", encounter.Pokemon.Name)		
	}
	return nil
}
func commandCatch(cfg *config, args []string) error{
	if len(args) < 1 {
		return fmt.Errorf("usage: catch <pokemon-name>")
	}
	
	selected_pokemon := strings.ToLower(args[0])
	fmt.Printf("Throwing a Pokeball at %s...\n", selected_pokemon)

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", selected_pokemon)
	data, err := fetchWithCache(url, cfg.cache)
	if err != nil{
		return fmt.Errorf("failed to fetch Pokémon data: %v", err)
	}
	var p pokemonEntry
	if err := json.Unmarshal(data, &p); err != nil {
		return fmt.Errorf("failed to decode Pokémon data: %v", err)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	chance := r.Intn(p.BaseExperience + 1)
	if chance > 50 {
		fmt.Printf("%s escaped!\n", p.Name)
		return nil
	}
	cfg.pokedex[p.Name] = p
	fmt.Printf("%s was caught!\n", p.Name)
	return nil
}
func commandInspect(cfg *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: inspect <pokemon>")
	}
	name := strings.ToLower(args[0])

	p, ok := cfg.pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)
	fmt.Println("Stats:")
	for _, stat := range p.Stats {
		fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf(" -%s\n", t.Type.Name)
	}
	return nil
}
func commandPokedex(cfg *config, args []string) error {
	if len(cfg.pokedex) == 0 {
		fmt.Println("You haven’t caught any Pokémon yet!")
		return nil
	}

	fmt.Println("Your Pokémon:")
	for name := range cfg.pokedex {
		fmt.Printf("  - %s\n", name)
	}
	return nil
}
func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	return strings.Fields(text)
}

func main() {
	cfg := &config{
		cache:   pokecache.NewCache(5 * time.Second),
		pokedex: make(map[string]pokemonEntry),
	}

	commandMap = map[string]cliCommand{
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
			description: "Explore the Pokémon world by listing location areas",
			callback:    commandMapExplore,
		},
		"mapb": {
			name:        "mapb",
			description: "Go back to the previous page of location areas",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore Pokémon at given location",
			callback:    commandMapPokemon,
		},
		"catch": {
			name:        "catch",
			description: "Throw a Pokéball and try to catch a Pokémon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Show details about a caught Pokémon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all Pokémon you've caught",
			callback:    commandPokedex,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokédex > ")
		if !scanner.Scan() {
			break
		}

		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		cmdName := input[0]
		cmdArgs := input[1:]
		cmd, exists := commandMap[cmdName]
		if !exists {
			fmt.Printf("Unknown command: %s. Type 'help' to see available commands.\n", cmdName)
			continue
		}

		if err := cmd.callback(cfg, cmdArgs); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
