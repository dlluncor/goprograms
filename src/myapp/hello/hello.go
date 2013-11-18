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
    http.HandleFunc("/hello", handler)
    http.HandleFunc("/wordracer_json", handlerWordRacer)
    http.HandleFunc("/", handlerWrPage)
    http.HandleFunc("/sockets", handlerSocketPage)
}

type mywriter struct {
  lines []string
}

func (w *mywriter) Write(p []byte) (n int, err error) {
  w.lines = append(w.lines, string(p))
  return len(p), nil
}

/*
// Another way to read files if ever need be.
func ReadLines(path string) ([]string, error) {
  t, err := template.ParseFiles(path)
  mysaver := &mywriter{
    lines: []string{},
  }
  t.Execute(mysaver, "")
  if err != nil {
    return nil, err
  }
  return mysaver.lines, nil
}
*/

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

func handlerWrPage(w http.ResponseWriter, r *http.Request) {
  handleStaticPage(w, r, "word_racer.html")
}

func handlerSocketPage(w http.ResponseWriter, r *http.Request) {
  handleStaticPage(w, r, "websocket.html")
}

// JSON handler.
func handlerWordRacer(w http.ResponseWriter, r *http.Request) {
    //c := appengine.NewContext(r)
    queryMap := r.URL.Query()
    content := queryMap.Get("board")
    length := queryMap.Get("length")
    lines := getLines(content, length)
    //words := lines
    words := spoj.WordRacerFromServer(lines)
    // Should use JSON here.
    output := strings.Join(words, ",")
    fmt.Fprintf(w, output)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
}

// Useful notes.
// JS beautifier: http://jsbeautifier.org/ 
