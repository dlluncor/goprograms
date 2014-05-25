// change to package mr
// MapReduce
// []V0 -> Map -> (K, V1)
// {Key: []V1} -> Reduce -> Output(V2)
package ir

import (
  "fmt"
  "reflect"
)

var (
  intType = reflect.TypeOf(1)
)

type Key string

type kv struct {
  k Key
  v interface{}
}

// Mapper and saving intermediate data.
type buffer struct {
  vals []kv
}

func (b *buffer) Emit(k Key, v interface{}) {
  b.vals = append(b.vals, kv{k, v})
}

type EmitFn interface {
  Emit(k Key, v interface{})  
}

type Mapper interface {
  Map(v interface{}, fn EmitFn)
}

// Reducer and shuffling.
type Reducer interface {
  Reduce(k Key, vals[]interface{}) interface{}
}

type Output struct {
  kind string // e.g., "map"
  v0 reflect.Type // e.g., Value produced by reducer.
}

// Full controller.

type Spec struct{
  Input []interface{}
  Mapper Mapper // provide functions which return mappers and reducers?
  Reducer Reducer
  Output Output 
}

type mrCtrl struct{
  Spec *Spec
}

func (m *mrCtrl) Run() interface{} {
  // Run mapper "in parallel".
  b := &buffer{}  // where to store the shuffle data locally.
  for _, in := range m.Spec.Input {
    mpr := m.Spec.Mapper
    mpr.Map(in, b)
  }
  // Shuffle.
  shuffled := make(map[Key][]interface{})
  for _, kv := range b.vals {
    _, ok := shuffled[kv.k] 
    if !ok {
      shuffled[kv.k] = make([]interface{}, 0)
    }
    vals := shuffled[kv.k]
    vals = append(vals, kv.v)
    shuffled[kv.k] = vals
  }
 
  // Reduce.
  output := make(map[Key]interface{})
  for k, values := range shuffled {
    reducer := m.Spec.Reducer
    out := reducer.Reduce(k, values)
    output[k] = out
  }
  switch m.Spec.Output.kind {
   case "map":
     return output
   default:
     panic(fmt.Sprintf("Unsupported output to MR: %v", m.Spec.Output)) 
  }
  return nil
}
