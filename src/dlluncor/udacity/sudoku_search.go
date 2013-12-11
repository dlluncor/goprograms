package udacity

import(
)

/* This file deals with doing an A* search using utilities found in
 * sudoku.go.
 */

type SNode struct {
  state *CellState
  cost int32
  f int32 // how many hops to get to this state.
  h int32 // how far am I from the goal.
}


/*
func (s *SFrontier) () {
  
}
*/

type SFrontier struct {

}

func (s *SFrontier) IsEmpty() bool {
  return true
}

func (s *SFrontier) RemoveChoice() interface{} {
  return &SNode{}
}

func (s *SFrontier) Contains(node interface{}) bool {
  return true
}

func (s *SFrontier) Add(node interface{}) {

}


type SExplored struct {

}

/*
func (s *SExplored) () {
  
}
*/

func (s *SExplored) Add(node interface{}) {
  
}

func (s *SExplored) Contains(node interface{}) bool {
  return true
}

type SudokuSolver struct {
  frontier *SFrontier
  explored *SExplored
}

/*
func (s *SudokuSolver) () {
  
}
*/

func (s *SudokuSolver) IsGoal(node interface{}) bool {
  return false 
}

func (s *SudokuSolver) NextActions(node interface{}) []interface{} {
  arr := make([]interface{}, 0)
  return arr
}

func (s *SudokuSolver) Init(s0 *CellState) {
  s.frontier = &SFrontier{}
  s.explored = &SExplored{}
}