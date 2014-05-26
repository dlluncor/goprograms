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
        panic("Cannot non DocInfo")
    }
  }
  panic("Unreachable DocReducer Reduce")
}

// Per term information.
type TermMapper struct {
}

func allWords(d *types.DocMetadata) []string {
  ws := []string{}
  ws = append(ws, sc.Tokenize(d.Description)...)  // dont remove stop words
  ws = append(ws, sc.Tokenize(d.Title)...)  // dont remove stop words
  return ws
}

func toTF(w string, d *types.DocMetadata) types.TF {
  return types.TF{
    Num: 1,
  }
}

func (m *TermMapper) Map(i interface{}, emitFn mr.EmitFn) {
  switch i.(type) {
    case *types.DocMetadata:
      d := i.(*types.DocMetadata)
      words := allWords(d)
      for _, word := range words {
         emitFn.Emit(mr.Key(word), toTF(word, d))
      }
    default:
      panic("Cant tokenize non string.")
  }
}

type TermReducer struct {
}

// Reduce for each term will sum up all docs it was found in.
func (r *TermReducer) Reduce(k mr.Key, vals []interface{}) reflect.Value {
  t := types.TF{}
  for _, val := range vals {
    switch val.(type) {
      case types.TF:
      v := val.(types.TF)
      t.Num = t.Num + v.Num
      default:
        panic("Cannot reduce TermReducer.")
    }
  }
  return mr.ToValue(t)
}
