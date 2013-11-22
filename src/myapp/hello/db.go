package hello

import (
    "reflect"
    "appengine/datastore"
)
// My code for saving maps in a db.
type DbMap map[string]interface{}
 
func (m DbMap) Load(c <-chan datastore.Property) error {
    for p := range c {
        if p.Multiple {
            value := reflect.ValueOf(m[p.Name])
            if value.Kind() != reflect.Slice {
                m[p.Name] = []interface{}{p.Value}
            } else {
                m[p.Name] = append(m[p.Name].([]interface{}), p.Value)
            }
        } else {
            m[p.Name] = p.Value 
        }
    }
    return nil
}
 
func (m DbMap) Save(c chan<- datastore.Property) error {
    for k, v := range m {
        c <- datastore.Property {
            Name: k,
            Value: v,
        }
    }
    return nil
}

// This map will now contain everything b/c I don't know how to save
// or retrieve any other fields from this stupid struct besides the
// map values WTF!!!
// TODO(dlluncor): Figure out how to get a map type in here. For now,
// not worth it!!
type MyGame struct {
   Tokens []string
   States []string
   Tables []string
   // Current Table Info.
   CurTable string
   CurRound int // Starts off at 1, ends at 4. (not used kept by client!)
}

func defaultGame() *MyGame {
    g := &MyGame{
        Tokens: []string{},
        States: []string{},
        Tables: []string{},
    }
    g.AddState("notStarted")
    return g
}

/*
func (g MyGame) Load(c <-chan datastore.Property) error {
  err := g.Users.Load(c)
  if err != nil {
    return err
  }
  return nil
}

func (g MyGame) Save(c chan<- datastore.Property) error {
  defer close(c)
  err := g.Users.Save(c)
  if err != nil {
    return err
  }
  return nil
}
*/

// TODO(dlluncor): Pretty inefficient, but oh well! Go is fast I think...
func inArr(items []string, item string) bool {
  has := false
  for _, aItem := range items {
    if aItem == item {
        return true
    }
  }
  return has
}

func (g *MyGame) AddState(state string) {
  g.States = append(g.States, state)
}

func (g *MyGame) HadState(state string) bool {
  return inArr(g.States, state)
}

func (g *MyGame) AddUserToken(user, token string) {
  g.Tokens = append(g.Tokens, token)
}

// SetTable("table1", "dsdfsdfs\nXXXX\nDDFD")
// We'll use the index for now to generate which table to serve up.
func (g *MyGame) SetTable(tableRoundKey, tableVal string) {
  g.Tables = append(g.Tables, tableVal)
}

func (g *MyGame) GetUserTokens() []string {
  return g.Tokens
}

type TableInfo struct {
  Table string
  Round int
}

func (g *MyGame) CreateTableInfo(round int) {
  // Keep track of a new slate of data for this particular round.
  g.CurTable = g.Tables[round-1]
  g.CurRound = round
}

// Gets the table information for this round.
func (g *MyGame) GetTableInfo() *TableInfo {
  info := &TableInfo{
    Table: g.CurTable,
    Round: g.CurRound, 
  }
  return info
}
