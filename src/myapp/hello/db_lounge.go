package hello

import (
    "appengine"
    "appengine/datastore"
    "net/http"
    "fmt"
    "strings"
)

// Not in DB but constructed from DB calls.
type LoungeResp struct {
  Lounges []MyLounge
  GameInfo map[string]MyGame
}


type MyLounge struct {
  Name string
  Games []string
}

// Database for the entire lounge.
type changeLoungeFunc func(l *MyLounge) bool


// List of lounges stored in the DB.
type MyLounges struct {
  LoungeNames []string
}

func defaultLoungeResp() LoungeResp {
  return LoungeResp {
    Lounges: []MyLounge{},
    GameInfo: make(map[string]MyGame),
  }
}

func defaultLounge() *MyLounge {
  return &MyLounge{
    Name: "",
    Games: []string{},
  }
}

func defaultLounges() *MyLounges {
  return &MyLounges {
    LoungeNames: []string{},
  }
}

func setUpDb(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)

  // Setup a key which contains the list of loounge names.
  ls := defaultLounges()
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
      k := datastore.NewKey(c, "WrData", "loungeNames", 0, nil)
      if _, err := datastore.Put(c, k, ls); err != nil {
        return err;
      }
      return nil
  }, nil)
  if err != nil {
    c.Infof("Problem with setup db.")
  }
  fmt.Fprintf(w, "Successfully setup db.")
}

func getLoungeNamesDb(c appengine.Context) []string {
  ls := defaultLounges()
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
      k := datastore.NewKey(c, "WrData", "loungeNames", 0, nil)
      if err := datastore.Get(c, k, ls); err != nil {
        return err;
      }
      return nil
  }, nil)
  if err != nil {
    c.Infof("Problem with getting all the lounge names.")
  }
  return ls.LoungeNames
}

func addLoungeNameDb(c appengine.Context, loungeName string) error {
  ls := defaultLounges()
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
      k := datastore.NewKey(c, "WrData", "loungeNames", 0, nil)
      if err := datastore.Get(c, k, ls); err != nil {
        return err;
      }
      ls.LoungeNames = append(ls.LoungeNames, loungeName)
      if _ , err := datastore.Put(c, k, ls); err != nil {
        return err;
      }
      return nil
  }, nil)
  return err
}

func clearLoungeNamesDb(c appengine.Context) error {
  ls := defaultLounges()
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
      k := datastore.NewKey(c, "WrData", "loungeNames", 0, nil)
      if _ , err := datastore.Put(c, k, ls); err != nil {
        return err;
      }
      return nil
  }, nil)
  return err
}

func deleteLounges(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  hadError := false
  loungeNames := getLoungeNamesDb(c)
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
  clearLoungeNamesDb(c)
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
  addLoungeNameDb(c, loungeName)
  // Create the lounges.
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
  // Create the games as well.
  lang := queryMap.Get("lang")
  for _, tableKey := range games {
    g := defaultGame()
    g.Language = lang
    createGame(c, tableKey, g)
  }
}

// Returns all information about the lounges and their associated games.
func getLounges(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)

  resp := defaultLoungeResp()
  lounges := []MyLounge{}
  loungeNames := getLoungeNamesDb(c)
  for _, loungeName := range loungeNames {
    loungeChanger := func(l *MyLounge) bool {
      return false
    }
    l := ChangeLounge(c, loungeName, loungeChanger)
    gameChanger := func (g *MyGame) bool {
      return false
    }
    for _, tableName := range l.Games {
      g := ChangeGame(c, tableName, gameChanger)
      if g != nil {
        resp.GameInfo[tableName] = *g
      }
    }
    lounges = append(lounges, *l)
  }
  resp.Lounges = lounges
  sendJSON(w, resp)
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
