package main

import (
	"math/rand"
	"time"

	"github.com/fbundle/sudoku/sudoku"
)

func main() {
	sudoku.ReduceBase(sudoku.N)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	r, store, run := setup()
	sudoku.RegisterRoutes(r, store, rnd)
	run()
}
