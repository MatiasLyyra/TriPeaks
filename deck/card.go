package deck

import (
	"fmt"
	"strconv"
)

const (
	Hearts = iota
	Spades
	Clubs
	Diamonds
)

type Card struct {
	Rank     int
	Suit     int
	FaceDown bool
}

func (c Card) String() string {
	var suit string
	var rank string
	if c.Rank < 10 {
		rank = strconv.Itoa(c.Rank)
	} else {
		switch c.Rank {
		case 10:
			rank = "T"
		case 11:
			rank = "J"
		case 12:
			rank = "Q"
		case 13:
			rank = "K"
		case 14:
			rank = "A"
		default:
			rank = "?"
		}
	}
	switch c.Suit {
	case Hearts:
		suit = "♥"
	case Spades:
		suit = "♠"
	case Diamonds:
		suit = "♦"
	case Clubs:
		suit = "♣"
	}
	if c.FaceDown {
		return fmt.Sprintf("[    ]")
	}
	return fmt.Sprintf("[%s  %s]", string(rank), string(suit))

}

func (c Card) HashCode() int {
	return c.Suit*100 + c.Rank
}
