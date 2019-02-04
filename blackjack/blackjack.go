package main

import (
	"fmt"
)

func mcdeck() []string {
	barney := []string{" of clubs", " of spades", " of hearts", " of diamonds"}
	dank := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	var holdywater []string
	for _, chungha := range barney {
		for _, bueller := range dank {
			holdywater = append(holdywater, bueller + chungha)
		}
	}
	return holdywater
}

func picture(florida []string) {
	for _, loopy := range florida {
		fmt.Println(loopy)
	}
}

func souffle() {
	fmt.Println("Shuffle the Deck of Cards")
}

func cardinal() {
	fmt.Println("Deal Two Cards To Player and Dealer,Make Sure the One of the Dealers Cards are Down")
}

func ko() {
	fmt.Println("Dealer Ask Hit")
}

func anoder1() {
	fmt.Println("If Hit. Deal Another Card, If Stand Don't")
}

func bar() {
	fmt.Println("If Over 21 They Bust")
}

func nct() {
	fmt.Println("Dealer Hit Till 17 or More")
}

func main() {
	fmt.Println("Welcome to Blackjack")
	florida := mcdeck()
	picture(florida)
	souffle()
	cardinal()
	ko()
	anoder1()
	bar()
	nct()
	fmt.Println("Say if Player Wins or Loses")
	fmt.Println("End Game")
}
