package ir

import(
  "fmt"
  "dlluncor/ir/types"
)

type tfIdf struct {
  method string
}

var(
  // okapi constants.
  k1 = 1.4 // [1.2, 2.0]
  b = 0.75
)

func okapi(tf float64, idf float64, numDocs float64, avgDl float64) float64 {
  n := tf * (k1 + 1)
  d := tf + k1 * (1 - b + (b * numDocs / avgDl))
  return idf * n / d
}

func (s *tfIdf) Score(q *types.Query, d *doc) score {
  sum := 0.0
  numDocs := float64(12)  // TODO(dlluncor): Unhardcode these.
  avgDl := 10.0
  fmt.Printf("Method: %v\n", s.method)
  //Num docs: %.1f. Avg doc len: %.1f\n", numDocs, avgDl)
  for _, node := range q.Nodes {
    // TF for doc?
    tfMap := d.data.Inf.Terms
    tfInf, ok := tfMap[node.Token]
    if !ok {
      // Never seen this term before no overlap.
      continue
    }
    idf := 1.0 / float64(node.DF.Num)
    tf := float64(tfInf.Num)
    v := 0.0
    switch s.method {
      case "okapi":
        v = okapi(tf, idf, numDocs, avgDl)
      case "tfidf":
        v = tf * idf
      default:
        panic("Unrecognized method tfidf Score")
    }
    sum += v
    fmt.Printf("Term overlap: %v. TF: %.1f. IDF: %.1f -> %.2f\n", node.Token, tf, idf, v)
  }
  return score{
    weight: 2.0,
    value: sum,
  }
}

var tfIdfListeners = []listener{
  &tfIdf{"okapi"},
  &tfIdf{"tfidf"},
}
