package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/MatiasLyyra/TriPeaks/deck"
	"github.com/MatiasLyyra/TriPeaks/game"

	"github.com/MatiasLyyra/TriPeaks/mcts"
)

type BenchmarkOptions struct {
	Name             string
	Threads          int
	N                int
	Determinizations int
	Trajectories     int
	Eval             mcts.SimulationtEval
}
type ToCSV interface {
}
type BenchmarkResult struct {
	Name             string
	N                int
	Determinizations int
	Trajectories     int
	GamesWon         int
	CardsCleared     int
	Points           int
}

func WriteCsv(results []BenchmarkResult, w io.Writer) {
	_, err := w.Write([]byte("name,n,determinizations,trajectories,games_won,cards_cleared,points\n"))
	if err != nil {
		log.Printf("write error: %s", err)
	}
	for _, r := range results {
		csv := fmt.Sprintf("%s,%d,%d,%d,%d,%d,%d\n", r.Name, r.N, r.Determinizations, r.Trajectories, r.GamesWon, r.CardsCleared, r.Points)
		_, err = w.Write([]byte(csv))
		if err != nil {
			log.Printf("write error: %s", err)
		}
	}
}

type AiFunc func(*game.TriPeaks, BenchmarkOptions) int

func benchmarkSearch(options BenchmarkOptions, ai AiFunc) BenchmarkResult {
	r := BenchmarkResult{
		Name:             options.Name,
		Determinizations: options.Determinizations * options.Threads,
		Trajectories:     options.Trajectories,
		N:                options.N,
	}
	for i := 0; i < options.N; i++ {
		stock := deck.New()
		stock.Shuffle()
		triGame := game.NewTripeaks(*stock)
		for !triGame.GameOver() {
			move := ai(triGame, options)
			if move == -1 {
				triGame.Draw()
			} else {
				triGame.Select(move)
			}
		}
		r.Points += triGame.Score
		r.CardsCleared += 28 - triGame.CardsLeft
		if triGame.CardsLeft == 0 {
			r.GamesWon++
		}
		fmt.Printf("%s progress: %.2f %%\n", options.Name, math.Round(float64(i)/float64(options.N)*10000)/100)
	}
	return r
}
func random(triGame *game.TriPeaks, options BenchmarkOptions) int {
	legals, _ := triGame.LegalMoves()
	return legals[rand.Intn(len(legals))]
}
func mctsSearch(triGame *game.TriPeaks, options BenchmarkOptions) int {
	movesMap := make(map[int]float64)
	movesChan := make(chan []mcts.SearchResult, 2)
	for i := 0; i < options.Threads; i++ {
		go func() {
			movesChan <- mcts.Search(triGame, options.Determinizations, options.Trajectories, options.Eval)
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
		if i >= options.Threads {
			break
		}
	}
	return action
}

func main() {
	runtime.GOMAXPROCS(10)
	results := make([]BenchmarkResult, 0, 10)

	rand.Seed(time.Now().UnixNano())
	options := BenchmarkOptions{
		Name:             "Random",
		N:                500,
		Threads:          1,
		Determinizations: 0,
		Trajectories:     0,
		Eval:             nil,
	}
	results = append(results, benchmarkSearch(options, random))

	options = BenchmarkOptions{
		Name:             "LinearEval 1",
		N:                500,
		Threads:          10,
		Determinizations: 1,
		Trajectories:     1500,
		Eval:             mcts.LinearEval,
	}
	results = append(results, benchmarkSearch(options, mctsSearch))
	options = BenchmarkOptions{
		Name:             "LinearEval 2",
		N:                500,
		Threads:          10,
		Determinizations: 5,
		Trajectories:     2500,
		Eval:             mcts.LinearEval,
	}
	results = append(results, benchmarkSearch(options, mctsSearch))
	options = BenchmarkOptions{
		Name:             "LinearEval 3",
		N:                500,
		Threads:          10,
		Determinizations: 10,
		Trajectories:     3500,
		Eval:             mcts.LinearEval,
	}
	results = append(results, benchmarkSearch(options, mctsSearch))

	options = BenchmarkOptions{
		Name:             "BinaryEval 1",
		N:                500,
		Threads:          10,
		Determinizations: 1,
		Trajectories:     1500,
		Eval:             mcts.BinaryEval,
	}
	results = append(results, benchmarkSearch(options, mctsSearch))
	options = BenchmarkOptions{
		Name:             "BinaryEval 2",
		N:                500,
		Threads:          10,
		Determinizations: 5,
		Trajectories:     2500,
		Eval:             mcts.BinaryEval,
	}
	results = append(results, benchmarkSearch(options, mctsSearch))
	options = BenchmarkOptions{
		Name:             "BinaryEval 3",
		N:                500,
		Threads:          10,
		Determinizations: 10,
		Trajectories:     3500,
		Eval:             mcts.BinaryEval,
	}
	results = append(results, benchmarkSearch(options, mctsSearch))

	options = BenchmarkOptions{
		Name:             "ScoreEval 1",
		N:                500,
		Threads:          10,
		Determinizations: 1,
		Trajectories:     1500,
		Eval:             mcts.ScoreEval,
	}
	results = append(results, benchmarkSearch(options, mctsSearch))
	options = BenchmarkOptions{
		Name:             "ScoreEval 2",
		N:                500,
		Threads:          10,
		Determinizations: 5,
		Trajectories:     2500,
		Eval:             mcts.ScoreEval,
	}
	results = append(results, benchmarkSearch(options, mctsSearch))
	options = BenchmarkOptions{
		Name:             "ScoreEval 3",
		N:                500,
		Threads:          10,
		Determinizations: 10,
		Trajectories:     3500,
		Eval:             mcts.ScoreEval,
	}
	results = append(results, benchmarkSearch(options, mctsSearch))

	options = BenchmarkOptions{
		Name:             "ScoreSigmoidEval 1",
		N:                500,
		Threads:          10,
		Determinizations: 1,
		Trajectories:     1500,
		Eval:             mcts.ScoreSigmoidEval,
	}
	results = append(results, benchmarkSearch(options, mctsSearch))
	options = BenchmarkOptions{
		Name:             "ScoreSigmoidEval 2",
		N:                500,
		Threads:          10,
		Determinizations: 5,
		Trajectories:     2500,
		Eval:             mcts.ScoreSigmoidEval,
	}
	results = append(results, benchmarkSearch(options, mctsSearch))
	options = BenchmarkOptions{
		Name:             "ScoreSigmoidEval 3",
		N:                500,
		Threads:          10,
		Determinizations: 10,
		Trajectories:     3500,
		Eval:             mcts.ScoreSigmoidEval,
	}
	results = append(results, benchmarkSearch(options, mctsSearch))

	saveResults(results)
}

func saveResults(r []BenchmarkResult) {
	t := time.Now()
	path := fmt.Sprintf("./benchmarks/benchmark_eval_%d_%02d_%02d_%02d_%02d.csv", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		log.Printf("failed to open / create file %s, writing to stdout\n", path)
		WriteCsv(r, os.Stdout)
	} else {
		WriteCsv(r, f)
	}
}
