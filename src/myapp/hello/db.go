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
    defer close(c)
    for k, v := range m {
        c <- datastore.Property {
            Name: k,
            Value: v,
        }
    }
    return nil
}

// Token associated with this particular user.
type User struct {
  Token string
}

type MyGame struct {
   Users DbMap
}

func (g MyGame) Load(c <-chan datastore.Property) error {
  g.Users.Load(c)
  return nil
}

func (g MyGame) Save(c chan<- datastore.Property) error {
  g.Users.Save(c)
  return nil
}
