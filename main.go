package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		prompt := cleanInput(scanner.Text())
		for key, value := range commands {
			if len(prompt) == 0 {
				continue
			}
			if prompt[0] == key {
				err := value.callback(&conf, prompt[1:]...)
				if err != nil {
					fmt.Printf("Error has occurred: %v\n", err.Error())
				}
			}
		}

		if !isCommand(prompt) && len(prompt) != 0 {
			fmt.Println("Unknown command")
		}

	}
}
