package udacity

import(
  "dlluncor/myio"
  "strings"
  "strconv"
  "fmt"
)

/* This file deals with doing an A* search using utilities found in
 * sudoku.go.
 */

 // board is a string where the first 9 characters are the first row,
// the next 9 are the second row, etc.
// a dot represents an unknown entry.
func (s *SudokuB) Create(board string) {
  chars := strings.Split(board, "")
  s0 := newCellState()
  for row := 0; row < 9; row++ {
    oneRow := []string{}
    for col := 0; col < 9; col++ {
      ind := (row * 9) + col
      char := chars[ind]
      s0.InitCell(ind, char)
      oneRow = append(oneRow, char)
    }
    fmt.Printf("%v\n", oneRow)
  }
  s0.Visualize()
  s0.UpdatePossib()
  s0.Visualize()
}

func (s *SudokuB) Solve() {
  
}

/*
 * Sudoku solver that uses an astar search found in search.go.
 * Using http://www.goobix.com/games/sudoku/ as a source for puzzles.
 */

 // Solve the problem of putting 15 tiles in order when you only have one blank space.
func Sudoku() {
  r := myio.NewReader()
  board := &SudokuB{}
  board.Init()
  T, _ := strconv.Atoi(r.Read())
  for i := 0; i < T; i++ {
    board.Create(r.Read())
  }
  //PrintBoard(board.board)
  board.Solve()
  fmt.Println("End of sudoku program.")
  SudokuChecks()
}