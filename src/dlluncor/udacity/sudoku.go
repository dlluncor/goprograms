package udacity

import (
	"dlluncor/myio"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const numSquares = 81 // Not 80 haha!!!

// Utilities

// GetNumber returns the single possible answer for a cell. bool == False
// if there is more than one number, so there are many possibilities.
// Returns the length of the array if there is more than one.
func GetNumber(nums *[]int) (int, bool) {
	if len(*nums) == 1 {
		return (*nums)[0], true
	}
	return len(*nums), false
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

// Indices:
// 0 1 2 3 4 5 6 7 8
// 9 10 11 12 13 14 15 16 17
// ...
// 72 73 74 75 76 77 78 79 80

// The state of the board.
type CellState struct {
	possibAns map[int]*[]int // [0] -> []int{4, 5, 6} if the upper left corner can have a 4, 5, or 6.
	unsolved  map[int]bool   // Map of integers which are not solved for yet.
}

func copyInts(fromArr *[]int) *[]int {
	// TODO(dlluncor): replace when realizing what the eff is going on.
	newArr := make([]int, len(*fromArr))
	for index, el := range *fromArr {
		newArr[index] = el
	}
	return &newArr
}

func (c *CellState) copy() *CellState {
	newS := newCellState()
	for key, value := range c.possibAns {
		newS.possibAns[key] = copyInts(value)
	}
	for key, value := range c.unsolved {
		newS.unsolved[key] = value
	}
	return newS
}

// stringified version of this cell state.
func (c *CellState) ToString() string {
	vals := make([]string, numSquares)
	for i := 0; i < numSquares; i++ {
		vals[i] = fmt.Sprintf("%v", *c.possibAns[i])
	}
	return strings.Join(vals, "-")
}

// InitCell(0, ".") if the top left hand cell is unknown.
// InitCell(1, "3") if the second from the top left is the number 3.
func (c *CellState) InitCell(index int, value string) {
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		// All values are possible this one is unknown.
		nums := &[]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		c.possibAns[index] = nums
	} else {
		c.possibAns[index] = &[]int{valueAsInt}
	}
}

// A state is invalid if ANY of the squares have 0 possibilities, meaning that we can't
// do anything meaningful (not solvable.)
func (c *CellState) IsInvalid() bool {
	for i := 0; i < numSquares; i++ {
		noAnswer := len(*c.possibAns[i]) == 0
		if noAnswer {
			return true
		}
	}
	return false
}

// Number of squares which have not been already solved for.
func (c *CellState) NumUnsolved() (int32, int32) {
	unsolved := int32(0)
	possibilities := int32(0)
	for i, _ := range c.unsolved {
		// TODO(dlluncor): Number of solved or unsolved values should be cached
		// somewhere...
		numRemaining, isAns := GetNumber(c.possibAns[i])
		if !isAns {
			unsolved++
			possibilities += int32(numRemaining)
		}
	}
	return unsolved, possibilities
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
func (c *CellState) prune(index int, otherInds []int) bool {
	hasAnswer := false
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
			// Did we generate a new answer here?
			_, okAfter := GetNumber(c.possibAns[index])
			if okAfter {
				// We generated an example in which there is now a solution where
				// there was not one before.
				// No need to check other neighbors because there is nothing left
				// to check for.
				hasAnswer = true
			}
		}
	}
	return hasAnswer
}

func (c *CellState) pruneOther(index int, otherInds []int) bool {
	// There is another case where you are the only person
	// who contains the possibility to answer with that number in which case
	// you get it.
	// 
	// Example:
	// 1 2 3
	// 4 5 6
	// (7/9) (7/8) (7/9)
	// 8 9
	//
	//
	//     8
	exists := make(map[int]bool) // indices which exist in the other quadrants.
	for _, otherInd := range otherInds {
		if otherInd == index {
			// Don't consider thyself.
			continue
		}
		for _, possibOtherAns := range *c.possibAns[otherInd] {
			exists[possibOtherAns] = true
		}
	}

	myNums := c.possibAns[index]
	for _, myNum := range *myNums {
		if _, ok := exists[myNum]; !ok {
			// If only see this number in our square and in nobody else's then
			// we have an answer for us and we need to clear it from everyone
			// else's list of possibilities.
			c.possibAns[index] = &[]int{myNum}
			return true
		}
	}
	return false
}

// IsSolved means every single cell has only one possibility.
func (c *CellState) IsSolved() bool {
	for i, _ := range c.unsolved {
		// TODO(dlluncor): Really only want the cells which are not solved yet.
		_, isAns := GetNumber(c.possibAns[i])
		if !isAns {
			return false
		}
	}
	return true
}

// Verifies that the Sudoku board is indeed a valid one by checking all quadrant
// and rows that the indices represent the numbers 0 to 9.
func (c *CellState) DidISolveTheBoard() bool {
	hasNineNums := func(indices []int) bool {
		numsMap := make(map[int]bool)
		for _, index := range indices {
			num := (*c.possibAns[index])[0]
			numsMap[num] = true
		}
		return len(numsMap) == 9
	}
	for i := 0; i < numSquares; i++ {
		if hasNineNums(horizInds[i]) && hasNineNums(verticalInds[i]) && hasNineNums(quadrantInds[i]) {
			continue
		}
		fmt.Printf("Index %d is wrong.", i)
		return false
	}
	return true
}

// Neighbors produces all neighbor states which result from making one move.
func (c *CellState) Neighbors() []*CellState {
	neighStates := []*CellState{}
	for i, _ := range c.unsolved {
		// TODO(dlluncor): Really only want the cells which are not solved yet.
		possibs := c.possibAns[i]
		if len(*possibs) == 1 {
			// This is already solved for don't need to create neighbors.
			continue
		} else {
			// Need to pick one of the answers and then run with it.
			for _, possib := range *possibs {
				newState := c.copy()
				newState.possibAns[i] = &[]int{possib} // So now we've chosen this to be the answer.
				neighStates = append(neighStates, newState)
			}
		}
	}
	return neighStates
}

// Update which numbers are feasible given the state of this board, e.g.
// prune numbers which are no longer possible given this new configuration.
func (c *CellState) UpdatePossib() {
	//TODO(dlluncor): There is something wrong with this logic... as in
	// if you change when we break out of the loop, it keeps finding incorrect
	// states not sure why.

	// TODO(dlluncor): Stop immediately when any of the states is invalid
	// and do not continue further nor even put the object on the frontier.
	hasNewAnswer := false
	for i, _ := range c.unsolved {
		// TODO(dlluncor): Only iterate over unsolved indices.
		// Is this cell already solved for? In which case, skip over it.
		_, isAns := GetNumber(c.possibAns[i])
		if isAns {
			delete(c.unsolved, i)
			continue
		}
		ans := c.prune(i, horizInds[i])
		if ans {
			hasNewAnswer = true
		}
		ans = c.prune(i, verticalInds[i])
		if ans {
			hasNewAnswer = true
		}
		ans = c.prune(i, quadrantInds[i])
		if ans {
			hasNewAnswer = true
		}
		ans = c.pruneOther(i, horizInds[i])
		if ans {
			hasNewAnswer = true
		}
		ans = c.pruneOther(i, verticalInds[i])
		if ans {
			hasNewAnswer = true
		}
		ans = c.pruneOther(i, quadrantInds[i])
		if ans {
			hasNewAnswer = true
		}
	}
	if hasNewAnswer {
		// If we found a solution in any one of these results, we need to
		// re-run pruning.
		c.UpdatePossib()
	}
}

func (c *CellState) RunForInitialState() {
	// Things that only need to be run once.
	for i := 0; i < numSquares; i++ {
		_, isAns := GetNumber(c.possibAns[i])
		if !isAns {
			c.unsolved[i] = true
		}
	}
	fmt.Printf("Num unsolved to start: %d\n", len(c.unsolved))
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

// Visualizes board and all possibilites.
func (c *CellState) VisualizeAll() {
	c.Visualize()
	for row := 0; row < 9; row++ {
		rowInf := []string{}
		for col := 0; col < 9; col++ {
			i := (row * 9) + col
			possibs := *c.possibAns[i]
			if len(possibs) == 1 {
				curInfo := fmt.Sprintf("(%d)", possibs[0])
				rowInf = append(rowInf, curInfo)
			} else {
				curInfo := fmt.Sprintf("(%v)", possibs)
				rowInf = append(rowInf, curInfo)
			}
		}
		fmt.Printf("%v\n", rowInf)
	}
}

func (c *CellState) PrintAsInput() {
	output := ""
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
		output += strings.Join(rowInf, "")
	}
	fmt.Printf("%v\n", output)
}

func newCellState() *CellState {
	return &CellState{
		possibAns: make(map[int]*[]int),
		unsolved:  make(map[int]bool),
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

// board is a string where the first 9 characters are the first row,
// the next 9 are the second row, etc.
// a dot represents an unknown entry.
func (s *SudokuB) Create(board string) *CellState {
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
	return s0
}

func (s *SudokuB) Solve(boardStr string, defHasSolution bool) {
	sol := &SudokuSolver{
		guess:   int32(0),
		prevNow: time.Now(),
	}
	state0 := s.Create(boardStr)
	sol.Init(state0)
	idest, numGuesses := GraphSearch(sol.frontier, sol.explored, sol)
	if idest != nil {
		dest := idest.(*SNode)
		if !dest.state.DidISolveTheBoard() {
			dest.state.VisualizeAll()
			log.Fatalf("You didn't really solve the board correctly dummy!!")
		}
		cost := dest.h + dest.f
		fmt.Printf("***********")
		fmt.Printf("Solved it with cost %v. f: %v. Guesses: %v\n",
			cost, dest.f, numGuesses)
		fmt.Printf("Solution board:\n")
		dest.state.Visualize()
	} else {
		if defHasSolution {
			log.Fatalf("Error there is a solution to this puzzle!!!\n")
		}
		fmt.Println("There is no way to solve this puzzle.\n")
	}
}

func SudokuTest(b *SudokuB) {
	SudokuChecks()
	easyTwoSteps := "6.2.5.........3.4..........43...8....1....2........7..5..27...........81...6....."

	b.Solve(easyTwoSteps, true)
	fmt.Println("Passes test.")
}

/*
 * Sudoku solver that uses an astar search found in search.go.
 * Using http://www2.warwick.ac.uk/fac/sci/moac/people/students/peter_cock/python/sudoku/ as a source for puzzles.
 */

// Solve the problem of putting 15 tiles in order when you only have one blank space.
func Sudoku() {
	board := &SudokuB{}
	board.Init()

	SudokuTest(board)

	solveCommandLine := func() {
		r := myio.NewReader()
		T, _ := strconv.Atoi(r.Read())
		for i := 0; i < T; i++ {
			before := time.Now()
			board.Solve(r.Read(), true)
			after := time.Now()
			delta := after.Sub(before)
			fmt.Printf("Puzzle %d took %v\n", i, delta)
		}
		//PrintBoard(board.board)
		fmt.Println("End of sudoku program.")
	}
	solveCommandLine()
}
