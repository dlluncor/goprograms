package hello

// TODO(dlluncor): Need ability to create lounges too...
import (
    "appengine"
    "appengine/datastore"
    "net/http"
    "fmt"
    "strings"
)

type MyLounge struct {
  Name string
  Games []string
}

// Database for the entire lounge.
type changeLoungeFunc func(l *MyLounge) bool


func defaultLounge() *MyLounge {
  return &MyLounge{
    Name: "",
    Games: []string{},
  }
}

var loungeNames = []string{"Intermediate Lounge", "Beginner Lounge"}

func deleteLounges(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  hadError := false
  for _, loungeName := range loungeNames {
    err := datastore.RunInTransaction(c, func(c appengine.Context) error {
      k := datastore.NewKey(c, "WrLounge", loungeName, 0, nil)
      if err := datastore.Delete(c, k); err != nil {
        return err;
      }
      return nil
    }, nil)
    if err != nil {
      hadError = true
      c.Errorf("Error deleting a lounge: %v", err)
    }
  }
  if !hadError {
    fmt.Fprintf(w, "Deleted lounges: %v", loungeNames)
  }
}

func createLounge(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  queryMap := r.URL.Query()
  loungeName := queryMap.Get("l")
  gamesStr := queryMap.Get("g")
  games := strings.Split(gamesStr, ",")
  l := defaultLounge()
  l.Name = loungeName
  l.Games = games
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
      k := datastore.NewKey(c, "WrLounge", loungeName, 0, nil)
      if _, err := datastore.Put(c, k, l); err != nil {
        return err;
      }
      return nil
    }, nil)
  if err != nil {
    fmt.Fprintf(w, "Error creating a lounge: %v", err)
  } else {
    fmt.Fprintf(w, "Success in creating lounge: %v with games: %s", loungeName, games)
  }
}

func getLounges(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)

  lounges := []MyLounge{}
  for _, loungeName := range loungeNames {
    loungeChanger := func(l *MyLounge) bool {
      return false
    }
    l := ChangeLounge(c, loungeName, loungeChanger)
    lounges = append(lounges, *l)
  }
  sendJSON(w, lounges)
}

// TODO(dlluncor): Merge this with ChangeGame as they are the same except for the key
// and the type of the changeEntityFunc, and the defaultLounge() thing.
// Utility function for reading from and updating a game before then
// doing further processing.
func ChangeLounge(c appengine.Context, loungeId string, clf changeLoungeFunc) *MyLounge {
  // Store all tables as part of the game state and send
  // a "startTimer" response.
  loungeKey := loungeId
  l := defaultLounge()
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
    k := datastore.NewKey(c, "WrLounge", loungeKey, 0, nil)
    if err := datastore.Get(c, k, l); err != nil {
      return err;
    }

    // Perform special logic here to manipulate game.
    shouldRunUpdate := clf(l)
    if !shouldRunUpdate {
      // Sometimes we don't need to update the database so don't.
      return nil
    }

    if _, err := datastore.Put(c, k, l); err != nil {
      return err
    }
    return nil
  }, nil)

  if err != nil {
    c.Errorf("Err in db transaction %v", err)
  }
  return l
}