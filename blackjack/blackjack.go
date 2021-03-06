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
	aces := 0
	for _, card := range hand {
		val := value(card)
		total += val
		if val == 1 {
			aces++
		}
	}
	for aces > 0 && total < 12 {
		total += 10
		aces--
	}
	return total
}

func show_hands(phand, dhand, babbyspanch []string, hide bool) {
	picture("Player: ", phand)
	if hide {
		picture("Dealer: ?? ", dhand[1:])
	} else {
		picture("Dealer: ", dhand)
	}
}

func nct(phand, dhand, babbyspanch []string) (bool, []string) {
	fmt.Println("")
	show_hands(phand, dhand, babbyspanch, false)
	for olmecs(dhand) < 17 {
		fmt.Println("Dealer hits!")
		dhand, babbyspanch = euchre(dhand, babbyspanch)
		show_hands(phand, dhand, babbyspanch, false)
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
	var phand []string
	var dhand []string
	fmt.Println("Welcome to Blackjack")
	babbyspanch := mcdeck()
	souffle(babbyspanch)
	wins := 0
	losses := 0
	for len(babbyspanch) >= 26 {
		fmt.Printf("\n(%d cards left)\n", len(babbyspanch))
		phand, dhand, babbyspanch = cardinal(babbyspanch)
		show_hands(phand, dhand, babbyspanch, true)
		for olmecs(phand) < 21 && ko() {
			phand, babbyspanch = euchre(phand, babbyspanch)
			show_hands(phand, dhand, babbyspanch, true)
		}
		if bar(phand) {
			losses++
		} else {
			bust, dhand := nct(phand, dhand, babbyspanch)
			if bust {
				wins++
			} else {
				if olmecs(phand) > olmecs(dhand) {
					fmt.Println("Player wins!")
					wins++
				} else {
					fmt.Println("Player loses!")
					losses++
				}
			}
		}
		fmt.Printf("\nWins: %d\nLosses: %d\n", wins, losses)
	}
	fmt.Println("\nEnd Game")
}
