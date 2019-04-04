package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
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

func (c Card) HashCode() int {
	return c.Suit*100 + c.Rank
}

type Deck struct {
	Cards []Card
}

func (d *Deck) Copy() *Deck {
	newDeck := &Deck{
		Cards: make([]Card, len(d.Cards)),
	}
	for i, card := range d.Cards {
		newDeck.Cards[i] = card
	}
	return newDeck
}
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

func NewDeck() *Deck {
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

type PeakCard struct {
	Card
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
	Deck      Deck
	Discards  []Card
	Cards     TriPeaksDeck
	CardsLeft int
}

func NewTripeaks(deck Deck) *TriPeaks {
	if deck.Len() != 52 {
		panic("deck requires 52 cards")
	}
	cardsLeft := 0
	_, discard := deck.Pop()
	discard.FaceDown = false
	game := TriPeaks{
		Deck:     deck,
		Discards: []Card{discard},
	}
	for i := 0; i < len(game.Cards); i++ {
		_, card := game.Deck.Pop()
		card.FaceDown = true
		cardsLeft++
		if i >= 18 {
			card.FaceDown = false
		}
		game.Cards[i] = PeakCard{
			Card:      card,
			Removed:   false,
			ChildLeft: 2,
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
		Discards:  make([]Card, len(tri.Discards)),
		CardsLeft: tri.CardsLeft,
	}
	CopyCards(newTri.Discards, tri.Discards)
	newTri.Deck = *tri.Deck.Copy()
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
	return true
}

func (tri *TriPeaks) Discard() Card {
	return tri.Discards[0]
}

func (tri *TriPeaks) AddDiscard(card Card) {
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
	canDraw := tri.Deck.Len() > 0
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

func (tri *TriPeaks) UsedCards() []Card {
	cards := make([]Card, 0, 10)
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
		!card.Removed &&
		(card.Rank-1 == tri.Discard().Rank ||
			card.Rank+1 == tri.Discard().Rank)
}

func (tri *TriPeaks) Draw() bool {
	ok, card := tri.Deck.Pop()
	if ok {
		tri.AddDiscard(card)
	}
	return ok
}

type Deter struct {
	Pos         int
	Card        Card
	Initialized bool
}

type NodeData struct {
	CardsLeft []Card
}

type Node struct {
	X        float64
	N        int
	Pos      int
	leftDet  Deter
	rightDet Deter
	Parent   *Node
	Children []*Node
	Data     *NodeData
}

func (n *Node) SetLeft(det Deter) {
	if n.leftDet.Initialized {
		panic("Overwriting old determinization")
	} else {
		n.leftDet = det
	}
}
func (n *Node) SetRight(det Deter) {
	if n.rightDet.Initialized {
		panic("Overwriting old determinization")
	} else {
		n.rightDet = det
	}
}

func (n *Node) GetUnvisitedChild() *Node {
	unvisited := make([]*Node, 0)
	for _, child := range n.Children {
		if child.N == 0 {
			unvisited = append(unvisited, child)
		}
	}
	if len(unvisited) == 0 {
		return nil
	}
	return unvisited[0]
}

func (n *Node) ChildPos(pos int) int {
	for i, child := range n.Children {
		if child.Pos == pos {
			return i
		}
	}
	return -1
}

func (n *Node) GetParentDeterminization(pos int, left bool) bool {
	if n.Parent == nil {
		return false
	}
	assignL := func(p *Node) {
		if left {
			n.SetLeft(p.leftDet)
		} else {
			n.SetRight(p.leftDet)
		}
	}
	assignR := func(p *Node) {
		if left {
			n.SetLeft(p.rightDet)
		} else {
			n.SetRight(p.rightDet)
		}
	}
	for parent := n.Parent; parent != nil; parent = parent.Parent {
		if parent.leftDet.Pos == pos {
			assignL(parent)
			return true
		} else if parent.rightDet.Pos == pos {
			assignR(parent)
			return true
		}
	}
	return false
}

func NewNode() *Node {
	return &Node{
		X:        0,
		N:        0,
		Pos:      -2,
		Parent:   nil,
		Children: make([]*Node, 0, 5),
	}
}

type SearchResult struct {
	Move  int
	Score float64
}

func MctsSearch(game *TriPeaks, determinizations, trajectories int) []SearchResult {
	random := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	unusedCards := UnusedCards(game)
	rootRewards := make(map[int]float64)
	gameCopy := &TriPeaks{}
	var (
		root            *Node
		unusedCardsCopy []Card
	)
	for i := 0; i < determinizations; i++ {
		root = NewNode()

		for j := 0; j < trajectories; j++ {
			unusedCardsCopy = make([]Card, len(unusedCards))
			CopyCards(unusedCardsCopy, unusedCards)
			data := &NodeData{
				CardsLeft: unusedCardsCopy,
			}
			root.Data = data
			gameCopy = game.Copy()
			var node *Node
			node = MctsSelect(gameCopy, root)
			if !gameCopy.GameOver() {
				node = DeterminizeState(node, gameCopy, random)
			}
			reward := MctsSimulation(gameCopy, node, random)
			MctsBackpropagate(node, reward)
		}
		for _, child := range root.Children {
			if _, exists := rootRewards[child.Pos]; exists {
				rootRewards[child.Pos] += child.X
			} else {
				rootRewards[child.Pos] = child.X
			}
		}
	}
	searchResult := make([]SearchResult, 0)
	for move, score := range rootRewards {
		searchResult = append(searchResult, SearchResult{
			Move:  move,
			Score: score,
		})
	}
	return searchResult
}

func UnusedCards(game *TriPeaks) []Card {
	allCards := NewDeck().Cards
	usedCardsMap := make(map[int]struct{})
	usedCards := game.UsedCards()
	unusedCards := make([]Card, 0, 10)
	for _, card := range usedCards {
		usedCardsMap[card.HashCode()] = struct{}{}
	}
	for _, card := range allCards {
		if _, contains := usedCardsMap[card.HashCode()]; !contains {
			unusedCards = append(unusedCards, card)
		}
	}
	return unusedCards
}

func MctsSelect(game *TriPeaks, node *Node) *Node {
	selected := node
	for game.CardsLeft > 0 {
		moves, _ := game.LegalMoves()
		totalMoves := len(moves)
		if node.GetUnvisitedChild() == nil && len(node.Children) == totalMoves {
			cNode := Ucb1(selected)
			cNode.Data = selected.Data
			selected = cNode
			ApplyNode(game, selected)
		} else {
			break
		}
	}
	return selected
}
func DeterminizeState(node *Node, game *TriPeaks, random *rand.Rand) *Node {
	cNode := NewNode()
	moves, _ := game.LegalMoves()
	usedMoves := make(map[int]struct{})
	unusedMoves := make([]int, 0)
	for _, child := range node.Children {
		usedMoves[child.Pos] = struct{}{}
	}
	for _, move := range moves {
		if _, contains := usedMoves[move]; !contains {
			unusedMoves = append(unusedMoves, move)
		}
	}
	if len(unusedMoves) == 0 {
		// if len(node.Children) == 0 {
		// 	panic("Node doesn't have any children")
		// }
		ind := random.Intn(len(node.Children))
		cNode = node.Children[ind]
		cNode.Data = node.Data
		ApplyNode(game, cNode)
		return cNode
	}
	cNode.Pos = unusedMoves[random.Intn(len(unusedMoves))]
	cNode.Parent = node
	node.Children = append(node.Children, cNode)
	cNode.Data = node.Data

	if cNode.Pos == -1 {
		ind := random.Intn(len(cNode.Data.CardsLeft))
		randCard := cNode.Data.CardsLeft[ind]
		cNode.Data.CardsLeft = RemoveCard(cNode.Data.CardsLeft, ind)
		cNode.SetLeft(Deter{
			Card:        randCard,
			Initialized: true,
		})
	} else {
		leftPos, rightPos := game.CheckReveals(cNode.Pos)
		var (
			leftFound  bool
			rightFound bool
		)
		if leftPos != -1 {
			leftFound = cNode.GetParentDeterminization(leftPos, true)
		}
		if rightPos != -1 {
			rightFound = cNode.GetParentDeterminization(rightPos, false)

		}
		if !leftFound && leftPos != -1 && game.Cards[leftPos].ChildLeft-1 == 0 {
			ind := random.Intn(len(cNode.Data.CardsLeft))
			randCard := cNode.Data.CardsLeft[ind]
			cNode.Data.CardsLeft = RemoveCard(cNode.Data.CardsLeft, ind)
			cNode.SetLeft(Deter{
				Card:        randCard,
				Pos:         leftPos,
				Initialized: true,
			})
		}
		if !rightFound && rightPos != -1 && game.Cards[rightPos].ChildLeft-1 == 0 {
			ind := random.Intn(len(cNode.Data.CardsLeft))
			randCard := cNode.Data.CardsLeft[ind]
			cNode.Data.CardsLeft = RemoveCard(cNode.Data.CardsLeft, ind)
			cNode.SetRight(Deter{
				Card:        randCard,
				Pos:         rightPos,
				Initialized: true,
			})
		}
	}
	ApplyNode(game, cNode)
	return cNode
}
func MctsSimulation(game *TriPeaks, node *Node, random *rand.Rand) float64 {
	for !game.GameOver() {
		node = DeterminizeState(node, game, random)
	}

	clearBonus := 0.0
	if game.CardsLeft == 0 {
		clearBonus = 0.2
	}

	return 0.8 - (float64(game.CardsLeft)/28.0)*0.8 + clearBonus
}

func MctsBackpropagate(node *Node, reward float64) {
	for ; node != nil; node = node.Parent {
		node.X += reward
		node.N++
	}
}

func CopyCards(to, from []Card) {
	for i, card := range from {
		to[i] = card
	}
}

func RemoveCard(cards []Card, i int) []Card {
	cards[i] = cards[len(cards)-1]
	cards = cards[:len(cards)-1]
	return cards
}

func RemoveCardVal(cards []Card, card Card) []Card {
	for i, val := range cards {
		if val.Suit == card.Suit && val.Rank == card.Rank {
			cards = RemoveCard(cards, i)
			return cards
		}
	}
	return cards
}

func Ucb1(node *Node) *Node {
	highest := -1.0
	var selected *Node
	for _, cNode := range node.Children {
		score := math.MaxFloat64
		if cNode.N > 0 {
			score = cNode.X + 2*math.Sqrt(math.Log(float64(node.N))/float64(cNode.N))
		}
		if score > highest {
			selected = cNode
			highest = score
		}
	}
	return selected
}

func ApplyNode(game *TriPeaks, node *Node) {
	if node.Pos == -1 && node.leftDet.Initialized {
		deckLen := game.Deck.Len()
		game.Deck.Cards[deckLen-1] = node.leftDet.Card
		game.Draw()
		if node.leftDet.Card.HashCode() != game.Discard().HashCode() {
			panic("Discard card differs from determinization, should not happen")
		}
		// unusedCards = RemoveCardVal(unusedCards, node.LeftDet.Card)
	} else {
		if leftDet := node.leftDet; leftDet.Initialized {
			game.Cards[leftDet.Pos].Card = leftDet.Card
			// unusedCards = RemoveCardVal(unusedCards, leftDet.Card)
		}
		if rightDet := node.rightDet; rightDet.Initialized {
			game.Cards[rightDet.Pos].Card = rightDet.Card
			// unusedCards = RemoveCardVal(unusedCards, rightDet.Card)
		}
		legalMove := game.Select(node.Pos)
		if !legalMove {
			panic("Game Tree contained illegal move, should not happen")
		}
	}
}

func main() {
	threads := 3 //runtime.NumCPU()
	runtime.GOMAXPROCS(threads)
	deck := NewDeck()
	deck.Shuffle()
	game := NewTripeaks(*deck)
	determinizations := 30 / threads
	trajectories := 45000
	fmt.Printf("Running %d determinizations wtih %d trajectories using %d cores\n", determinizations, trajectories, threads)
	for {
		legalMoves, _ := game.LegalMoves()
		if game.CardsLeft == 0 {
			fmt.Printf("AI won the game!\n")
			break
		} else if len(legalMoves) == 0 {
			fmt.Printf("AI lost the game :(\n")
			break
		}
		fmt.Printf("%s", game)
		fmt.Printf("Cards in deck: %d\t\t\t\tDiscard: %s\n", game.Deck.Len(), game.Discard())
		movesMap := make(map[int]float64)
		movesChan := make(chan []SearchResult, 2)
		for i := 0; i < threads; i++ {
			go func() {
				movesChan <- MctsSearch(game, determinizations, trajectories)
			}()
		}
		highestScore := -1.0
		action := -1
		i := 0
		for moves := range movesChan {
			fmt.Printf("Thread %d results:\n", i)
			for _, move := range moves {
				fmt.Printf("Action %d Score: %f\n", move.Move, move.Score)
				if _, exists := movesMap[move.Move]; exists {
					movesMap[move.Move] += move.Score
				} else {
					movesMap[move.Move] = move.Score
				}
				if movesMap[move.Move] > highestScore {
					highestScore = movesMap[move.Move]
					action = move.Move
				}
			}
			i++
			if i >= threads {
				break
			}
		}
		if action == -1 {
			fmt.Printf("AI Chose to draw a card\n")
			game.Draw()
		} else {
			fmt.Printf("AI Chose to discard card on position: %d\n", action)
			game.Select(action)
		}
	}

}
