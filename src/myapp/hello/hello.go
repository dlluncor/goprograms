package hello

import (
    "appengine"
    "fmt"
    "net/http"
    "dlluncor/myio"
    "dlluncor/spoj"
    "strings"
    "strconv"
)

func init() {
    InitMulti()
    // JSON handlers.
    http.HandleFunc("/wordracer_json", handlerWordRacer)
    http.HandleFunc("/getallwords", handleGetAllWords)

    // HTML pages.
    http.HandleFunc("/", handleSigninPage) // View for signing in.
    http.HandleFunc("/seeLounges", handleWrLoungePage)  // View for all lounges.
    http.HandleFunc("/enterTable", handleWrPage)  // View for a table.
}

type mywriter struct {
  lines []string
}

func (w *mywriter) Write(p []byte) (n int, err error) {
  w.lines = append(w.lines, string(p))
  return len(p), nil
}

func staticPage(fileName string) (string, error) {
  lines, err := myio.ReadLines(fileName)
  if err != nil {
    return "", err
    //log.Fatalf("Could not find static page: %v", err)
  }
  return strings.Join(lines, "\n"), nil
}

// Transforms the puzzle passed in to a valid puzzle.
func getLines(content, length string) []string{
  l, _ := strconv.Atoi(length)
  lines := []string{}
  for i := 0; i < len(content); {
    line := content[i:i+l]
    lines = append(lines, line)
    i += l
  }
  return lines
}

func handleStaticPage(w http.ResponseWriter, r *http.Request, page string) {
  c := appengine.NewContext(r)
  content, err := staticPage(page)
  if err != nil {
    c.Criticalf("Could not not serve static page: %v", err)
  }
  fmt.Fprintf(w, content)
}

// Handlers.

func handleSigninPage(w http.ResponseWriter, r *http.Request) {
  handleStaticPage(w, r, "wr_signin.html")
}

func handleWrLoungePage(w http.ResponseWriter, r *http.Request) {
  handleStaticPage(w, r, "wr_lounge.html")
}

// JSON handler for getting all solutions. Still needed!
func handlerWordRacer(w http.ResponseWriter, r *http.Request) {
    checker := spoj.NewChecker("allWords.txt")
    //c := appengine.NewContext(r)
    queryMap := r.URL.Query()
    content := queryMap.Get("board")
    length := queryMap.Get("length")
    lines := getLines(content, length)
    //words := lines
    words := spoj.WordRacerFromServer(checker, lines)
    // Should use JSON here.
    output := strings.Join(words, ",")
    fmt.Fprintf(w, output)
}

// content should have newlines, that's why we need length!
// Deprecated cannot use currently.
func solveForWords(content string, length int) string {
  checker := spoj.NewChecker("allWords.txt")
  lines := getLines(content, string(length))
  words := spoj.WordRacerFromServer(checker, lines)
  return strings.Join(words, ",")
}

// Useful notes.
// JS beautifier: http://jsbeautifier.org/ 
// Channels: https://developers.google.com/appengine/docs/go/channel/
 
func handleGetAllWords(w http.ResponseWriter, r *http.Request) {
    checker := spoj.NewChecker("allWords.txt")
    words := checker.AllWords()
    output := strings.Join(words, ",")
    fmt.Fprintf(w, output)
}
