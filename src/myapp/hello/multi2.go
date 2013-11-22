package hello

import (
    "strconv"
    "net/http"

    "appengine"
    "appengine/datastore"
    "appengine/channel"
)

type changeGameFunc func(g *MyGame) bool

// When the first user starts the game and gives us tables that defines
// this entire game.
//
// Notifies all users to start their timers for round 1.
func sendTables(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  tableKey := r.FormValue("g")
  // Store all tables as part of the game state and send
  // a "startTimer" response.
  timerStarted := false
  gameChanger := func(g *MyGame) bool {
    timerStarted = g.HadState("timerStarted")
    if timerStarted {
      return false
    }
    g.AddState("timerStarted")
    tableStrKeys := []string{"table1", "table2", "table3", "table4"}
    for _, tableStrKey := range tableStrKeys {
      tableVal := r.FormValue(tableStrKey)
      g.SetTable(tableStrKey, tableVal)
    }
    return true
  }
  g := ChangeGame(c, tableKey, gameChanger)

  if timerStarted {
    return
  }

  resp := &Resp{
    Action: "startTimers",
    Payload: "",
  }

  // Send an update to everyone.
  for _, token := range g.GetUserTokens() {
    err := channel.SendJSON(c, token, resp)
    if err != nil {
      c.Errorf("Err with sendTables response: %v", err)
    }
  }
}

func getRoundInfo(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  tableKey := r.FormValue("g")
  round := r.FormValue("r")
  isRoundFetched := false
  gameChanger := func(g *MyGame) bool {
    val := "roundFetched" + round
    isRoundFetched = g.HadState(val)
    if isRoundFetched {
      return false
    }
    roundInt, _ := strconv.Atoi(round)
    g.CreateTableInfo(roundInt)
    g.AddState(val)
    return true
  }
  g := ChangeGame(c, tableKey, gameChanger)   // Just need to read the game.

  if isRoundFetched {
    // We can only fetch a round once...
    return
  }
  tableInfo := g.GetTableInfo()
  resp := &Resp{
    Action: "aboutToStartRound",
    Payload: tableInfo,
  }
  // Send table information to everyone in the room (need a solution
  // for when someone randomly jumps into the game).
  for _, token := range g.GetUserTokens() {
    err := channel.SendJSON(c, token, resp)
    if err != nil {
      c.Errorf("Err with sendTables response: %v", err)
    }
  }
}


// Utility function for reading from and updating a game before then
// doing further processing.
func ChangeGame(c appengine.Context, gameId string, cgf changeGameFunc) *MyGame {
  // Store all tables as part of the game state and send
  // a "startTimer" response.
  tableKey := gameId
  g := defaultGame()
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
    k := datastore.NewKey(c, "WrGame", tableKey, 0, nil)
    if err := datastore.Get(c, k, g); err != nil {
      return err;
    }

    // Perform special logic here to manipulate game.
    shouldRunUpdate := cgf(g)
    if !shouldRunUpdate {
      // Sometimes we don't need to update the database so don't.
      return nil
    }

    if _, err := datastore.Put(c, k, g); err != nil {
      return err
    }
    return nil
  }, nil)

  if err != nil {
    c.Errorf("Err in db transaction %v", err)
  }
  return g
}