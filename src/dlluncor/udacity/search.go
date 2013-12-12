package udacity

import (
  "dlluncor/myio"
  "strconv"
  "fmt"
  "strings"
)

// Interfaces here are really nodes which the implementation should pass around.
type Fronteir interface {
  IsEmpty() bool
  RemoveChoice() interface{}
  Contains(node interface{}) bool
  Add(node interface{})
}

type Explored interface {
  Add(node interface{})
  Contains(node interface{}) bool
}

type Searcher interface {
  IsGoal(node interface{}) bool
  NextActions(node interface{}) []interface{}
}

// GraphSearch implements a general graph search application, such as A* or BFS
// or DFS.
func GraphSearch(fronteir Fronteir, explored Explored, searcher Searcher) (interface{}, int) {
  var i = 0
  for {
    //if i == 100000 {
    //  break
   // }
    if fronteir.IsEmpty() {
    	return nil, i
    }
    node := fronteir.RemoveChoice()
    /*mnode := node.(*BNode)
    if i % 10000 == 0 {
      fmt.Printf("At node. Cost: %v F: %v H: %v\n", mnode.cost, mnode.f, mnode.h)
    }*/
    explored.Add(node) // We've now seen this node.
    if searcher.IsGoal(node) {
    	return node, i
    }
    actions := searcher.NextActions(node)
    for _, a := range actions {
    	if !fronteir.Contains(a) && !explored.Contains(a) { 
    	  fronteir.Add(a)
    	}
    }
    i += 1
  }
  return nil, i
}

// General utilities for the 15 tile problem.

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

// Prints the board in human readable format.
func PrintBoard(board string) {
  output := ""
  boardArr := ToBoard(board)
  for i := 0; i < 4; i++ {
    line := ""
    for j := 0; j < 4; j++ {
      tab := "\t"
      if j == 0 {
        tab = ""
      }
      k := i * 4 + j
      number := " "
      tileVal := boardArr[k]
      if tileVal != 16 {
        number = fmt.Sprintf("%d", tileVal)
      }
      line += tab + number 
    }
    output += line + "\n"
  }
  fmt.Printf(output)
}

// Path we took to solve the puzzle.
func PrintPath(node *BNode) {
  defer fmt.Println("--------------------------")
  defer PrintBoard(node.state)
  if node.parent == nil {
    return
  }
  PrintPath(node.parent)
}



// Takes two indices and swaps their values, returning a new board.

func SwapThem(blankIndex int, filledIndex int, board []int) []int {
  newBoard := make([]int, 16)
  copy(newBoard, board)
  newBoard[blankIndex] = board[filledIndex]
  newBoard[filledIndex] = BLANK
  return newBoard
}

// The number 16 indicates blank.
var (
  BLANK = 16
  BLANK_TO_ADJACENTS = map[int][]int{
    0: {1, 4},
    1: {0, 2, 5},
    2: {1, 3, 6},
    3: {2, 7},
    4: {0, 5, 8},
    5: {1, 4, 6, 9},
    6: {2, 5, 7, 10},
    7: {3, 6, 11},
    8: {4, 9, 12},
    9: {5, 8, 10, 13},
    10: {6, 9, 11, 14},
    11: {7, 10, 15},
    12: {8, 13},
    13: {9, 12, 14},
    14: {10, 13, 15},
    15: {11, 14},
  }
)

// Returns an array of possible new boards given the current board state.
func FindAdjacents(board string) []string {
  // Find the location of the blank.
  numToLoc := make(map[int]int)
  boardArr := ToBoard(board)
  for index, num := range boardArr {
    numToLoc[num] = index
  }
  blankIndex := numToLoc[16]

  // Find adjacent locations to this one.
  possibIndices := BLANK_TO_ADJACENTS[blankIndex]
  newBoards := make([]string, 0)
  for _, possibInd := range possibIndices {
    // Here we can swap.
    // Create a new array of elements which we can swap.
    newBoard := SwapThem(blankIndex, possibInd, boardArr)
    newBoards = append(newBoards, FromBoard(newBoard))
  }
  return newBoards
}

// General board problem driver.

type Board struct {
  board string
}

func (b *Board) Solve() {
  bs := &BoardSolver{}
  bs.Init(b.board)
  idest, numGuesses := GraphSearch(bs.fronteir, bs.explored, bs)
  dest := idest.(*BNode)
  if dest != nil {
    fmt.Printf("***********")
    fmt.Printf("Solved it with cost %v. f: %v. Guesses: %v\n", 
               dest.cost, dest.f, numGuesses)
    fmt.Printf("Path to get there:\n")
    PrintPath(dest)
  } else {
    fmt.Println("There is no way to solve this puzzle.\n")
  }
}

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

// Solve the problem of putting 15 tiles in order when you only have one blank space.
func FifteenNums() {
  r := myio.NewReader()
  T, _ := strconv.Atoi(r.Read())
  for i := 0; i < T; i++ {
    board := &Board{}
    board.Create(r)
    //PrintBoard(board.board)
    board.Solve()
  }
}
