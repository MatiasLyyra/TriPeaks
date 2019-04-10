package game

import (
	"fmt"

	"github.com/MatiasLyyra/TriPeaks/deck"
)

type PeakCard struct {
	deck.Card
	Removed   bool
	ChildLeft int
}

func (pk *PeakCard) SubChild() {
	if pk.ChildLeft > 0 {
		pk.ChildLeft--
		if pk.ChildLeft == 0 {
			pk.Card.FaceDown = false
		}
	}
}

func (c PeakCard) String() string {
	if c.Removed {
		return "      "
	}
	return c.Card.String()
}

type TriPeaksDeck [28]PeakCard

type TriPeaks struct {
	Stock     deck.Deck
	Discards  []deck.Card
	Cards     TriPeaksDeck
	CardsLeft int
	Score     int
	Streak    int
}

func NewTripeaks(stock deck.Deck) *TriPeaks {
	if stock.Len() != 52 {
		panic("deck requires 52 cards")
	}
	cardsLeft := 0
	_, discard := stock.Pop()
	discard.FaceDown = false
	game := TriPeaks{
		Stock:    stock,
		Discards: []deck.Card{discard},
	}
	for i := 0; i < len(game.Cards); i++ {
		_, card := game.Stock.Pop()

		cardsLeft++
		if i >= 18 {
			game.Cards[i] = PeakCard{
				Card:      card,
				Removed:   false,
				ChildLeft: 0,
			}
		} else {
			card.FaceDown = true
			game.Cards[i] = PeakCard{
				Card:      card,
				Removed:   false,
				ChildLeft: 2,
			}
		}

	}
	game.CardsLeft = cardsLeft
	return &game
}

func (tri *TriPeaks) GameOver() bool {
	legalMoves, _ := tri.LegalMoves()
	return len(legalMoves) == 0
}

func (tri *TriPeaks) Copy() *TriPeaks {
	newTri := &TriPeaks{
		Discards:  make([]deck.Card, len(tri.Discards)),
		CardsLeft: tri.CardsLeft,
		Score:     tri.Score,
		Streak:    tri.Streak,
	}
	deck.Copy(newTri.Discards, tri.Discards)
	newTri.Stock = *tri.Stock.Copy()
	for i, card := range tri.Cards {
		newTri.Cards[i] = card
	}
	return newTri
}
func (tri *TriPeaks) String() string {
	var gameState string
	for i := 0; i < 3; i++ {
		if i == 0 {
			gameState += fmt.Sprintf("         %s", tri.Cards[i])
		} else {
			gameState += fmt.Sprintf("            %s", tri.Cards[i])
		}
	}
	gameState += "\n"
	for i := 3; i < 9; i++ {
		if (i-3)%2 == 0 {
			gameState += fmt.Sprintf("      %s", tri.Cards[i])
		} else {
			gameState += fmt.Sprintf("%s", tri.Cards[i])

		}
	}
	gameState += "\n   "
	for i := 9; i < 18; i++ {
		gameState += fmt.Sprintf("%s", tri.Cards[i])
	}
	gameState += "\n"
	for i := 18; i < 28; i++ {
		gameState += fmt.Sprintf("%s", tri.Cards[i])
	}
	gameState += "\n"
	return gameState
}
func (tri *TriPeaks) Surrender() {
	for i, card := range tri.Cards {
		if !card.Removed {
			tri.Score -= 5
		}
		tri.Cards[i].Removed = true
	}
	tri.CardsLeft = 0
}
func (tri *TriPeaks) Select(pos int) bool {
	if pos < 0 || pos >= len(tri.Cards) {
		return false
	}
	card := &tri.Cards[pos]
	if !tri.IsLegal(*card) {
		return false
	}
	card.Removed = true
	tri.AddDiscard(card.Card)
	tri.ApplyReveals(pos)
	tri.CardsLeft--
	tri.incrementScore()
	if pos < 3 {
		tri.Score += 15
	}
	if tri.Cards[0].Removed && tri.Cards[1].Removed && tri.Cards[2].Removed {
		tri.Score += 15
	}
	return true
}

func (tri *TriPeaks) incrementScore() {
	tri.Streak += 1
	tri.Score += tri.Streak
}

func (tri *TriPeaks) Discard() deck.Card {
	return tri.Discards[0]
}

func (tri *TriPeaks) AddDiscard(card deck.Card) {
	temp := tri.Discards[0]
	tri.Discards = append(tri.Discards, temp)
	tri.Discards[0] = card
}

func (tri *TriPeaks) LegalMoves() ([]int, bool) {
	legalMoves := make([]int, 0)
	for pos, card := range tri.Cards {
		if tri.IsLegal(card) {
			legalMoves = append(legalMoves, pos)
		}
	}
	canDraw := tri.Stock.Len() > 0
	if canDraw {
		legalMoves = append(legalMoves, -1)
	}
	return legalMoves, canDraw
}

func (tri *TriPeaks) CheckReveals(pos int) (int, int) {
	leftPos := -1
	rightPos := -1
	if pos >= 3 && pos < 9 {
		switch pos {
		case 3, 4:
			leftPos = 0
		case 5, 6:
			leftPos = 1
		case 7, 8:
			leftPos = 2
		}
	} else if pos >= 9 && pos < 18 {
		switch pos {
		case 9:
			rightPos = 3
		case 10:
			leftPos = 3
			rightPos = 4
		case 11:
			leftPos = 4
		case 12:
			rightPos = 5
		case 13:
			leftPos = 5
			rightPos = 6
		case 14:
			leftPos = 6
		case 15:
			rightPos = 7
		case 16:
			leftPos = 7
			rightPos = 8
		case 17:
			leftPos = 8
		}
	} else if pos >= 18 && pos < 28 {
		switch pos {
		case 18, 19, 20, 21:
			rightPos = pos - 18 + 9
			if pos != 18 {
				leftPos = rightPos - 1
			}
		case 22, 23, 24:
			rightPos = pos - 22 + 13
			leftPos = rightPos - 1
		case 25, 26, 27:
			rightPos = pos - 25 + 16
			leftPos = rightPos - 1
			if pos == 27 {
				rightPos = -1
			}
		}
	}

	return leftPos, rightPos
}

func (tri *TriPeaks) ApplyReveals(pos int) {
	leftPos, rightPos := tri.CheckReveals(pos)
	if leftPos != -1 {
		tri.Cards[leftPos].SubChild()
	}
	if rightPos != -1 {
		tri.Cards[rightPos].SubChild()
	}
}

func (tri *TriPeaks) UsedCards() []deck.Card {
	cards := make([]deck.Card, 0, 10)
	for _, card := range tri.Cards {
		if !card.FaceDown {
			cards = append(cards, card.Card)
		}
	}
	for _, card := range tri.Discards {
		cards = append(cards, card)
	}
	return cards
}

func (tri *TriPeaks) IsLegal(card PeakCard) bool {
	return !card.FaceDown &&
		card.ChildLeft == 0 &&
		!card.Removed &&
		(card.Rank-1 == tri.Discard().Rank ||
			card.Rank+1 == tri.Discard().Rank ||
			(card.Rank == 2 && tri.Discard().Rank == 14) ||
			(card.Rank == 14 && tri.Discard().Rank == 2))
}

func (tri *TriPeaks) Draw() bool {
	ok, card := tri.Stock.Pop()
	if ok {
		tri.Score -= 5
		tri.Streak = 0
		tri.AddDiscard(card)
	}
	return ok
}
