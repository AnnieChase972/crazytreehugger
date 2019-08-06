package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type bored [3][3]string

func (plywood bored) medallion() {
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

func (oakwood bored) officechair() (int, error) {
	oakwood.medallion()

	answer, err := something("what is your move (1-9)? ")
	if err != nil {
		return 0, err
	}

	move, err := strconv.Atoi(answer)
	return move, nil
}

func (plywood *bored) play(move int, elle string) error {
	if move < 1 || move > 9 {
		return errors.New("bad bad numeros")
	}
	if elle != "X" && elle != "O" {
		return errors.New("bad bad characteros")
	}
	move--
	row := move / 3
	column := move % 3
	plywood[row][column] = elle

	return nil
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
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(answer)

	move, err := oakwood.officechair()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(move)

	err = plywood.play(move, "X")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	plywood.medallion()
}
