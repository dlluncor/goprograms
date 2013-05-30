package udacity

import (
  "dlluncor/myio"
  "dlluncor/container"
  "container/heap"
  "strconv"
  "fmt"
  "strings"
  "time"
)

type Fronteir interface {
  IsEmpty() bool
  RemoveChoice() *BNode
  Contains(node *BNode) bool
  Add(node *BNode)
}

type Explored interface {
  Add(node *BNode)
  Contains(node *BNode) bool
}

type Searcher interface {
  IsGoal(node *BNode) bool
  NextActions(node *BNode) []*BNode
}

// GraphSearch implements a general graph search application, such as A* or BFS
// or DFS.
func GraphSearch(fronteir Fronteir, explored Explored, searcher Searcher) *BNode {
  var i = 0
  for {
    if i == 100000 {
      break
    }
    if fronteir.IsEmpty() {
    	return nil
    }
    node := fronteir.RemoveChoice()
    explored.Add(node) // We've now seen this node.
    if searcher.IsGoal(node) {
    	return node
    }
    actions := searcher.NextActions(node)
    for _, a := range actions {
    	if !fronteir.Contains(a) && !explored.Contains(a) { 
    	  fronteir.Add(a)
    	}
    }
    i += 1
  }
  return nil
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

// The distance this board is from being the correct answer. Just count the 
// number of misplaced tiles.
func H1Dist(board string) int {
  numArr := ToBoard(board)
  distance := 0
  for index, tileVal := range numArr {
    if tileVal - 1 != index {
      distance += 1
    }
  }
  return distance
}

// Takes two indices and swaps their values, returning a new board.

func SwapThem(blankIndex int, filledIndex int, board []int) []int {
  newBoard := make([]int, 16)
  copy(newBoard, board)
  newBoard[blankIndex] = board[filledIndex]
  newBoard[filledIndex] = BLANK
  return newBoard
}


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


type BNode struct {
  parent *BNode // which state did I eminate from.
  state string // string representation of the board.
  cost int
  f int // number of hops to get to this state.
  h int // how far am I from my goal.
}

type BExplored struct {
  boardMap map[string] bool // Whether this board state has already been explored.
}

func (e *BExplored) Add(node *BNode) {
  e.boardMap[node.state] = true
}

func (e *BExplored) Contains(node *BNode) bool {
  fmt.Printf("Explored size: %v\n", len(e.boardMap))
  _, ok := e.boardMap[node.state]
  return ok 
}

// Fronteir of what has not been explored yet.
type BFronteir struct {
  boardMap map[string] *BNode // Map of boards to their state in the fronteir.
  queue *container.PriorityQueue // Keeps track of the next node with the least cost.
}

// Adds a board to the fronteir.
func (b *BFronteir) Add(node *BNode) {
  // Calculate costs before adding to the fronteir.
  if node.parent != nil {
    node.f = node.parent.f + 1  // We took one more hop, or step to get here.
    node.h = H1Dist(node.state)
    node.cost = node.f + node.h
  }
  b.boardMap[node.state] = node
  // Add the state to a priority queue as well.
  item := &container.Item{
    Value: node.state,
    Priority: node.cost,
  }
  heap.Push(b.queue, item)
}

func (b *BFronteir) IsEmpty() bool {
  return len(b.boardMap) == 0
}

func (b *BFronteir) Contains(node *BNode) bool {
  _, ok := b.boardMap[node.state]
  return ok
}


func(b *BFronteir) RemoveChoice() *BNode {
  before := time.Now()
  // Pop the item with the lowest cost from the heap, fast!!!
  item := heap.Pop(b.queue).(*container.Item)
  lNode := b.boardMap[item.Value] // Node with the lowest cost.
  fmt.Printf("Fronteir size: %v\n", len(b.boardMap))
  fmt.Printf("Best choice. Cost: %v. H: %v\n", lNode.cost, lNode.h)
  PrintBoard(lNode.state)
  delete(b.boardMap, lNode.state)
  fmt.Printf("Time elapsed: %v\n", time.Since(before))
  return lNode
}


// Helps solve the board problem.
type BoardSolver struct {
  fronteir *BFronteir
  explored *BExplored
}

// Init initializes the board solver given the first state.
func (bs *BoardSolver) Init(board string) {
  // Fronteir consists of current board.
  bs.fronteir = &BFronteir{
    boardMap: make(map[string] *BNode),
    queue: &container.PriorityQueue{},
  }
  heap.Init(bs.fronteir.queue)
  node := &BNode{
    parent: nil,
    state: board,
    f: 0,
    h: H1Dist(board), 
  }
  node.cost = node.f + node.h
  bs.fronteir.Add(node)
  bs.explored = &BExplored{
    boardMap: make(map[string] bool),
  }
}

func (bs *BoardSolver) IsGoal(node *BNode) bool {
  return H1Dist(node.state) == 0
}

func (bs *BoardSolver) NextActions(node *BNode) []*BNode {
  // Returns the list of BNodes to explore next.
  nextNodes := make([]*BNode, 0)
  adjacentStates := FindAdjacents(node.state)
  fmt.Println("Next states:")
  for _, state := range adjacentStates {
    nextNode := &BNode{
      parent: node,
      state: state,
    }
    //PrintBoard(state)
    nextNodes = append(nextNodes, nextNode)
  }
  return nextNodes
}

// General board problem driver.

type Board struct {
  board string
}


func (b *Board) Solve() {
  bs := &BoardSolver{}
  bs.Init(b.board)
  dest := GraphSearch(bs.fronteir, bs.explored, bs)
  if dest != nil {
    fmt.Printf("***********")
    fmt.Printf("Solved it with cost %v. f: %v\n", dest.cost, dest.f)
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

func FifteenNums() {
  r := myio.NewReader()
  T, _ := strconv.Atoi(r.Read())
  for i := 0; i < T; i++ {
    board := &Board{}
    board.Create(r)
    PrintBoard(board.board)
    board.Solve()
  }
}
