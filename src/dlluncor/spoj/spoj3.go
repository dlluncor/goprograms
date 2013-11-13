package spoj

// Solves Word Racer boards.
// Solves Scrabble as well.

import (
  "dlluncor/myio"
  "fmt"
  "strconv"
  "strings"
  "sort"
  "log"
  //"unicode/utf8"
)
const (
  EMPTY = "X"
)

func toArr(chars string) []string {
 s := []string{}
 for j := 0; j < len(chars); j++ {
   s = append(s, string(chars[j]))
  }
  return s 
}

type Node struct {
  val string
  nodes []*Node
  visited bool
}

func (n *Node) addEdge(edgeNode *Node) {
  n.nodes = append(n.nodes, edgeNode)
}

func createNode(val string) *Node {
  return &Node {
    val:val,
    nodes:[]*Node{},
    visited:false,
  } 
}

type wordValFunc func(word string) int

type graph struct {
  checker *Checker
  filters []filterFunc
  wordVal wordValFunc
  nodeMap map[string]*Node
  words map[string]bool
}

func (g *graph) getId(row int, col int, val string) string {
  return fmt.Sprintf("%d-%d-%s", row, col, val)
}

func (g *graph) getValFromId(id string) string {
  return strings.Split(id, "-")[2]
}

func (g *graph) getById(id string) *Node {
  // Constructs the node if one does not exist.
  node, ok := g.nodeMap[id]
  if !ok {
    val := g.getValFromId(id)
    anode := createNode(val)
    g.nodeMap[id] = anode
    return anode
  }
  return node
}

func (g *graph) Add(sourceId, toId string) {
  fromNode := g.getById(sourceId)
  toNode := g.getById(toId)
  if fromNode == toNode {
    return
  }
  fromNode.addEdge(toNode)
}

func (g *graph) RemapVal(val string) string {
  if val == "Q" {
    return "qu"
  }
  return val
}

func (g *graph) Connect(prev, cur, next []string, curLineInd int) {
  // All the arrays are the same length.
  n := len(prev)
  arrs := [][]string{prev, cur, next}
  for i := 0; i < n; i++ {
    // Connect to neighbors.
    neighs := []int{i}
    if (i - 1) >= 0 {
      neighs = append(neighs, i -1)
    }
    if (i + 1) < n {
      neighs = append(neighs, i + 1)
    }
    curChar := cur[i]
    curChar = g.RemapVal(curChar)
    //fmt.Printf("Cur char: %v. Neighs: %v\n", curChar, neighs)
    if curChar == EMPTY {
      continue
    }
    sourceId := g.getId(curLineInd, i, curChar)
    for _, neighCol := range neighs {
      // Now toId comes from the previous and next arrays, as well
      // as the current line.
      for index, arr := range arrs {
        neighVal := arr[neighCol]
        neighVal = g.RemapVal(neighVal)
        if neighVal == EMPTY {
          continue
        }
        toId := g.getId(curLineInd-1+index, neighCol, neighVal)
        g.Add(sourceId, toId)
      }
    }
  }
}

func emptyLine(length int) []string {
  l := []string{}
  for i := 0; i < length; i++ {
    l = append(l, EMPTY)
  }
  return l
}

func (g *graph) ConnectLines(lines []string) {
  numChars := len(lines[0])
  prev := emptyLine(numChars)
  num := len(lines)
  for i := 0; i < num; i++ {
    cur := toArr(lines[i])
    next := emptyLine(numChars)
    if i + 1 < num {
      next = toArr(lines[i+1])
    }
    //fmt.Printf("%v %v %v\n", prev, cur, next)
    // Now connect the items.
    g.Connect(prev, cur, next, i)
    prev = cur
  }
}

// Print lists the node value and its edges values.
func (g *graph) Print() {
  for id, node := range g.nodeMap {
    edges := []string{}
    for i := 0; i < len(node.nodes); i++ {
      edges = append(edges, node.nodes[i].val)
    }
    fmt.Printf("Id: %v. Edges: %v\n", id, edges)
  }
}

// A function that returns false if the word is invalid.
type filterFunc func(string) bool

func (g *graph) passesFilter(word string) bool {
  // Default filters.
  for _, filter := range g.filters {
    if !filter(word) {
      return false 
    }
  }
  if len(word) <= 2 {
    return false
  } 
  return true
}

func (g *graph) CheckAndAdd(word string) {
  if g.passesFilter(word) {
    g.words[word] = true
  }
}

// Explore looks at a graph of letters and determines the words that they
// can form.
// NOTE: copying lots of arrays can be slow and inefficient.
func (g *graph) Explore(n *Node, chars []string) {
  n.visited = true
  chars = append(chars, n.val)
  // Check if it works...
  word := strings.Join(chars, "")
  g.CheckAndAdd(word)
  isValidPrefix := g.checker.IsValidPrefix(word)
  if isValidPrefix {
    // Don't explore this word further if never occurs in any
    // possible words!!!
    for _, neighNode := range n.nodes {
      if neighNode.visited {
        continue
      }
      g.Explore(neighNode, chars)
    }
  }
  n.visited = false
}

func (g *graph) Solve() {
  // Does a search on the graph of words to find English words.
  // Register constraints here eventually for Scrabble.
  for _, node := range g.nodeMap {
    chars := []string{}
    g.Explore(node, chars)
  }
}

type wordSortFunc func(w1, w2 *string) bool

type wordSorter struct {
  words []string
  by wordSortFunc
}

// Len is part of sort.Interface.
func (s *wordSorter) Len() int {
	return len(s.words)
}

// Swap is part of sort.Interface.
func (s *wordSorter) Swap(i, j int) {
	s.words[i], s.words[j] = s.words[j], s.words[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *wordSorter) Less(i, j int) bool {
	return s.by(&s.words[i], &s.words[j])
}


// customSort sorts the words based on their matter of importance, most
// important going first.
func (g *graph) customSort(words []string) {
  by := func(w1, w2 *string) bool {
    return g.wordVal(*w1) < g.wordVal(*w2)
  }
  aWordSorter := &wordSorter{
    words:words,
    by:by,
  }
  sort.Sort(aWordSorter)
}

// Answer prints out the words we found. 
func (g *graph) Answer() []string {
  fmt.Println("Words:")
  words := []string{}
  for word, _ := range g.words {
    words = append(words, word)
  }
  g.customSort(words)
  return words
}

// Filters. 
func noShortWords(word string) bool {
  if len(word) <= 2 {
    return false
  }
  return true
}

// CreateWordRacerGraph constructs a graph given a representation
// of a board as an array of strings where each line
// corresponds to a set of characters spread out horizontally.
func CreateWordRacerGraph(fileName string, lines []string) *graph {
  checker := newChecker(fileName)
  filters := defaultFilters(checker)
  // Construct the graph from the lines.
  g := newGraph(checker, filters)
  g.wordVal = func(word string) int {
    return len(word)
  } 
  g.ConnectLines(lines)
  fmt.Println("Word racer graph in memory.")
  return g
}

func newChecker(fileName string) *Checker {
  // Initialize a checker object that validates whether words
  // are words or not.
  checker := &Checker{
    wordMap:make(map[string]bool),
    prefixes:make(prefixMapType),
    filePath:fileName,
  }
  checker.Initialize()
  fmt.Println("Checker initialized.")
  return checker
}

func defaultFilters(checker *Checker) []filterFunc {
  filters := []filterFunc{}
  filters = append(filters, noShortWords)
  isWord := func(word string) bool {
    return checker.IsWord(word)
  }
  filters = append(filters, isWord)
  return filters
}

func newGraph(checker *Checker, filters []filterFunc) *graph {
  return &graph{
    checker:checker,
    filters:filters,
    nodeMap: make(map[string]*Node),
    words: make(map[string]bool),
  }
}

func scrabbleGraph(g *graph, word string) {
  chars := toArr(word)
  for i := 0; i < len(chars); i++ {
    val := chars[i]
    nodeId := g.getId(0, i, val)
    for j := 0; j < len(chars); j++ {
      toVal := chars[j]
      toNodeId := g.getId(0, j, toVal)
      if i == j {
        continue
      }
      // Connect every letter with every other letter.
      g.Add(nodeId, toNodeId)
    }
  }
}

var SCRABBLE_VAL = map[string]int{
  "a": 1,
  "b": 3,
  "c": 3,
  "d": 2,
  "e": 1,
  "f": 4,
  "g": 2,
  "h": 4,
  "i": 1,
  "j": 8,
  "k": 5,
  "l": 1,
  "m": 3,
  "n": 1,
  "o": 1,
  "p": 3,
  "q": 10,
  "r": 1,
  "s": 1,
  "t": 1,
  "u": 1,
  "v": 4,
  "w": 4,
  "x": 8,
  "y": 4,
  "z": 10,
}

func scrabbleWordVal (word string) int {
  chars := toArr(word)
  val := 0
  for _, char := range chars {
    val += SCRABBLE_VAL[char]
  }
  return val
}

func CreateScrabbleGraph() *graph {
  fileName := "dlluncor/spoj/allWords.txt" 
  checker := newChecker(fileName)
  filters := defaultFilters(checker)
  g := newGraph(checker, filters)
  g.wordVal = func(word string) int {
    return scrabbleWordVal(word)
  }
  r := myio.NewReader()
  word := r.Read()
  scrabbleGraph(g, word)
  return g
}

type prefixMapType map[int]map[string]bool

type Checker struct {
  wordMap map[string]bool
  prefixes prefixMapType
  filePath string
}

func (c *Checker) Initialize() {
  lines, err := myio.ReadLines(c.filePath)
  if err != nil {
    log.Fatalf("Could not read all words text file! %v", err)
  }
  //fmt.Printf("%v", lines)
  for _, line := range lines {
    // Keep track of what are actual words.
    c.wordMap[line] = true
    // Also fill the valid prefix map.
    s := []string{}
    for i := 0; i < len(line); i++ {
      s = append(s, string(line[i]))
      prefixLen := len(s)
      prefix := strings.Join(s, "")
      prefixMap, ok := c.prefixes[prefixLen]
      if !ok {
        c.prefixes[prefixLen] = make(map[string]bool)
        prefixMap = c.prefixes[prefixLen]
      }
      prefixMap[prefix] = true
    }
  }
}

func (c *Checker) IsWord(word string) bool {
  _, ok := c.wordMap[word]
  return ok
}

func (c *Checker) IsValidPrefix(word string) bool {
  prefixMap, ok := c.prefixes[len(word)]
  if !ok {
    return false
  }
  _, ok = prefixMap[word]
  return ok
}


func linesFromFile() []string {
  // Read the board and save them as strings.
  r := myio.NewReader()
  num, _ := strconv.Atoi(r.Read())
  lines := []string{}
  for i := 0; i < num; i++ {
    chars := r.Read()
    fmt.Printf("%v\n", chars)
    lines = append(lines, chars)
  }
  return lines
} 

func WordRacer() {
  lines := linesFromFile()
  fileName := "dlluncor/spoj/allWords.txt" 
  g := CreateWordRacerGraph(fileName, lines)
  g.Solve()
  words := g.Answer()
  for _, word := range words {
    defer fmt.Printf("%v\n", word)
  }
}

// Takes in a board and returns a list of words that solves the
// puzzle.
func WordRacerFromServer(lines []string) []string {
  fileName := "allWords.txt" 
  g := CreateWordRacerGraph(fileName, lines)
  g.Solve()
  return g.Answer()
}

func Scrabble() {
  g := CreateScrabbleGraph()
  g.Solve()
  words := g.Answer()
  for _, word := range words {
    defer fmt.Printf("%v\n", word)
  }
}
