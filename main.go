package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
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
	callback    func(*config, []string) error
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("help: Displays a help message")
	fmt.Println("map: Displays 20 locations")
	fmt.Println("mapb: Displays last 20 locations")
	fmt.Println("explore <area>: Display Pokemon in area")
	fmt.Println("catch: Catch a Pokemon with a pokeball")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func commandMap(cfg *config, args []string) error {
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

func commandMapBack(cfg *config, args []string) error {
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

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provie a location name")
	}

	areaName := args[0]

	fmt.Printf("Exploring %s", areaName)
	fmt.Printf("Found Pokemon:\n")

	details, err := pokeapi.GetLocationAreaDetails(areaName)
	if err != nil {
		return err
	}

	for _, encounter := range details.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Please provide a pokemon name")
	}

	name := args[0]

	fmt.Printf("Throwing a Pokeball at %s...", name)

	pokemon, err := pokeapi.GetPokemon(name)
	if err != nil {
		return err
	}

	roll := rand.Intn(pokemon.BaseExperience + 1)

	if roll < 50 {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
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
	"explore": {
		name:        "explore",
		description: "Explore a location area area",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Throw a Pokeball to catch a Pokemon",
		callback:    commandCatch,
	},
}

var pokedex = make(map[string]pokeapi.Pokemon)

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
		args := input[1:]

		command, ok := commands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := command.callback(cfg, args); err != nil {
			fmt.Println("Error", err)
		}
	}
	if err := words.Err(); err != nil {
		log.Fatal(err)
	}
}
