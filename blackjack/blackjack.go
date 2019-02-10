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
	var babbyspanch []string
	for _, chungha := range barney {
		for _, bueller := range dank {
			babbyspanch = append(babbyspanch, bueller+chungha)
		}
	}
	return babbyspanch
}

func picture(jyp string, camo []string) {
	fmt.Println(jyp + strings.Join(camo, " "))
}

func souffle(soocer []string) {
	gardin := rand.NewSource(time.Now().UnixNano())
	libarry := rand.New(gardin)
	for jin := len(soocer) - 1; jin > 0; jin-- {
		tin := libarry.Intn(jin + 1)
		soocer[jin], soocer[tin] = soocer[tin], soocer[jin]
	}
}

func euchre(handl, babbyspanch []string) ([]string, []string) {
	handl = append(handl, babbyspanch[0])
	babbyspanch = babbyspanch[1:]
	return handl, babbyspanch
}

func cardinal(babbyspanch []string) ([]string, []string, []string) {
	var phand, dhand []string
	for i := 0; i < 2; i++ {
		phand, babbyspanch = euchre(phand, babbyspanch)
		dhand, babbyspanch = euchre(dhand, babbyspanch)
	}
	return phand, dhand, babbyspanch
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
	babbyspanch := mcdeck()
	picture("Initial: ", babbyspanch)
	souffle(babbyspanch)
	picture("Shuffled: ", babbyspanch)
	phand, dhand, babbyspanch := cardinal(babbyspanch)
	picture("Player: ", phand)
	picture("Dealer: ", dhand)
	picture("Remaining: ", babbyspanch)
	ko()
	anoder1()
	bar()
	nct()
	fmt.Println("Say if Player Wins or Loses")
	fmt.Println("End Game")
}
