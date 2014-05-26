// package ir. idf describes building an index counting
// up terms in all documents and such.
package ir

import(
  "fmt"
  "reflect"

  "dlluncor/ir/mr"
  "dlluncor/ir/types"
  "dlluncor/ir/mappers"
)

type indCounter struct {
  docs []*types.DocMetadata
}

var tInfoType = reflect.TypeOf(types.DocInfo{})

func (i *indCounter) Count() {
  // Map doc to word counts.

  spec := &mr.Spec{
    Input: docToInterface(i.docs),
    Mapper: &mappers.DocMapper{},
    Reducer: &mappers.DocReducer{},
    Output: mr.Output{"map", tInfoType},
  }
  out := (mr.Run(spec).Interface()).(map[mr.Key]types.DocInfo)
  for id, inf := range out {
    fmt.Printf("\n*******\nDoc: %v\n", id)
    for t, tInf := range inf.Terms {
      fmt.Printf("%v: %v, ", t, tInf)
    }
    fmt.Printf("\n")
  }
  // Reduce word counts to idf scores.
}

func BuildIndex() {
  c := &indCounter{
    docs: allDocs,
  }
  c.Count()
}
