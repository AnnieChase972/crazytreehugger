package main

import (
	"fmt"
	"math/rand"
	"strconv"
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
	fmt.Printf("%s%s (%d)\n", jyp, strings.Join(camo, " "), olmecs(camo))
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

func ko() bool {
	var input string
	for {
		fmt.Print("Hit? ")
		fmt.Scanln(&input)
		switch {
		case len(input) > 0 && strings.EqualFold(input[:1], "y"):
			return true
		case len(input) > 0 && strings.EqualFold(input[:1], "n"):
			return false
		}
	}

}

func bar(hand []string) bool {
	if olmecs(hand) > 21 {
		fmt.Println("Player busts!")
		return true
	} else {
		return false
	}
}

func value(card string) int {
	switch card[0] {
	case 'A':
		return 1
	case 'J', 'Q', 'K', '1':
		return 10
	default:
		value, err := strconv.Atoi(card[:1])
		if err != nil {
			fmt.Println(err)
			return 0
		}
		return value
	}
}

func olmecs(hand []string) int {
	total := 0
	for _, card := range hand {
		total += value(card)
	}
	return total
}

func show_hands(phand, dhand, babbyspanch []string) {
	picture("Player: ", phand)
	picture("Dealer: ", dhand)
	picture("Remaining: ", babbyspanch)
}

func nct(phand, dhand, babbyspanch []string) (bool, []string) {
	for olmecs(dhand) < 17 {
		fmt.Println("Dealer hits!")
		dhand, babbyspanch = euchre(dhand, babbyspanch)
		show_hands(phand, dhand, babbyspanch)
	}
	if olmecs(dhand) > 21 {
		fmt.Println("Dealer busts!")
		return true, dhand
	} else {
		fmt.Println("Dealer stands.")
		return false, dhand
	}
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
	for olmecs(phand) < 21 && ko() {
		phand, babbyspanch = euchre(phand, babbyspanch)
		show_hands(phand, dhand, babbyspanch)
	}
	if !bar(phand) {
		bust, dhand := nct(phand, dhand, babbyspanch)
		if !bust {
			if olmecs(phand) > olmecs(dhand) {
				fmt.Println("Player wins!")
			} else {
				fmt.Println("Player loses!")
			}
		}
	}
	fmt.Println("End Game")
}
