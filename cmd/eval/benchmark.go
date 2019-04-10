package main

import (
	"github.com/MatiasLyyra/TriPeaks/deck"
	"github.com/MatiasLyyra/TriPeaks/game"

	"github.com/MatiasLyyra/TriPeaks/mcts"
)

func benchmarkSearch(treeSamples, trajectories, n int, eval mcts.SimulationtEval) {

	for i := 0; i < n; i++ {
		stock := deck.New()
		stock.Shuffle()
		triGame := game.NewTripeaks(*stock)
		for !triGame.GameOver() {
			r := mcts.Search(triGame, treeSamples, trajectories, eval)
			move := r.BestMove()
			if move == -1 {
				triGame.Draw()
			} else {
				triGame.Select(move)
			}
		}
	}
}

func main() {
	benchmarkSearch(10, 3500, 500, mcts.LinearEval)
}
