package hello

import(
  "fmt"
  "log"
  "net/http"
  "encoding/json"
)

func sendJSON(w http.ResponseWriter, message interface{}) {
  b, err := json.Marshal(message)
  if err != nil {
      fmt.Println("error encoding the response a request")
      log.Fatal(err)
  }
  fmt.Fprintf(w, string(b))
}