package mappers

import(
  "dlluncor/ir/mr"
  "dlluncor/ir/types"
  sc "dlluncor/ir/score"

  "reflect"
  //"strings"
)

// Per doc information.
type DocMapper struct {
}

func toDocInfo(d *types.DocMetadata) types.DocInfo {
  words := sc.Tokenize(d.Description)

  terms := map[string]*types.TInfo{}
  for _, w := range words {
    _, ok := terms[w]
    if !ok {
      terms[w] = &types.TInfo{
        Num: 0,
      }
    }
    t := terms[w]
    t.Num = t.Num + 1 
  }

  return types.DocInfo{
    Terms: terms,
  }
}


func (m *DocMapper) Map(i interface{}, emitFn mr.EmitFn) {
  switch i.(type) {
    case *types.DocMetadata:
      d := i.(*types.DocMetadata)
      inf := toDocInfo(d)
      emitFn.Emit(mr.Key(d.Id), inf)
    default:
      panic("Cant tokenize non string.")
  }
}

type DocReducer struct {
}

// Reduce is an Identiy reducer.
func (r *DocReducer) Reduce(k mr.Key, vals []interface{}) reflect.Value {
  for _, val := range vals {
    switch val.(type) {
      case types.DocInfo:
        return mr.ToValue(val)
      default:
        panic("Cannot sum non-int.")
    }
  }
  panic("Unreachable")
}
