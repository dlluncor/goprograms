package udacity

import(
  "dlluncor/container"
  "container/heap"
)

/* This file deals with doing an A* search using utilities found in
 * sudoku.go.
 */

type SNode struct {
  state *CellState
  f int32 // how many hops to get to this state.
  h int32 // how far am I from the goal.
}


/*
func (s *SFrontier) () {
  
}
*/

type SFrontier struct {
  queue *container.PriorityQueue
}

func (s *SFrontier) IsEmpty() bool {
  return s.queue.Len() == 0
}

func (s *SFrontier) RemoveChoice() interface{} {
  item := heap.Pop(s.queue).(*container.Item)
  node := item.Value.(*SNode)
  return node
}

func (s *SFrontier) Contains(node interface{}) bool {
  //TODO(dlluncor): did I see this state.
  return false
}

func (s *SFrontier) Add(inode interface{}) {
  node := inode.(*SNode)
  item := &container.Item{
    Value: node,
    Priority: node.f + node.h,
  }
  heap.Push(s.queue, item)
  // TODO(dlluncor): Add to list of explored states.
}


type SExplored struct {
  seen map[string]bool
}

/*
func (s *SExplored) () {
  
}
*/

func (s *SExplored) Add(node interface{}) {
  //TODO(dlluncor): did I see this state.
}

func (s *SExplored) Contains(node interface{}) bool {
  //TODO(dlluncor): did I see this state.
  return false
}

type SudokuSolver struct {
  frontier *SFrontier
  explored *SExplored
}

/*
func (s *SudokuSolver) () {
  
}
*/

func (s *SudokuSolver) IsGoal(inode interface{}) bool {
  node := inode.(*SNode)
  return node.state.IsSolved() 
}

func (s *SudokuSolver) NextActions(node interface{}) []interface{} {
  arr := make([]interface{}, 0)
  return arr
}

func (s *SudokuSolver) Init(s0 *CellState) {
  s.frontier = &SFrontier{
    queue: &container.PriorityQueue{},
  }
  heap.Init(s.frontier.queue)
  s.explored = &SExplored{
    seen: make(map[string]bool),
  }
  node0 := &SNode{
    f: 0,
    h: 0,
    state: s0,
  }
  s.frontier.Add(node0)
}