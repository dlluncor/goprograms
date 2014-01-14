package hello

import (
    "fmt"
    "strconv"
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
  queryMap := r.URL.Query()
  startStr := queryMap.Get("s")
  if startStr == "" {
    startStr = "0"
  }
  endStr := queryMap.Get("e")
  if endStr == "" {
    endStr = "5"
  }
  start, _ := strconv.Atoi(startStr)
  end, _ := strconv.Atoi(endStr) 
  // table Id's I'm futzing with.
  tableKeys := []string{}
  for i := start; i < end; i++ {
    tableKeys = append(tableKeys, fmt.Sprintf("game%d", i))
  }

  opts := &datastore.TransactionOptions{
    XG: true,
  }
  // Clear all tables in the DB.
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
    c.Infof("Resetting the values for a table in the database.")
    for _, tableKey := range tableKeys {
        k := datastore.NewKey(c, "WrGame", tableKey, 0, nil)
        g := defaultGame()
        if err := datastore.Get(c, k, g); err != nil {
          return err
        }
        g2 := defaultGame()
        g2.Language = g.Language // Only preserve the language.
        if _, err := datastore.Put(c, k, g2); err != nil {
          return err
        }
    }
    return nil
  }, opts)
  if err != nil {
    msg := fmt.Sprintf("Error clearing database: %v", err)
    c.Errorf(msg)
    fmt.Fprintf(w, msg)
  } else {
    msg := fmt.Sprintf("Deleted keys: %v", tableKeys)
    fmt.Fprintf(w, msg)
  }
}


var tableTemplate = template.Must(template.ParseFiles("word_racer.html"))

func handleWrPage(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  queryMap := r.URL.Query()
  table := queryMap.Get("t")
  cookObj, err := r.Cookie("ww-user")
  if err != nil {
    c.Errorf("Error reading user cookie: %v", err)
    return
  }
  id := cookObj.Value
  c.Infof("User that just entered the page is: %v", id)
  if len(id) == 0 {
    c.Errorf("Could not find out who the user is!!!! Got as the user: %v", id)
    return
  }

  token, err := channel.Create(c, table+ "-" + id)
  err = tableTemplate.Execute(w, map[string]string{
    "userToken": token,
  })
  if err != nil {
      c.Errorf("tableTemplate: %v", err)
  }
}

// Creates a game and saves it to the database.
func createGame(c appengine.Context, tableKey string, g *MyGame) error {
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
    k := datastore.NewKey(c, "WrGame", tableKey, 0, nil)
    if _, err := datastore.Put(c, k, g); err != nil {
        return err
    }
    return nil
  }, nil)
  return err
}

// Might want to split game up into two pieces, info needed for real-time updates,
// and the big one needed to update the user of all of the interactions which
// have happened thus far we will see...

func opened(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)

  // Make sure the table is in the database here and provide the entire state
  // to the user of what is going on right now.
  // The user needs to know.
  // users and their points.
  // current puzzles associated with game.
  // current words found in the puzzle.
  g := defaultGame()
  tableKey := r.FormValue("g")
  c.Infof("Got a message from a client that they connected to a table: %v", tableKey)
  token := r.FormValue("t")
  userExists := false
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
    k := datastore.NewKey(c, "WrGame", tableKey, 0, nil)
    if err := datastore.Get(c, k, g); err != nil {
      c.Infof("We have never entered this game into the database. Should have one already!!!")
      return err
    }
    // Now write this user to the list of connect users to this table.
    user := r.FormValue("u")
    userExists = g.AddUserToken(user, token)
    if userExists {
      c.Infof("User %v already exists, replacing their token.", user)
    }
    if _, err := datastore.Put(c, k, g); err != nil {
      return err
    }
    // Update user count for lounge.
    return nil
  }, nil)

  /*
  if !userExists {
    // Update the lounge with the number of current players for this table.
    // when a new player was added.
    lounge := r.FormValue("l")
    loungeChanger := func(l *MyLounge) bool {
      // TODO(dlluncor): In the middle of something will need to resolve this later...
      return true
    }
    ChangeLounge(c, lounge, loungeChanger)
  }
  */

  if err != nil {
    c.Errorf("Error in db with connect to table. %v", err)
    return
  }

  g.SetNow()  // Let users who jump in randomly to figure out what time it is from
  // when the last round started.
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