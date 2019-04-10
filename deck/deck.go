package deck

import (
	"math/rand"
	"time"
)

// Deck represents normal playing card deck
type Deck struct {
	// TODO: Define this as type Deck []Card instead
	Cards []Card
}

// New creates 52 card deck with 13 cards of each suit
func New() *Deck {
	deck := Deck{
		Cards: make([]Card, 0, 52),
	}
	for i := 2; i <= 14; i++ {
		deck.Cards = append(deck.Cards, Card{
			Rank:     i,
			Suit:     Hearts,
			FaceDown: false,
		})
		deck.Cards = append(deck.Cards, Card{
			Rank:     i,
			Suit:     Spades,
			FaceDown: false,
		})
		deck.Cards = append(deck.Cards, Card{
			Rank:     i,
			Suit:     Diamonds,
			FaceDown: false,
		})
		deck.Cards = append(deck.Cards, Card{
			Rank:     i,
			Suit:     Clubs,
			FaceDown: false,
		})
	}
	return &deck
}

// Creates a copy of the deck
func (d *Deck) Copy() *Deck {
	newDeck := &Deck{
		Cards: make([]Card, len(d.Cards)),
	}
	for i, card := range d.Cards {
		newDeck.Cards[i] = card
	}
	return newDeck
}

// Shuffles the deck using rand.Shuffle
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		temp := d.Cards[i]
		d.Cards[i] = d.Cards[j]
		d.Cards[j] = temp
	})
}

func (d *Deck) Len() int {
	return len(d.Cards)
}

func (d *Deck) Pop() (bool, Card) {
	var card Card
	if d.Len() <= 0 {
		return false, card
	}
	card = d.Cards[d.Len()-1]
	d.Cards = append(d.Cards[:d.Len()-1])
	return true, card
}

func Copy(to, from []Card) {
	for i, card := range from {
		to[i] = card
	}
}

func Remove(cards []Card, i int) []Card {
	cards[i] = cards[len(cards)-1]
	cards = cards[:len(cards)-1]
	return cards
}

func RemoveVal(cards []Card, card Card) []Card {
	for i, val := range cards {
		if val.Suit == card.Suit && val.Rank == card.Rank {
			cards = Remove(cards, i)
			return cards
		}
	}
	return cards
}
