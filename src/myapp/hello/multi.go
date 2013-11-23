package hello

import (
    "fmt"
    "html/template"
    "net/http"

    "appengine"
    "appengine/datastore"
    "appengine/channel"
    //"appengine/user"
)

func InitMulti() {
    http.HandleFunc("/multi", main)

    // Getting ready to start the game.
    http.HandleFunc("/opened", opened)
    http.HandleFunc("/getToken", getToken)
    http.HandleFunc("/startGame", startGame)
    http.HandleFunc("/_ah/channel/disconnected/", leaving)

    // Starting round 0.
    http.HandleFunc("/sendTables", sendTables)
    http.HandleFunc("/getRoundInfo", getRoundInfo)
    http.HandleFunc("/submitWord", submitWord)
    http.HandleFunc("/gameOver", gameOver)

    // Debug.
    http.HandleFunc("/clearAll", clearAll)
}

type Resp struct {
    Action string
    Payload interface{}
}

func clearAll(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  // table Id's I'm futzing with.
  tableKeys := []string{}
  for i := 0; i < 50; i++ {
    tableKeys = append(tableKeys, fmt.Sprintf("table%d", i))
  }

  // Clear all tables in the DB.
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
    c.Infof("Deleting all keys in database.")
    for _, tableKey := range tableKeys {
        k := datastore.NewKey(c, "WrGame", tableKey, 0, nil)
        // Delete each ones.
        datastore.Delete(c, k)
    }
    return nil
  }, nil)

  if err != nil {

  }
}


var tableTemplate = template.Must(template.ParseFiles("word_racer.html"))

func handleWrPage(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  queryMap := r.URL.Query()
  table := queryMap.Get("t")
  id := queryMap.Get("u")
  token, err := channel.Create(c, table+id)
  err = tableTemplate.Execute(w, map[string]string{
    "userToken": token,
  })
  if err != nil {
      c.Errorf("tableTemplate: %v", err)
  }
}

// Might have to store this stuff as a property list.
/*
func defaultProps() datastore.PropertyList {
  var plist datastore.PropertyList = make(datastore.PropertyList, 1)
  plist = append(plist, datastore.Property { "name", "Mat", false, false })
  return &plist
}
*/

// Might want to split game up into two pieces, info needed for real-time updates,
// and the big one needed to update the user of all of the interactions which
// have happened thus far we will see...

// Step 1. Get a token to open a channel when the user joins a table.
func getToken(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  queryMap := r.URL.Query()
  table := queryMap.Get("t")
  id := queryMap.Get("u")
  tok, err := channel.Create(c, table+id)

  if err != nil {
        c.Errorf("getToken error: %v", err)
  } else {
    fmt.Fprintf(w, tok)
  }
}

func opened(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  c.Infof("Got a message from a client that they connected to a table.")

  // Make sure the table is in the database here and provide the entire state
  // to the user of what is going on right now.
  // The user needs to know.
  // users and their points.
  // current puzzles associated with game.
  // current words found in the puzzle.
  g := defaultGame()
  tableKey := r.FormValue("g")
  token := r.FormValue("t")
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
    k := datastore.NewKey(c, "WrGame", tableKey, 0, nil)
    if err := datastore.Get(c, k, g); err != nil {
      c.Infof("We have never entered this game into the database so create a new one.")
      // Put the game in the database, this should basically happen only once.
      if _, err := datastore.Put(c, k, g); err != nil {
        return err
      }
    }
    // Now write this user to the list of connect users to this table.
    user := r.FormValue("u")
    g.AddUserToken(user, token)
    if _, err := datastore.Put(c, k, g); err != nil {
      return err
    }
    return nil
  }, nil)

  if err != nil {
    c.Errorf("Error in db with connect to table. %v", err)
  }
  resp := &Resp{
    Action: "join",
    Payload: g,
  }
  c.Infof("Notifying these tokens that a user entered: %v", g.GetUserTokens())
  // Let everyone know that they joined the game!
  for _, token := range g.GetUserTokens() {
    err := channel.SendJSON(c, token, resp)
    if err != nil {
      c.Errorf("sending Start game updates: %v", err)
    }
  }
}

func startGame(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  c.Infof("Got a message from a client that they want to start a game.")

  // Need to keep track of all users associated with the games, then
  // we must broadcast an update to all of them.
  tableKey := r.FormValue("g")
  isStarted := false
  g := defaultGame()
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
    k := datastore.NewKey(c, "WrGame", tableKey, 0, nil)
    if err := datastore.Get(c, k, g); err != nil {
      return err
    }
    isStarted = g.HadState("justStarted")
    if isStarted {
        return nil
    }
    g.AddState("justStarted")  // Set to true when a user first starts the game.
    if _, err := datastore.Put(c, k, g); err != nil {
      return err
    }
    return nil
  }, nil)
  
  if err != nil {
    c.Errorf("Error in start game db call. %v", err)
    return
  }

  if isStarted {
    // Do nothing if the game is already started.
    // TODO(dlluncor): Update the user's UI that the game is already started.
    c.Infof("Game has already started!!!")
    return
  }

  resp := &Resp{
    Action: "startGame",
    Payload: "start",
  }
  c.Infof("Tokens: %v", g.Tokens)
  for _, token := range g.GetUserTokens() {
    err := channel.SendJSON(c, token, resp)
    if err != nil {
      c.Errorf("sending Start game updates: %v", err)
    }
  }
  return
}

// Need to start a table and broadcast it to all users.

var mainTemplate = template.Must(template.ParseFiles("multi_main.html"))

func main(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    err := mainTemplate.Execute(w, map[string]string{})
    if err != nil {
        c.Errorf("mainTemplate: %v", err)
    }
}