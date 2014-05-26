package main

import (
	"fmt"
	"os"

	"dlluncor/ir"
	"dlluncor/server"
	"dlluncor/spoj"
	"dlluncor/udacity"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(`Usage: ./cmd 0 {args for your program}`)
	}
	prog := os.Args[1]
	switch prog {
	case "0":
		//spoj.Bitmap()
		//spoj.Recaman()
		//spoj.EditDistance()
		//spoj.Party()
		//spoj.GreatBall()
		//spoj.Sqrt()
		spoj.MoveToInvert()
	case "1":
		//udacity.Search()
		//udacity.FifteenNums()
	case "2":
		server.Serve()
		//spoj.Concurrency()
		//spoj.WordRacer()
		//spoj.Scrabble()
	case "3":
		udacity.Sudoku()
	case "4":
		ir.BuildIndex()
	case "5":
		ir.MainScorer(1, os.Args)
	default:
		panic(fmt.Sprintf("Unrecognized program int: %v\n", prog))
	}
}
