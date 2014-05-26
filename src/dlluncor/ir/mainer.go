package ir

import (
	"fmt"
	"os"
	"sort"

	"dlluncor/ir/types"
)

type query struct {
	raw string
	num int
}

type doc struct {
	score float64 /* returned by mustang after first pass. */
	name  string
	data  *types.DocMetadata
}

type mustang struct {
	index Index
}

func (m *mustang) Retrieve(q *query) []*doc {
	return m.index.Find(q)
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
	if len(os.Args) != 2 {
		fmt.Printf(`Usage: ./cmd "Angry birds"` + "\n")
		return
	}
	fmt.Printf("Hi main scorer.\n")
	rawQuery := os.Args[1]
	q := &query{
		raw: rawQuery,
		num: 10,
	}
	m := &mustang{}
	docs := m.Retrieve(q)
	for _, doc := range docs {
		as := &ascorer{}
		RegisterListeners(as)
		as.Score(q, doc)
	}
	s := &docSorter{
		docs: docs,
	}
	sort.Sort(s)
	// Print.
	fmt.Printf("***Results***\n")
	fmt.Printf("Pos\tName\tScore\n")
	for i, doc := range docs {
		fmt.Printf("%d. %v\t%.2f\n", i, doc.name, doc.score)
	}
}
