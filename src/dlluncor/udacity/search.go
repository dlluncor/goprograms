package udacity

import (
  "dlluncor/myio"
  "strconv"
  "fmt"
  "strings"
)

// GraphSearch implements a general graph search application, such as A* or BFS
// or DFS.
/*
func GraphSearch(fronteir Fronteir, explored Explored, searcher Searcher) bool {
  for {
    if fronteir.IsEmpty() {
    	return false
    }
    node := fronteir.RemoveChoice()
    s = node.state // What node we are at currently in this path.
    explored.Add(s) // We've now seen this node.
    if searcher.IsGoal(node) {
    	return true
    }
    actions := searcher.NextActions(node)
    for a, _ range actions {
    	if !fronteir.Contains(a) && !explored.Contains(a) { 
    	  fronteir.AddAction(a)
    	}
    }
  }
}
*/

 // The number 16 indicates blank.
var (
  BLANK = 16
)

// General utilities.

// FromBoard returns the board as a string.
func FromBoard(board []int) string{
  str := ""
  for index, integer := range board {
    space := " "
    if index == 0 {
      space = ""
    }
    str += fmt.Sprintf(space + "%d", integer)
  }
  return str
}

// Converts a string to a integer board.
func ToBoard(board string) []int {
  boardNumArr := make([]int, 16)
  numStrs := strings.Split(board, " ")
  for index, numStr := range numStrs {
    boardNumArr[index], _ = strconv.Atoi(numStr)
  }
  return boardNumArr
}


// Takes two indices and swaps their values, returning a new board.
/*
func SwapThem(blankIndex int, filledIndex int, []int board) []int {
  newBoard := make([]int, 16)
  copy(newBoard, board)
  newBoard[blankIndex] = board[filledIndex]
  newBoard[filledIndex] = BLANK
  return newBoard
}
*/

/*
// Returns an array of possible new boards given the current board state.
func FindAdjacents(board []int) [][]int {
  // Find the location of the blank.
  numToLoc := make(map[int]int)
  for index, num := range board {
    numToLoc[num] = index
  }
  blankIndex := numToLoc[16]

  // Find adjacent locations to this one.
  possibIndices := []int{blankIndex-1, blankIndex-4, blankIndex+1, blankIndex+4}
  newBoards := make([][]int, 0)
  for _, possibInd := range possibIndices {
    if possibInd < 16 {
      // Here we can swap.
      // Create a new array of elements which we can swap.
      newBoard := SwapThem(blankIndex, possibInd, board)
      newBoards = append(newBoards, newBoard)
    }
  }
  return newBoards
}

type BFronteir struct {
  strToBoard map[string] []int // Maps from string rep of board to actual board.
} 

// Adds a board to the fronteir.
func (b *BFronteir) Add(board []int) {
  stringRep := StringBoard(board)
  b.strToBoard[stringRep] = board
}

func (b *BFronteir) IsEmpty() bool {
  return len(b.strToBoard) == 0
}

type BoardSolver struct {
  fronteir BFronteir
}

func (bs *BoardSolver) Init(board string) {
  // Fronteir consists of current board.
  bs.fronteir = &BFronteir{}
  bs.fronteir.Add(board)
}

func (bs *BoardSolver) IsEmpty() {
  return bs.fronteir.IsEmpty()
}

func (bs *BoardSolver) RemoveChoice() {
  return bs.fronteir.RemoveChoice()
}
*/

// General board.

type Board struct {
  board string
}

/*
func (b *Board) Solve() {
  bs := &BoardSolver{}
  bs.Init(b.board)
  //GraphSearch(bs, bs, bs)
}
*/

func (b *Board) Create(r myio.Reader) {
  boardArr := make([]int, 16)

  // Read board from input.
  for i := 0; i < 4; i++ {
    nums := strings.Split(r.Read(), " ")
	  for j, num := range nums {
	    numI, _ := strconv.Atoi(num)
      k := i * 4 + j
      boardArr[k] = numI
	   }
  }

  b.board = FromBoard(boardArr)
}

func (b *Board) Print() {
  output := ""
  boardArr := ToBoard(b.board)
  for i := 0; i < 4; i++ {
    line := ""
    for j := 0; j < 4; j++ {
      tab := "\t"
      if j == 0 {
        tab = ""
      }
      k := i * 4 + j
      line += fmt.Sprintf(tab + "%d", boardArr[k])
    }
    output += line + "\n"
  }
  fmt.Printf("Current board: \n")
  fmt.Printf(output)
}

func FifteenNums() {
  r := myio.NewReader()
  board := &Board{}
  board.Create(r)
  board.Print()
  //board.Solve()
}
