package ir

import(
  "fmt"
  "sort"
)

type query struct {
  raw string
  num int
}

type doc struct {
  score int /* returned by mustang after first pass. */
  name string
  data *docMetadata
}

type mustang struct {
  index Index
}

func (m *mustang) Retrieve(q query) []*doc {
  return m.index.Find(q)
}

type ascorer struct {
  ind int
}

func (a *ascorer) Score(d *doc) {
  d.score = a.ind 
}

type docSorter struct {
  docs []*doc
}

func (d *docSorter) Len() int {
  return len(d.docs)
}

func (d *docSorter) Swap(i int, j int) {
  d.docs[i], d.docs[j] = d.docs[j], d.docs[i]
}

func (d *docSorter) Less(i int, j int) bool {
  return d.docs[i].score > d.docs[j].score
}

func MainScorer() {
  fmt.Printf("Hi main scorer.\n")
  q := query{
    raw: "lucene",
    num: 10,
  }
  m := &mustang{}
  docs := m.Retrieve(q)
  for i, doc := range docs {
    as := &ascorer{i}
    as.Score(doc)
  }
  s := &docSorter{
    docs: docs,
  }
  sort.Sort(s)
  // Print.
  fmt.Printf("***Results***\n")
  fmt.Printf("Pos\tName\tScore\n")
  for i, doc := range docs {
    fmt.Printf("%d. %v\t%v\n", i, doc.name, doc.score)
  }
}
