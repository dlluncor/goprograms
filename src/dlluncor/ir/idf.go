// package ir. idf describes building an index counting
// up terms in all documents and such.
package ir

type indCounter struct {
  docs []*docMetadata
}

func (i *indCounter) Count() {

}

func BuildIndex() {
  c := &indCounter{
    docs: allDocs,
  }
  c.Count()
}
