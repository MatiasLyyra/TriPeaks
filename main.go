package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
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

type Deck struct {
	Cards []Card
}

func (d *Deck) Suffle() {
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

func NewDeck() *Deck {
	deck := Deck{
		Cards: make([]Card, 0, 52),
	}
	for i := 1; i <= 13; i++ {
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

type PeakCard struct {
	Card
	Removed   bool
	ChildLeft int
}

func (c PeakCard) String() string {
	if c.Removed {
		return "      "
	}
	return c.Card.String()
}

type TriPeaks struct {
	Deck      Deck
	Discard   Card
	peaks     [3][3][]PeakCard
	CardsLeft int
	flat      [10]PeakCard
}

func NewTripeaks(deck Deck) *TriPeaks {
	if deck.Len() != 52 {
		panic("deck requires 52 cards")
	}
	cardsLeft := 0
	_, discard := deck.Pop()
	discard.FaceDown = false
	game := TriPeaks{
		Deck:    deck,
		Discard: discard,
	}
	for i := 0; i < 3; i++ {
		for level := 0; level < 3; level++ {
			cards := make([]PeakCard, 0, 1+level)
			for j := 0; j < 1+level; j++ {
				_, card := game.Deck.Pop()
				card.FaceDown = true
				cardsLeft++
				cards = append(cards, PeakCard{
					Card:      card,
					Removed:   false,
					ChildLeft: 2,
				})
			}
			game.peaks[i][level] = cards
		}
	}
	var flat [10]PeakCard
	for i := 0; i < len(flat); i++ {
		_, card := game.Deck.Pop()
		cardsLeft++
		card.FaceDown = false
		flat[i] = PeakCard{Card: card, Removed: false}
	}
	game.flat = flat
	game.CardsLeft = cardsLeft

	return &game
}

func (tri *TriPeaks) String() string {
	var gameState string
	for level := 0; level < 3; level++ {
		for i := 0; i < 3; i++ {
			// TODO: Clean this mess
			switch level {
			case 0:
				if i == 0 {
					gameState += "         "
				} else {
					gameState += "            "
				}
			case 1:
				gameState += "      "
			case 2:
				if i == 0 {
					gameState += "   "
				}
			}
			for j := 0; j < len(tri.peaks[i][level]); j++ {
				card := tri.peaks[i][level][j]
				gameState += fmt.Sprintf("%s", card)
			}
		}
		gameState += "\n"
	}
	for _, card := range tri.flat {
		gameState += fmt.Sprintf("%s", card)
	}
	return gameState
}

type TriPeakMove struct {
	Level int
	Pos   int
}

func (tri *TriPeaks) LegalMoves() (int, []TriPeakMove, bool) {
	legalMoves := make([]TriPeakMove, 0)
	n := 0
	for i := 0; i < 3; i++ {
		for level := 0; level < 3; level++ {
			for j := 0; j < len(tri.peaks[i][level]); j++ {
				card := tri.peaks[i][level][j]
				if tri.IsLegal(card) {
					legalMoves = append(legalMoves, TriPeakMove{Level: level, Pos: j})
					n++
				}
			}
		}
	}
	for i := 0; i < len(tri.flat); i++ {
		card := tri.flat[i]
		if tri.IsLegal(card) {
			n++
			legalMoves = append(legalMoves, TriPeakMove{Level: 3, Pos: i})
		}
	}
	canDrawFromDeck := false
	if tri.Deck.Len() > 0 {
		n++
		canDrawFromDeck = true
	}
	return n, legalMoves, canDrawFromDeck
}

func (tri *TriPeaks) IsLegal(card PeakCard) bool {
	return !card.FaceDown &&
		!card.Removed &&
		(card.Rank-1 == tri.Discard.Rank ||
			card.Rank+1 == tri.Discard.Rank)
}

func (tri *TriPeaks) Draw() bool {
	ok, card := tri.Deck.Pop()
	if ok {
		tri.Discard = card
	}
	return ok
}

func (tri *TriPeaks) MoveToDiscard(level, pos int) {
	var card *PeakCard
	if level < 0 || level > 3 {
		return
	}
	if pos < 0 {
		return
	}
	if level == 3 {
		if pos > 9 {
			return
		}
		card = &tri.flat[pos]
	} else {
		if pos >= 3*(level+1) {
			return
		}
		sec := pos / (level + 1)
		secPos := pos % (level + 1)
		card = &tri.peaks[sec][level][secPos]
	}
	if tri.IsLegal(*card) {
		card.Removed = true
		tri.CardsLeft--
		tri.Discard = card.Card
		tri.checkReveals(TriPeakMove{Level: level, Pos: pos})
	}
}

func (tri *TriPeaks) checkReveals(coord TriPeakMove) {
	rightCoord := TriPeakMove{Level: coord.Level - 1, Pos: coord.Pos}
	leftCoord := TriPeakMove{Level: coord.Level - 1, Pos: coord.Pos - 1}
	if rightCoord.Level < 0 {
		return
	}
	if leftCoord.Pos > 0 {
		sec := leftCoord.Pos / (leftCoord.Level + 1)
		secPos := leftCoord.Pos % (leftCoord.Level + 1)
		card := &tri.peaks[sec][leftCoord.Level][secPos]
		card.ChildLeft--
		if card.ChildLeft <= 0 {
			card.FaceDown = false
		}
	}
	if rightCoord.Pos < 3*(rightCoord.Level+1) {
		sec := rightCoord.Pos / (rightCoord.Level + 1)
		secPos := rightCoord.Pos % (rightCoord.Level + 1)
		card := &tri.peaks[sec][rightCoord.Level][secPos]
		card.ChildLeft--
		if card.ChildLeft <= 0 {
			card.FaceDown = false
		}
	}
}

type triNode struct {
	removed bool
	value   Card
	lhs     *Card
	rhs     *Card
	parent  *Card
}

func main() {
	deck := NewDeck()
	deck.Suffle()
	game := NewTripeaks(*deck)
	for game.CardsLeft != 0 {
		fmt.Printf("%s\n", game)
		fmt.Printf("Cards in the deck %d\tCards in peak: %d\tDiscard %s\n", game.Deck.Len(), game.CardsLeft, game.Discard)
		moveCnt, _, _ := game.LegalMoves()
		if moveCnt == 0 {
			fmt.Println("Game Over")
			break
		}
		var (
			level int
			pos   int
		)
		fmt.Println("Give level and pos '<level> <pos>':")
		n, err := fmt.Scanf("%d %d", &level, &pos)
		if err != nil || n != 2 {
			fmt.Println("error with reading input")
			continue
		}
		if level == -1 {
			game.Draw()
		} else {
			game.MoveToDiscard(level-1, pos-1)
		}
	}
}
