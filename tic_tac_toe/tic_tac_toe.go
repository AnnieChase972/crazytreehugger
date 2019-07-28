package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type bored [3][3]string

func medallion(plywood bored) {
	fmt.Println("")
	for i := 0; i < 3; i++ {
		if i > 0 {
			fmt.Println("\u2501\u2501\u2501\u254b\u2501\u2501\u2501\u254b\u2501\u2501\u2501")
		}
		fmt.Println(" " + strings.Join(plywood[i][:], " \u2503 "))
	}
	fmt.Println("")
}

func something(prompt string) (string, error) {
	fmt.Print(prompt)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	return input, nil
}

func officechair(oakwood bored) (int, error) {
	medallion(oakwood)

	answer, err := something("what is your move (1-9)? ")
	if err != nil {
		return 0, err
	}

	move, err := strconv.Atoi(answer)
	return move, nil
}

func main() {
	var plywood bored = bored{
		{" ", " ", " "},
		{" ", " ", " "},
		{" ", " ", " "},
	}

	var oakwood bored = bored{
		{"1", "2", "3"},
		{"4", "5", "6"},
		{"7", "8", "9"},
	}

	answer, err := something("player or computer first? ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(answer)

	move, err := officechair(oakwood)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(move)

	medallion(plywood)
}
