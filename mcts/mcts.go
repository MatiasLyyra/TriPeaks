package mcts

import (
	"math"
	"math/rand"
	"time"

	"github.com/MatiasLyyra/TriPeaks/deck"
	"github.com/MatiasLyyra/TriPeaks/game"
)

type SearchResult struct {
	Move  int
	Score float64
}

type SearchResults []SearchResult

func (sr SearchResults) BestMove() int {
	var (
		max    float64
		argMax int
	)
	for _, result := range sr {
		if result.Score > max {
			argMax = result.Move
			max = result.Score
		}
	}
	return argMax
}

func Search(tri *game.TriPeaks, determinizations, trajectories int, eval SimulationtEval) SearchResults {
	random := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	unusedCards := hiddenCards(tri)
	rootRewards := make(map[int]float64)
	gameCopy := &(game.TriPeaks{})
	var (
		root            *Node
		unusedCardsCopy []deck.Card
	)
	for i := 0; i < determinizations; i++ {
		root = NewNode()

		for j := 0; j < trajectories; j++ {
			unusedCardsCopy = make([]deck.Card, len(unusedCards))
			gameCopy = tri.Copy()
			deck.Copy(unusedCardsCopy, unusedCards)
			data := &NodeData{
				CardsLeft:          unusedCardsCopy,
				CardsLeftBeginning: gameCopy.CardsLeft,
			}
			root.Data = data
			node := Select(gameCopy, root)
			if !gameCopy.GameOver() {
				node = determinize(node, gameCopy, random)
			}
			reward := simulate(gameCopy, node, random, eval)
			backpropagate(node, reward)
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

func hiddenCards(game *game.TriPeaks) []deck.Card {
	allCards := deck.New().Cards
	usedCardsMap := make(map[int]struct{})
	usedCards := game.UsedCards()
	unusedCards := make([]deck.Card, 0, 10)
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

func Select(game *game.TriPeaks, node *Node) *Node {
	selected := node
	for game.CardsLeft > 0 {
		moves, _ := game.LegalMoves()
		totalMoves := len(moves)
		if node.GetUnvisitedChild() == nil && len(node.Children) == totalMoves {
			cNode := ucb1(selected)
			cNode.Data = selected.Data
			selected = cNode
			applyNode(game, selected)
		} else {
			break
		}
	}
	return selected
}
func determinize(node *Node, game *game.TriPeaks, random *rand.Rand) *Node {
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
		ind := random.Intn(len(node.Children))
		cNode = node.Children[ind]
		cNode.Data = node.Data
		applyNode(game, cNode)
		return cNode
	}
	cNode.Pos = unusedMoves[random.Intn(len(unusedMoves))]
	cNode.Parent = node
	node.Children = append(node.Children, cNode)
	cNode.Data = node.Data

	if cNode.Pos == -1 {
		ind := random.Intn(len(cNode.Data.CardsLeft))
		randCard := cNode.Data.CardsLeft[ind]
		cNode.Data.CardsLeft = deck.Remove(cNode.Data.CardsLeft, ind)
		cNode.LeftDet = Deter{
			Card:        randCard,
			Initialized: true,
		}
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
		if !leftFound && leftPos != -1 && game.Cards[leftPos].FaceDown && game.Cards[leftPos].ChildLeft-1 == 0 {
			ind := random.Intn(len(cNode.Data.CardsLeft))
			randCard := cNode.Data.CardsLeft[ind]
			cNode.Data.CardsLeft = deck.Remove(cNode.Data.CardsLeft, ind)
			cNode.LeftDet = Deter{
				Card:        randCard,
				Pos:         leftPos,
				Initialized: true,
			}
		}
		if !rightFound && rightPos != -1 && game.Cards[rightPos].FaceDown && game.Cards[rightPos].ChildLeft-1 == 0 {
			ind := random.Intn(len(cNode.Data.CardsLeft))
			randCard := cNode.Data.CardsLeft[ind]
			cNode.Data.CardsLeft = deck.Remove(cNode.Data.CardsLeft, ind)
			cNode.RightDet = Deter{
				Card:        randCard,
				Pos:         rightPos,
				Initialized: true,
			}
		}
	}
	applyNode(game, cNode)
	return cNode
}
func simulate(game *game.TriPeaks, node *Node, random *rand.Rand, eval SimulationtEval) float64 {
	for !game.GameOver() {
		node = determinize(node, game, random)
	}
	return eval(node, game)
}

func backpropagate(node *Node, reward float64) {
	for ; node != nil; node = node.Parent {
		node.X += reward
		node.N++
	}
}

func ucb1(node *Node) *Node {
	highest := -1.0
	var selected *Node
	for _, cNode := range node.Children {
		score := math.MaxFloat64
		if cNode.N > 0 {
			score = cNode.X + math.Sqrt(2*math.Log(float64(node.N))/float64(cNode.N))
		}
		if score > highest {
			selected = cNode
			highest = score
		}
	}
	return selected
}

func applyNode(game *game.TriPeaks, node *Node) {
	if node.Pos == -1 && node.LeftDet.Initialized {
		deckLen := game.Stock.Len()
		game.Stock.Cards[deckLen-1] = node.LeftDet.Card
		game.Draw()
		if node.LeftDet.Card.HashCode() != game.Discard().HashCode() {
			panic("Discard card differs from determinization, should not happen")
		}
		// unusedCards = RemoveCardVal(unusedCards, node.LeftDet.Card)
	} else {
		if leftDet := node.LeftDet; leftDet.Initialized {
			game.Cards[leftDet.Pos].Card = leftDet.Card
			// unusedCards = RemoveCardVal(unusedCards, leftDet.Card)
		}
		if rightDet := node.RightDet; rightDet.Initialized {
			game.Cards[rightDet.Pos].Card = rightDet.Card
			// unusedCards = RemoveCardVal(unusedCards, rightDet.Card)
		}
		legalMove := game.Select(node.Pos)
		if !legalMove {
			panic("Game Tree contained illegal move, should not happen")
		}
	}
}
