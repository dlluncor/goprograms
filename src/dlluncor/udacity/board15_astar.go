package udacity

import(
  "dlluncor/container"
  "container/heap"
  "fmt"
)

// Utilities for the board when doing the a-star search.

// The distance this board is from being the correct answer. Just count the 
// number of misplaced tiles.
func H1Dist(board string) int32 {
  numArr := ToBoard(board)
  distance := int32(0)
  for index, tileVal := range numArr {
    if tileVal - 1 != index {
      distance += 1
    }
  }
  return distance
}

// Map of tile number to the row and column that tile belongs to.
var TILE_TO_LOC = map[int][]int{
  1: {0, 0},
  2: {0, 1},
  3: {0, 2},
  4: {0, 3},
  5: {1, 0},
  6: {1, 1},
  7: {1, 2},
  8: {1, 3},
  9: {2, 0},
  10: {2, 1},
  11: {2, 2},
  12: {2, 3},
  13: {3, 0},
  14: {3, 1},
  15: {3, 2},
  16: {3, 3},
}

func Abs(num int) int{
  if num < 0 {
    return -num
  }
  return num
}

// The number of tiles which are misplaced, where the number of steps it takes to
// move that tile to the correct location is counted.
func H2Dist(board string) int32 {
  numArr := ToBoard(board)
  distance := int32(0)
  for row := 0; row < 4; row++ {
    for col := 0; col < 4; col++ {
      k := row * 4 + col
      tileVal := numArr[k]
      loc := TILE_TO_LOC[tileVal]
      correctRow := loc[0]
      correctCol := loc[1]
      distance += int32(Abs(correctRow - row) + Abs(correctCol - col))
    }
  }
  return distance
}

func Heuristic(board string) int32 {
  return H1Dist(board)
}


type BNode struct {
  parent *BNode // which state did I eminate from.
  state string // string representation of the board.
  cost int32
  f int32 // number of hops to get to this state.
  h int32 // how far am I from my goal.
}

type BExplored struct {
  boardMap map[string] bool // Whether this board state has already been explored.
}

func (e *BExplored) Add(node interface{}) {
  anode := node.(*BNode)
  e.boardMap[anode.state] = true
}

func (e *BExplored) Contains(node interface{}) bool {
  //fmt.Printf("Explored size: %v\n", len(e.boardMap))
  anode := node.(*BNode)
  _, ok := e.boardMap[anode.state]
  return ok 
}

// Fronteir of what has not been explored yet.
type BFronteir struct {
  boardMap map[string] *BNode // Map of boards to their state in the fronteir.
  queue *container.PriorityQueue // Keeps track of the next node with the least cost.
}

// Adds a board to the fronteir.
func (b *BFronteir) Add(inode interface{}) {
  // Calculate costs before adding to the fronteir.
  node := inode.(*BNode)
  if node.parent != nil {
    node.f = node.parent.f + 1  // We took one more hop, or step to get here.
    node.h = Heuristic(node.state)
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

func (b *BFronteir) Contains(inode interface{}) bool {
  node := inode.(*BNode)
  _, ok := b.boardMap[node.state]
  return ok
}

func(b *BFronteir) RemoveChoice() interface{} {
  //before := time.Now()
  // Pop the item with the lowest cost from the heap, fast!!!
  item := heap.Pop(b.queue).(*container.Item)
  lNode := b.boardMap[item.Value.(string)] // Node with the lowest cost.
  //fmt.Printf("Fronteir size: %v\n", len(b.boardMap))
  //fmt.Printf("Best choice. Cost: %v. H: %v\n", lNode.cost, lNode.h)
  //PrintBoard(lNode.state)
  delete(b.boardMap, lNode.state)
  //fmt.Printf("Time elapsed: %v\n", time.Since(before))
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
    h: Heuristic(board), 
  }
  node.cost = node.f + node.h
  bs.fronteir.Add(node)
  bs.explored = &BExplored{
    boardMap: make(map[string] bool),
  }
}

func (bs *BoardSolver) IsGoal(inode interface{}) bool {
  node := inode.(*BNode)
  return Heuristic(node.state) == 0
}

func (bs *BoardSolver) NextActions(inode interface{}) []interface{} {
  // Returns the list of BNodes to explore next.
  node := inode.(*BNode)
  nextNodes := make([]interface{}, 0)
  adjacentStates := FindAdjacents(node.state)
  //fmt.Println("Next states:")
  for _, state := range adjacentStates {
    nextNode := &BNode{
      parent: node,
      state: state,
    }
    //PrintBoard(state)
    nextNodes = append(nextNodes, nextNode)
  }
  // For these next actions, I need to see whether coming from this new
  // node actually results in a better f, or shorter path to me.
  for _, nextNode := range nextNodes {
    if bs.fronteir.Contains(nextNode) {
      // It's in the fronteir is there a better one?
      // I don't see this ever being true though...
      aNextNode := nextNode.(*BNode)
      if node.f + 1 < aNextNode.f {
        aNextNode.f = node.f + 1
        aNextNode.parent = node
        fmt.Printf("Found a shorter path to me!!")
      }
    }
  }

  return nextNodes
}