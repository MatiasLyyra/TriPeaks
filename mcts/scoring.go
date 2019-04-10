package mcts

import (
	"math"

	"github.com/MatiasLyyra/TriPeaks/game"
)

type SimulationtEval func(*Node, *game.TriPeaks) float64

func BinaryEval(node *Node, tri *game.TriPeaks) float64 {
	if tri.CardsLeft == 0 {
		return 1
	}
	return 0
}

func LinearEval(node *Node, tri *game.TriPeaks) float64 {
	return float64(tri.CardsLeft) / 28.0
}

func ScoreEval(node *Node, tri *game.TriPeaks) float64 {
	score := 0.0125*float64(tri.Score) + 0.25
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}
	return score
}

func ScoreLogEval(node *Node, tri *game.TriPeaks) float64 {
	return math.Log(1 + math.Exp(float64(tri.Score)))
}

func ScoreSigmoidEval(node *Node, tri *game.TriPeaks) float64 {
	return 1 / (1 + math.Exp(-float64(tri.Score)/15))
}
