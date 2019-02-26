package cards

import (
//	"fmt"
)

type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

type Rank int

const (
	Two Rank = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

type Scoring int

type Card struct {
	suit Suit
	rank Rank
}

type Deck struct {
	cards []Card
}

func make_deck() Deck {
	var d Deck
	d.cards = make([]Card, 52)
	return d
}
