package qrewrite

import(
  "fmt"

  "dlluncor/ir/types"
  "dlluncor/ir/util"
  "dlluncor/ir/mr" // silly just for mr.Key
  sc "dlluncor/ir/score"
)

type Rewriter struct{
  dfMap map[mr.Key]types.DF
}

func NewRewriter() *Rewriter{
  return &Rewriter{
     dfMap: make(map[mr.Key]types.DF),
  }
}

func (r *Rewriter) Init(dfFile string) {
  util.DecodeFile(&r.dfMap, dfFile) 
}

// Annotate adds the Term nodes and gives them annotations.
func (r *Rewriter) Annotate(q *types.Query) {
  // Add to the map of term to term information.
  terms := sc.Tokenize(q.Raw)
  nodes := []types.QNode{}
  for _, term := range terms {
    df, ok := r.dfMap[mr.Key(term)]
    if !ok {
      // How do we treat terms we have never seen before?
      df = types.DF{
        Num: -1,
      }
    }
    nodes = append(nodes, types.QNode{
      Token: term,
      DF: df,
    }) 
  }
  q.Nodes = nodes
  fmt.Printf("%v\n", q)
}
