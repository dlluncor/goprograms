package udacity

import(
  "dlluncor/container"
  "container/heap"
  "fmt"
  "log"
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
  seen map[string]bool
}

func (s *SFrontier) IsEmpty() bool {
  return s.queue.Len() == 0
}

func (s *SFrontier) RemoveChoice() interface{} {
  item := heap.Pop(s.queue).(*container.Item)
  node := item.Value.(*SNode)
  return node
}

func (s *SFrontier) Contains(inode interface{}) bool {
  //fmt.Println("!In frontier already!!")
  node := inode.(*SNode)
  _ , ok := s.seen[node.state.ToString()]
  return ok
}

func (s *SFrontier) Add(inode interface{}) {
  node := inode.(*SNode)
  item := &container.Item{
    Value: node,
    Priority: node.f + node.h,
  }
  heap.Push(s.queue, item)
  s.seen[node.state.ToString()] = true
}


type SExplored struct {
  seen map[string]bool
}

/*
func (s *SExplored) () {
  
}
*/

func (s *SExplored) Add(inode interface{}) {
  node := inode.(*SNode)
  s.seen[node.state.ToString()] = true
}

func (s *SExplored) Contains(inode interface{}) bool {
  //fmt.Println("!!Explored..")
  node := inode.(*SNode)
  _ , ok := s.seen[node.state.ToString()]
  return ok
}

type SudokuSolver struct {
  frontier *SFrontier
  explored *SExplored
  guess int32
}

/*
func (s *SudokuSolver) () {
  
}
*/

func (s *SudokuSolver) IsGoal(inode interface{}) bool {
  s.guess++
  node := inode.(*SNode)
  //node.state.Visualize()
  fmt.Printf("Number of unsolved squares: %d. Guess: %d\n", node.h, s.guess)
  return node.state.IsSolved() 
}

func SudokuHeuristic(s *CellState) int32{
  unsolved, possibs := s.NumUnsolved()
  if unsolved == 1 && possibs == 0 {
    // TODO(dlluncor): This state shouldn't be possible [1 ... 9]
    s.VisualizeAll()
    s.PrintAsInput()
    log.Fatalf("foobar state.")
  }
  return (unsolved * 10) + possibs
}

func (s *SudokuSolver) NextActions(inode interface{}) []interface{} {
  // TODO(dlluncor): If this is an unsolvable board, then continue no further. 

  arr := make([]interface{}, 0)
  node := inode.(*SNode)
  neighborStates := node.state.Neighbors()
  //fmt.Println("\nVisualizing neighbors...")
  for _, nState := range neighborStates {
    //nState.Visualize()
    nState.UpdatePossib()
    sNode := &SNode{
      state:nState,
      f:0,
      h:SudokuHeuristic(nState),
    }
    arr = append(arr, sNode)
  }
  return arr
}

func (s *SudokuSolver) Init(s0 *CellState) {
  s.frontier = &SFrontier{
    queue: &container.PriorityQueue{},
    seen: make(map[string]bool),
  }
  heap.Init(s.frontier.queue)
  s.explored = &SExplored{
    seen: make(map[string]bool),
  }
  s0.UpdatePossib()
  node0 := &SNode{
    f: 0,
    h: SudokuHeuristic(s0),
    state: s0,
  }
  s.frontier.Add(node0)
}