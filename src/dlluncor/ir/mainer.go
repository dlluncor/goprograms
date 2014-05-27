package ir

import (
	"fmt"
	"sort"

	"dlluncor/ir/types"
        "dlluncor/ir/qrewrite"
)

type doc struct {
	score float64 /* returned by mustang after first pass. */
	name  string
	data  *types.DocMetadata
}

type mustang struct {
	index *Index
}

func (m *mustang) Retrieve(q *types.Query) []*doc {
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

// MainScorer enters the Mustang scoring routine.
// pos that your arguments start at.
func MainScorer(pos int, args []string) {
        // Check params.
	if len(args) != 3 {
		fmt.Printf(`Usage: ./cmd "Angry birds"` + "\n")
		return
	}
	fmt.Printf("Hi main scorer.\n")
       
        // Init index.
        index := &Index{}
        index.Init(types.DocInfFile)

        // Construct query.
	rawQuery := args[pos + 1]
	q := &types.Query{
		Raw: rawQuery,
		Num: 10,
	}
        qe := qrewrite.NewRewriter()
        qe.Init(types.DFFile)
        qe.Annotate(q)
       
        // Fetch documents.
	m := &mustang{
          index: index,
        }
	docs := m.Retrieve(q)
        
        // Score top k.
	for _, doc := range docs {
		as := &ascorer{}
		RegisterListeners(as)
		as.Score(q, doc)
	}

	// Print results.
	s := &docSorter{
		docs: docs,
	}
	sort.Sort(s)
	fmt.Printf("***Results***\n")
	fmt.Printf("Pos\tName\tScore\n")
	for i, doc := range docs {
		fmt.Printf("%d. %v\t%.2f\n", i, doc.name, doc.score)
	}
}
