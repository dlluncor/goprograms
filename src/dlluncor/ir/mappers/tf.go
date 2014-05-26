package mappers

import(
  "dlluncor/ir/mr"

  "reflect"
  "strings"
)

type DocMapper struct {
}

func (m *DocMapper) Map(i interface{}, emitFn mr.EmitFn) {
  switch i.(type) {
    case *docMetadata:
      s := i.(string)
      for _, w := range strings.Split(s, " ") {
        emitFn.Emit(mr.Key(w), 1)
      }
    default:
      panic("Cant tokenize non string.")
  }
}

type DocReducer struct {
}

func (r *DocReducer) Reduce(k mr.Key, vals []interface{}) reflect.Value {
  sum := int(0)
  for _, val := range vals {
    switch val.(type) {
      case int:
        sum += val.(int)
      default:
        panic("Cannot sum non-int.")
    }
  }
  return mr.ToValue(sum) 
}
