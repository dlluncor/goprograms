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
type MyGame struct {
   Users DbMap
}

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

var notUserKeys = map[string]bool{
  "isStarted": true,
}

func (g MyGame) IsStarted() bool {
  val, ok := g.Users["isStarted"]
  if !ok {
    return false
  }
  return val.(bool)
}

func (g MyGame) SetIsStarted(started bool) {
    g.Users["isStarted"] = started
}

func (g MyGame) AddUserToken(user, token string) {
    g.Users[user] = token
}

func (g MyGame) GetUserTokens() []string {
    users := []string{}
    for key, val := range g.Users {
       if _, ok := notUserKeys[key]; ok {
        continue
       }
       users = append(users, val.(string))
    }
    return users
}
