package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"net/http"
	"encoding/json"
)
type cliCommand struct{
	name        string
	description string
	callback    func(*config) error
}

var commandMap map[string]cliCommand

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range commandMap {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil 
}

type locationAreaResponse struct {
	Results  []struct{ Name string } `json:"results"`
	Next     string                  `json:"next"`
	Previous string                  `json:"previous"`
} 



func commandMapExplore(cfg *config) error {
	var locationOffset = 0

	url := cfg.next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch locations: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received bad status: %s", resp.Status)
	}

	var data locationAreaResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if len(data.Results) == 0 {
		fmt.Println("No more locations.")
		return nil
	}

	for _, loc := range data.Results {
		fmt.Println(loc.Name)
	}

	locationOffset += 20
	cfg.next = data.Next
	cfg.previous = data.Previous
	return nil
}

func commandMapBack(cfg *config) error {

	if cfg.previous == "" {
		fmt.Println("You're already at the beginning")
		return nil
	}

	resp, err := http.Get(cfg.previous)
	if err != nil {
		return fmt.Errorf("failed to fetch locations: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	var data locationAreaResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode: %v", err)
	}

	for _, loc := range data.Results {
		fmt.Println(loc.Name)
	}

	cfg.next = data.Next
	cfg.previous = data.Previous

	return nil
}

type config struct {
	next     string
	previous string
} 

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	return strings.Fields(text)
}

func main() {
	cfg:= &config{}

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
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("Pokedex > ")
		//fmt.Println("Hello there! Feel free to enter something!")

		if !scanner.Scan() {
			break //exit loop if scanner cant read
		}

		// get input
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		//print first word as command
		//fmt.Println("Your command was:", words[0])
		cmdName := input[0]
		cmd, exists := commandMap[cmdName]
		if !exists {
			fmt.Printf("Unknown command: %s. Type 'help' to see available commands.\n", cmdName)
			continue
		}
		if err := cmd.callback(cfg); err != nil {
			fmt.Println("Error:, err")
		}
	}


}
	