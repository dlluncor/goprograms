package spoj

import (
	"dlluncor/myio"
	"fmt"
	"math"
	"strconv"
	"strings"
	"container/list"
	"unicode/utf8"
	"sort"
)

// GIRLSNBOYS
func getMaxDivers(G int, B int) int {
	var big = G
	var little = B
	if B > G {
		big = B
		little = G
	}
	if big == 0 && little == 0 {
		return 0
	}
	var left = big - 1 - little
	if left <= 0 {
		return 1
	}
	var beforeCeil = float64(left) / float64(little+1.0)
	var extra = int(math.Ceil(beforeCeil))
	return extra + 1
}

func GirlsBoys() {
	reader := myio.NewReader()
	for {
		line := reader.Read()
		if line == "0 0" {
			break
		}
		elements := strings.Split(line, " ")
		G, _ := strconv.Atoi(elements[0])
		B, _ := strconv.Atoi(elements[1])
		val := getMaxDivers(G, B)
		fmt.Printf("%d", val)
	}
}

// BITMAP
type Point struct {
  Row int
  Col int
}

type Info struct {
  Seen bool
  Dist int
}

type BitmapperI interface {
  ReadInput(in myio.Reader)
  Solve() string 
} 

type Bitmapper struct {
  valueMap map[Point] *Info
  q *list.List
  numRows int
  numCols int
}

func (b *Bitmapper) ReadInput(r myio.Reader) {
  b.q = list.New()
  b.valueMap = make(map[Point] *Info)
  line := r.Read()
  elements := strings.Split(line, " ")
  numRows, _ := strconv.Atoi(elements[0])
  numCols, _ := strconv.Atoi(elements[1])
  b.numRows = numRows
  b.numCols = numCols
  for i := 0; i < numRows; i++ {
    colLine := r.Read()
    colArr := strings.Split(colLine, "")
    for j := 0; j < numCols; j++ {
      val, _ := strconv.Atoi(colArr[j])
      point := Point{i, j}
      if val == 1 {
        b.valueMap[point] = &Info{true, 0}
	b.q.PushBack(point)
        //fmt.Printf("pushing a 1") 
      } else {
        b.valueMap[point] = &Info{false, -1}
      }
    }
  }
}

func (b *Bitmapper) Solve() string {
  for ; b.q.Len() > 0; {
    // Pop from the front of the queue.
    el := b.q.Front()
    p := el.Value.(Point) // a cast.
    b.q.Remove(el)
    curInfo := b.valueMap[p]
    newDistance := curInfo.Dist + 1
    
    // Create neighbors to look through.
    up := Point{p.Row-1, p.Col}
    down := Point{p.Row+1, p.Col}
    left := Point{p.Row, p.Col-1}
    right := Point{p.Row, p.Col+1}
    neighs := [4]Point{up, down, left, right}
    for _, neigh := range neighs {
      // Remove neighbors which are not in the map.
      if _, ok := b.valueMap[neigh]; !ok {
        continue
      }
      pointInfo, _ := b.valueMap[neigh]

      if pointInfo.Dist == -1 {
        pointInfo.Dist = newDistance 
      } else { 
        // For each neighbor, did I find a shorter path?
        if newDistance < pointInfo.Dist {
          pointInfo.Dist = newDistance
          pointInfo.Seen = false // Need to investigate this point again and put on queue
        }
      }

      // Add the neighbor to explore if we haven't looked at it yet.
      if !pointInfo.Seen {
        b.q.PushBack(neigh)
      }
    }
    curInfo.Seen = true // Ive explored this node now
  }

  // Print out the distances in a string.
  lines := make([]string, b.numRows)
  for i := 0; i < b.numRows; i++ {
    lineStrs := make([]string, b.numCols)
    for j := 0; j < b.numCols; j++ {
      p := Point{i, j}
      pointInfo, _ := b.valueMap[p]
      lineStrs[j] = strconv.Itoa(pointInfo.Dist)
    }
    lines[i] = strings.Join(lineStrs, " ")
  }
  answer := strings.Join(lines, "\n")
  return answer
}

func BitmapSolver(reader myio.Reader, bm BitmapperI) string {
  bm.ReadInput(reader)
  answer := bm.Solve()
  return answer
} 

// Spoj problems.

func Bitmap() {
  r := myio.NewReader()
  line := r.Read()
  T, _ := strconv.Atoi(line)
  for i := 0; i < T; i++ {
    bm := &Bitmapper{}
    answer := BitmapSolver(r, bm)
    fmt.Printf("%s", answer)
  }
}

// PARTY
// Woah this one is hard!!

type Partier struct {
  budget int
  numParties int
  r myio.Reader
}

func (p *Partier) Solve() (int, int) {
  // Create an array of maximum value for each dollar value
  // of your budget.
  costs := make([]int, p.budget + 1)
  maxFun := 0
  maxDollars := 0
  for i := 0;  i < p.numParties; i++ {
    line := p.r.Read()
    numStrs := strings.Split(line, " ")
    cost, _ := strconv.Atoi(numStrs[0])
    fun, _ := strconv.Atoi(numStrs[1])

    // nums[k] represents max fun at k + 1 dollars spent.
    maxBudgetToCheck := p.budget - cost 
    for k := maxBudgetToCheck; k >= 0; k-- {
      newBudget := cost + k
      curFun := costs[k]
      if k != 0 && curFun == 0 {
        // We can't spend this much anyway so ignore this.
        continue
      }
      possNewFun := curFun + fun
      if costs[newBudget] < possNewFun {
        costs[newBudget] = possNewFun
        // Save max fun right now?
        if possNewFun > maxFun {
          maxFun = possNewFun
          maxDollars= newBudget
        }
      }
    }
  }
  p.r.Read()
  return maxDollars, maxFun
}

func Party() {
  r := myio.NewReader()
  for {
    line := r.Read()
    if line == "0 0" {
      break
    }
    numStrs := strings.Split(line, " ")
    budget, _ := strconv.Atoi(numStrs[0])
    numParties, _ := strconv.Atoi(numStrs[1])
    partier := &Partier{budget, numParties, r}
    fees, fun := partier.Solve()
    fmt.Printf("%d %d\n", fees, fun)
  }
}

// PT07Z
// Another good one. Longest path in a tree.

// EDIST
// Edit distance. Need to figure that out.
// It is the Levenstein distance.

func myMin(a int, b int, c int) int {
  if a < b {
    if a < c {
      return a
    } else {
      return c
    }
  } else {
    if b < c {
      return b
    } else {
      return c
    }
  }
  fmt.Printf("Should never get here.")
  return 0
}

func editDist(input string, output string) int {
  inputStr := strings.Split(input, "")
  outputStr := strings.Split(output, "")
  m := utf8.RuneCountInString(input)
  n := utf8.RuneCountInString(output)

  // Allocate the 2D array.
  d := make([][]int, m + 1) // Allocate number of rows
  for i := range d {
    d[i] = make([]int, n + 1) // Number of columns.
  }

  // Source strings can become empty strings by dropping
  // each character.
  for i := 0; i < m + 1; i++ {
    d[i][0] = i
  }

  // To go from output to no source is the same.
  for j := 0; j < n + 1; j++ {
    d[0][j] = j
  } 

  for j := 1; j < n + 1; j++ {
    for i := 1; i < m + 1; i++ {
      // No cost to include this, so look at cost in our previous
      // iteration of the input and output.
      if inputStr[i-1] == outputStr[j-1] {
        d[i][j] = d[i-1][j-1]
      } else {
        // deletion, insertion, replacement
        minCost := myMin(d[i-1][j] + 1, d[i][j-1] + 1, d[i-1][j-1] + 1)
        d[i][j] = minCost
      }
    }
  }
  return d[m][n]
}

func EditDistance() {
  r := myio.NewReader()
  T, _ := strconv.Atoi(r.Read())
  for i := 0; i < T; i++ {
    word0 := r.Read()
    word1 := r.Read()
    fmt.Printf("%d\n", editDist(word0, word1))
  }
}

type Recmaner struct {
  inseq map[int]bool
  answers [] int
}

func (r *Recmaner) getRecaman(k int) int {
  return 2
}

// MRECAMAN
func Recaman() {
  r := myio.NewReader()
  inputs := []int{}
  for {
    line := r.Read()
    input, _ := strconv.Atoi(line)
    if input == -1 {
      break
    }
    inputs = append(inputs, input)
    //fmt.Printf("%d\n", getRecaman(input))
  }
  rec := &Recmaner{}
  for _, input := range inputs {
    ans := rec.getRecaman(input)
    fmt.Printf("%d\n", ans)
  }
}

var (
  START = "s"
  STOP = "e"
)

type GBaller struct {
  inds map[int]string
  keys []int
  entry int
}

func (b *GBaller) Reset(entries int) {
  b.inds = make(map[int]string)
  b.keys = make([]int, 2 * entries)
  b.entry = 0
}

func (b *GBaller) Process(start, stop int) {
  b.inds[start] = START
  b.inds[stop] = STOP
  b.keys[b.entry] = start
  b.entry++
  b.keys[b.entry] = stop
  b.entry++
}

func (b *GBaller) Answer() int {
  sort.Ints(b.keys)
  maxPeople := 0
  curPeople := 0
  for _, t := range b.keys {
    if b.inds[t] == START {
      curPeople++
      if curPeople > maxPeople {
        maxPeople = curPeople
      }
    } else {
      curPeople--
    }
  }
  return maxPeople 
}

func GreatBall() {
  r := myio.NewReader()
  b := &GBaller{}
  T, _ := strconv.Atoi(r.Read())
  for i := 0; i < T; i++ {
    N, _ := strconv.Atoi(r.Read())
    b.Reset(N)
    for j := 0; j < N; j++ {
      els := strings.Split(r.Read(), " ")
      num0, _ := strconv.Atoi(els[0])
      num1, _ := strconv.Atoi(els[1])
      b.Process(num0, num1)
    }
    fmt.Printf("%d\n", b.Answer()) 
  }
}

var (
 EPSILON = 0.00000001
)

func mysqrt(n float64) float64 {
  x := n
  err := x * x - n
  move := 0.0
  for ; err > EPSILON ; {
    move = (x * x - n) / (2 * x)
    x = x - move
    err = x * x - n
  }
  return x
}

func Sqrt() {
  r := myio.NewReader()
  T, _ := strconv.Atoi(r.Read())
  for i := 0; i < T; i++ {
    n, _ := strconv.Atoi(r.Read())
    fmt.Printf("%f\n", mysqrt(float64(n)))
  }
}

