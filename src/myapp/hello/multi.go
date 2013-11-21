package hello

import (
    "fmt"
    "errors"
    "html/template"
    "net/http"
    "strconv"
    "strings"

    "appengine"
    "appengine/datastore"
    "appengine/channel"
    //"appengine/user"
)

func InitMulti() {
    http.HandleFunc("/multi", main)
    http.HandleFunc("/move", move)
    http.HandleFunc("/opened", opened)
    http.HandleFunc("/getToken", getToken)
    http.HandleFunc("/startGame", startGame)
}

type Game struct {
    UserX  string
    UserO  string
    MoveX  bool
    Board  string
    Winner string
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

func defaultGame() *MyGame {
    return &MyGame{
        Users: make(DbMap),
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
        //http.Error(w, err.Error(), 500)
        return err
      }
    } else {
      // TODO(dlluncor):
      // Always write a new default game if need be when changing the schema
      // constantly.
    }
    if g == nil {
      c.Errorf("Game should never be null here.")
    }
    // Now write this user to the list of connect users to this table.
    user := r.FormValue("u")
    g.Users[user] = token
    if _, err := datastore.Put(c, k, g); err != nil {
      return err
    }
    return nil
  }, nil)

  if err != nil {
    c.Errorf("Error in db with connect to table. %v", err)
  }
  //c.Infof("Users in table: %v", g.Users)
  // Return the game state to the user.
  channel.SendJSON(c, token, g)
}

func startGame(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  c.Infof("Got a message from a client that they want to start a game.")

  // Need to keep track of all users associated with the games, then
  // we must broadcast an update to all of them.
  tableKey := r.FormValue("g")
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
    k := datastore.NewKey(c, "WrGame", tableKey, 0, nil)
    g := new(MyGame)
    if err := datastore.Get(c, k, g); err != nil {
      return err
    }
    return nil
  }, nil)

  /*
  myg := MyGame{
    Msg: "hello",
  }
  // Send the game state to both clients.
  for _, token := range tokens {
     err := channel.SendJSON(c, token, myg)
     if err != nil {
         c.Errorf("sending Game: %v", err)
     }
  }*/
  if err != nil {
    c.Errorf("Error in start game. %v", err)
  }
  return
}

// Second part.

func move(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get the user and their proposed move.
    pos, err := strconv.Atoi(r.FormValue("i"))
    if err != nil {
        http.Error(w, "Invalid move", http.StatusBadRequest)
        return
    }
    key := r.FormValue("gamekey")

    g := new(Game)
    err = datastore.RunInTransaction(c, func(c appengine.Context) error {
        // Retrieve the game from the Datastore.
        k := datastore.NewKey(c, "Game", key, 0, nil)
        if err := datastore.Get(c, k, g); err != nil {
            return err
        }

        // Make the move (mutating g).
        if !g.Move("", pos) {
            return errors.New("Invalid move")
        }

        // Update the Datastore.
        _, err := datastore.Put(c, k, g)
        return err
    }, nil)
    if err != nil {
        http.Error(w, "Couldn't make move", http.StatusInternalServerError)
        c.Errorf("move: %v", err)
        return
    }

    // Send the game state to both clients.
    for _, uID := range []string{g.UserX, g.UserO} {
        err := channel.SendJSON(c, uID+key, g)
        if err != nil {
            c.Errorf("sending Game: %v", err)
        }
    }
}

// Need to start a table and broadcast it to all users.

var mainTemplate = template.Must(template.ParseFiles("multi_main.html"))

func main(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    //user.Current(c) // assumes 'login: required' set in app.yaml
    key := r.FormValue("gamekey")
  
    id := "myid"
    newGame := key == ""
    if newGame {
        key = id
    }
    err := datastore.RunInTransaction(c, func(c appengine.Context) error {
        k := datastore.NewKey(c, "Game", key, 0, nil)
        g := new(Game)
        if newGame {
            // No game specified.
            // Create a new game and make this user the 'X' player.
            g.UserX = id
            g.MoveX = true
            g.Board = strings.Repeat(" ", 9)
        } else {
            // Game key specified, load it from the Datastore.
            if err := datastore.Get(c, k, g); err != nil {
                return err
            }
            if g.UserO != "" {
                // Both players already in game, skip the Put below.
                return nil
            }
            if g.UserX != id {
                // This game has no 'O' player.
                // Make the current user the 'O' player.
                g.UserO = id
            }
        }
        // Store the created or updated Game to the Datastore.
        _, err := datastore.Put(c, k, g)
        return err
    }, nil)
    if err != nil {
        http.Error(w, "Couldn't load Game", http.StatusInternalServerError)
        c.Errorf("setting up: %v", err)
        return
    }

    tok, err := channel.Create(c, id+key)
    if err != nil {
        http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
        c.Errorf("channel.Create: %v", err)
        return
    }

    err = mainTemplate.Execute(w, map[string]string{
        "token":    tok,
        "me":       id,
        "game_key": key,
    })
    if err != nil {
        c.Errorf("mainTemplate: %v", err)
    }
}

func (g *Game) Move(uID string, pos int) (ok bool) {
    // validate the move and update the board
    // (implementation omitted in this example)
    ok = true
    return
}