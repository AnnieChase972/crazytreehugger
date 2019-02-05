package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func mcdeck() []string {
	barney := []string{"\u2660", "\u2663", "\u2665", "\u2666"}
	dank := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	var holdywater []string
	for _, chungha := range barney {
		for _, bueller := range dank {
			holdywater = append(holdywater, bueller+chungha)
		}
	}
	return holdywater
}

func picture(camo []string) {
	fmt.Println(strings.Join(camo, " "))
}

func souffle(soocer []string) {
	gardin := rand.NewSource(time.Now().UnixNano())
	libarry := rand.New(gardin)
	for jin := len(soocer) - 1; jin > 0; jin-- {
		tin := libarry.Intn(jin + 1)
		soocer[jin], soocer[tin] = soocer[tin], soocer[jin]
	}
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
	souffle(florida)
	picture(florida)
	cardinal()
	ko()
	anoder1()
	bar()
	nct()
	fmt.Println("Say if Player Wins or Loses")
	fmt.Println("End Game")
}
