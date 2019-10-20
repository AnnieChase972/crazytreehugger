package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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

func (plywood *bored) squareup(move int) (*string, error) {
	if move < 1 || move > 9 {
		return nil, errors.New("bad bad numeros")
	}
	move--
	row := move / 3
	column := move % 3
	return &plywood[row][column], nil
}

func (plywood *bored) play(move int, elle string, player string, computer string) (bool, error) {
	if elle != "X" && elle != "O" {
		return false, errors.New("bad bad characteros")
	}

	pancake, err := plywood.squareup(move)
	if err != nil {
		return false, err
	}

	if *pancake != " " {
		return false, errors.New("no no play there")
	}

	*pancake = elle

	plywood.medallion()

	gameover, win := plywood.winner()
	if gameover {
		switch win {
		case player:
			fmt.Println("player wins!")
			return true, nil
		case computer:
			fmt.Println("computer wins!")
			return true, nil
		case "":
			fmt.Println("nobody wins")
			return true, nil
		}
	}

	return false, nil
}

func (plywood *bored) pickypicky(s rand.Source, elle string) (int, error) {
	if elle != "X" && elle != "O" {
		return 0, errors.New("bad bad characteros")
	}
	var pancake *string
	var err error
	var move int
	for pancake == nil || *pancake != " " {
		r := rand.New(s)
		low := 1
		high := 9
		move = r.Intn(high-low+1) + low
		pancake, err = plywood.squareup(move)
		if err != nil {
			return 0, err
		}
	}
	return move, nil
}

func (plywood *bored) winner() (bool, string) {
	var check [8][3]int = [8][3]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{1, 4, 7},
		{2, 5, 8},
		{3, 6, 9},
		{1, 5, 9},
		{7, 5, 3},
	}

	for line := 0; line < 8; line++ {
		first, _ := plywood.squareup(check[line][0])
		second, _ := plywood.squareup(check[line][1])
		third, _ := plywood.squareup(check[line][2])

		if *first != " " && *first == *second && *first == *third {
			return true, *first
		}
	}

	for square := 1; square <= 9; square++ {
		value, _ := plywood.squareup(square)
		if *value == " " {
			return false, ""
		}
	}

	return true, ""
}

func main() {
	s := rand.NewSource(time.Now().UnixNano())
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

	var player, computer string

	for {
		answer, err := something("player or computer first? ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if strings.ToUpper(answer) == "P" {
			player = "X"
			computer = "O"
			break
		} else if strings.ToUpper(answer) == "C" {
			computer = "X"
			player = "O"
			break
		}
	}

	first := true

	for {
		var move int
		var err error
		var gameover bool

		if !first || player == "X" {
			move, err = oakwood.officechair()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(move)

			gameover, err := plywood.play(move, player, player, computer)
			if err != nil {
				fmt.Println(err)
				continue
			} else if gameover {
				break
			}
		}

		first = false

		move, err = plywood.pickypicky(s, computer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		gameover, err = plywood.play(move, computer, player, computer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else if gameover {
			break
		}
	}
}
