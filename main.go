package main

import (
	"fmt"
	"runtime"

	"github.com/MatiasLyyra/TriPeaks/deck"
	"github.com/MatiasLyyra/TriPeaks/game"
	"github.com/MatiasLyyra/TriPeaks/mcts"
)

func main() {
	threads := runtime.NumCPU()
	runtime.GOMAXPROCS(threads)
	deck := deck.New()
	deck.Shuffle()
	game := game.NewTripeaks(*deck)
	determinizations := 72 / threads
	trajectories := 5000
	fmt.Printf("Running %d determinizations wtih %d trajectories using %d cores\n", determinizations, trajectories, threads)
	for {
		legalMoves, _ := game.LegalMoves()
		fmt.Printf("%s", game)
		fmt.Printf("Cards in deck: %d\tScore: %d\t\tDiscard: %s\n", game.Stock.Len(), game.Score, game.Discard())
		if game.CardsLeft == 0 {
			fmt.Printf("AI won the game!\n")
			break
		} else if len(legalMoves) == 0 {
			fmt.Printf("AI lost the game :(\n")
			break
		}

		movesMap := make(map[int]float64)
		movesChan := make(chan []mcts.SearchResult, 2)
		for i := 0; i < threads; i++ {
			go func() {
				movesChan <- mcts.Search(game, determinizations, trajectories, mcts.ScoreSigmoidEval)
			}()
		}
		highestScore := -1.0
		action := -1
		i := 0
		for moves := range movesChan {
			for _, move := range moves {
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
		for move, score := range movesMap {
			fmt.Printf("Move %d Score %f\n", move, score/float64(determinizations*threads*trajectories))
		}
		if action == -1 {
			fmt.Printf("AI Chose to draw a card\n")
			game.Draw()
		} else {
			fmt.Printf("AI Chose to discard %s on position: %d\n", game.Cards[action], action)
			game.Select(action)
		}
	}
}
