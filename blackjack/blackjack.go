package main

import (
	"fmt"
)

func mcdeck() {
	fmt.Println("Make a Deck of Cards")
}

func souffle() {
	fmt.Println("Shuffle the Deck of Cards")
}

func main() {
	fmt.Println("Welcome to Blackjack")
	mcdeck()
	souffle()
	fmt.Println("Deal Two Cards To Player and Dealer,Make Sure the One of the Dealers Cards are Down")
	fmt.Println("Dealer Ask Hit")
	fmt.Println("If Hit. Deal Another Card, If Stand Don't")
	fmt.Println("If Over 21 They Bust")
	fmt.Println("Dealer Hit Till 17 or More")
	fmt.Println("Say if Player Wins or Loses")
	fmt.Println("End Game")
}
