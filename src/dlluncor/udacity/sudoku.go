package udacity

import(
  "dlluncor/myio"
  "strings"
  "strconv"
  "fmt"
)

// Indices:
// 0 1 2 3 4 5 6 7 8
// 9 10 11 12 13 14 15 16 17
// ...
// 72 73 74 75 76 77 78 79 80

// The state of the board.
type CellState struct {
  possibAns map[int]*[]int // [0] -> []int{4, 5, 6} if the upper left corner can have a 4, 5, or 6.
}

// InitCell(0, ".") if the top left hand cell is unknown.
// InitCell(1, "3") if the second from the top left is the number 3.
func (c *CellState) InitCell(index int, value string) {
  valueAsInt, err := strconv.Atoi(value)
  if err != nil {
    // All values are possible this one is unknown.
    nums := &[]int{1,2,3,4,5,6,7,8,9}
    c.possibAns[index] = nums
  } else {
    c.possibAns[index] = &[]int{valueAsInt}
  }
}

// GetNumber returns the single possible answer for a cell. bool == False
// if there is more than one number, so there are many possibilities.
func GetNumber(nums *[]int) (int, bool) {
  if len(*nums) == 1 {
    return (*nums)[0], true
  }
  return -1, false
}

func deleteFromList(arr *[]int, elToRemove int) *[]int {
  newArr := make([]int, 0)
  for _, el := range *arr {
    if el == elToRemove {
      continue
    }
    newArr = append(newArr, el)
  }
  return &newArr
}

// Get populated when Sudoku first runs.
// horizInds list of other inds in your row, including yourself.
var horizInds = make(map[int][]int)
// verticalInds list of inds in your col, including yourself.
var verticalInds = make(map[int][]int)
// quadrantInds list of inds in your quadrant, including yourself.
var quadrantInds = make(map[int][]int)

// Decreases the number of possibilities based on what it knows about
// the other cells in this row.
func (c *CellState) prune(index int, otherInds []int) {
  for _, otherInd := range otherInds {
    if otherInd == index {
      // Don't consider thyself.
      continue
    }
    val, ok := GetNumber(c.possibAns[otherInd])
    if ok {
      // Now we know this was a solved square, so it can't be a possible
      // solution for this square anymore.
      c.possibAns[index] = deleteFromList(c.possibAns[index], val)
    }
  }
}

// Update which numbers are feasible given the state of this board, e.g.
// prune numbers which are no longer possible given this new configuration.
func (c *CellState) UpdatePossib() {
  for i := 0; i < 80; i++ {
    c.prune(i, horizInds[i])
    c.prune(i, verticalInds[i])
    c.prune(i, quadrantInds[i])
  }
}

// Visualize the board and its possibilities.
func (c *CellState) Visualize() {
  fmt.Println("Visualizing...")
  for row := 0; row < 9; row++ {
    rowInf := []string{}
    for col := 0; col < 9; col++ {
      i := (row * 9) + col
      val, ok := GetNumber(c.possibAns[i])
      curInfo := "."
      if ok {
       curInfo = fmt.Sprintf("%d", val)
      }
      rowInf = append(rowInf, curInfo)
    }
    fmt.Printf("%v\n", rowInf)
  }
}

func newCellState() *CellState {
  return &CellState{
    possibAns: make(map[int]*[]int),
  }
}

// A Sudoku board.
type SudokuB struct {
}

// Init global variables like the index maps.
func (s *SudokuB) Init() {
  byRow := make(map[int][]int) // map of row to indices in that row.
  byCol := make(map[int][]int) // map of col to indices in that col.
  for i := 0; i < 9; i++ {
    byRow[i] = []int{}
    byCol[i] = []int{}
  }

  for row := 0; row < 9; row++ {
    for col := 0; col < 9; col++ {
      index := (row * 9) + col
      byRow[row] = append(byRow[row], index)
      byCol[col] = append(byCol[col], index)
    }
  }

  // Now find all indices in the same row and col.
  for _, indices := range byRow {
    for _, index := range indices {
      horizInds[index] = indices
    }
  }

  for _, indices := range byCol {
    for _, index := range indices {
      verticalInds[index] = indices
    }
  }

  quadrantToInds := map[int][]int{
    0: []int{0, 1, 2, 9, 10, 11, 18, 19, 20},
    1: []int{3, 4, 5, 12, 13, 14, 21, 22, 23},
    2: []int{6, 7, 8, 15, 16, 17, 24, 25, 26},
    3: []int{27, 28, 29, 36, 37, 38, 45, 46, 47},
    4: []int{30, 31, 32, 39, 40, 41, 48, 49, 50},
    5: []int{33, 34, 35, 42, 43, 44, 51, 52, 53},
    6: []int{54, 55, 56, 63, 64, 65, 72, 73, 74},
    7: []int{57, 58, 59, 66, 67, 68, 75, 76, 77},
    8: []int{60, 61, 62, 69, 70, 71, 78, 79, 80},
  }
  for _, indices := range quadrantToInds {
    for _, index := range indices {
      quadrantInds[index] = indices
    }
  }
}

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

// DBG

func printInds() {
  for row := 0; row < 9; row++ {
    vals := []int{}
    for col := 0; col < 9; col++ {
      ind := (row * 9) + col
      vals = append(vals, ind)
    }
    fmt.Printf("%v\n", vals)
  }
}

func SudokuChecks() {
  //printInds()
  //fmt.Printf("%v", horizInds)
  //fmt.Printf("%v", verticalInds)
  //fmt.Printf("%v", quadrantInds)
}

// DBG

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