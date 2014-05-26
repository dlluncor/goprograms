// package ir. idf describes building an index counting
// up terms in all documents and such.
package ir

import(
  //"fmt"
  "reflect"

  "dlluncor/ir/mr"
  "dlluncor/ir/types"
  "dlluncor/ir/mappers"
)

type indCounter struct {
  docs []*docMetadata
}

var tInfoType = reflect.TypeOf(types.TInfo{})

func (i *indCounter) Count() {
  // Map doc to word counts.

  spec := &mr.Spec{
    Input: docToInterface(i.docs),
    Mapper: &mappers.DocMapper{},
    Reducer: &mappers.DocReducer{},
    Output: mr.Output{"map", tInfoType},
  }
   mr.Run(spec)
  // Reduce word counts to idf scores.
}

func BuildIndex() {
  c := &indCounter{
    docs: allDocs,
  }
  c.Count()
}
