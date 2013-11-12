package hello

import (
    "fmt"
    "log"
    "net/http"
    "dlluncor/myio"
    "dlluncor/spoj"
    "strings"
    "strconv"
)

func init() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/wordracer_json", handlerWordRacer)
    http.HandleFunc("/wordracer", handlerWrPage)
}

var(
  puzzle = []string{
    "Qtrfen",
    "ppmite",
    "tiXXow",
    "asXXmt",
    "phsehw",
    "ijrlnm",
  }
)

func staticPage(fileName string) string {
  lines, err := myio.ReadLines(fileName)
  if err != nil {
    log.Fatalf("Could not find static page: %v", err)
  }
  return strings.Join(lines, "\n")
}

func handlerWrPage(w http.ResponseWriter, r *http.Request) {
  content := staticPage("static/word_racer.html")
  fmt.Fprintf(w, content)
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

// JSON handler.
func handlerWordRacer(w http.ResponseWriter, r *http.Request) {
    queryMap := r.URL.Query()
    content := queryMap.Get("board")
    length := queryMap.Get("length")
    lines := getLines(content, length)
    words := spoj.WordRacerFromServer(lines)
    // Should use JSON here.
    output := strings.Join(words, ",")
    fmt.Fprintf(w, output)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
}
