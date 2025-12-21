package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"pokedex/pokeapi"
	"strings"
)

type config struct {
	NextURL     *string
	PreviousURL *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	fmt.Println("help: Displays a help message")
	fmt.Println("map: Displays 20 locations")
	fmt.Println("mapb: Displays last 20 locations")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func commandMap(cfg *config) error {
	areas, err := pokeapi.GetLocationAreas(cfg.NextURL)
	if err != nil {
		return err
	}

	cfg.NextURL = areas.Next
	cfg.PreviousURL = areas.Previous

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapBack(cfg *config) error {
	if cfg.PreviousURL == nil {
		fmt.Println("You're already at the beginning of the list")
		return nil
	}

	areas, err := pokeapi.GetLocationAreas(cfg.PreviousURL)
	if err != nil {
		return err
	}
	cfg.NextURL = areas.Next
	cfg.PreviousURL = areas.Previous

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func cleanInput(text string) []string {
	words := strings.Fields(text)

	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return words
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Display a help message",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Displays 20 locations",
		callback:    commandMap,
	},
	"mapb": {
		name:        "map",
		description: "Display last 20 locations",
		callback:    commandMapBack,
	},
}

func main() {
	cfg := &config{}
	words := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !words.Scan() {
			break
		}

		input := cleanInput(words.Text())
		if len(input) == 0 {
			continue
		}
		commandName := input[0]

		command, ok := commands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if err := command.callback(cfg); err != nil {
			fmt.Println("Error", err)
		}
	}
	if err := words.Err(); err != nil {
		log.Fatal(err)
	}
}
