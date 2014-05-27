package qrewrite

import(
  "dlluncor/ir/types"
  sc "dlluncor/ir/score"
)

type Rewriter struct{
  map[string]types.TF
}

func (r *Rewriter) Init() {

}

// Annotate adds the Term nodes and gives them annotations.
func (r *Rewriter) Annotate(q *types.Query) {
  // Add to the map of term to term information.
  terms := sc.Tokenize(q.Raw)
  nodes := []types.QNode{}
  for _, term := range terms {
    tf, ok := r.tfMap[term]
    if !ok {
      // How do we treat terms we have never seen before?
      tf = TF{
        Num: -1,
      }
    }
    nodes = append(nodes, types.QNode{
      Token: term,
      TF: tf,
    }) 
  }
  q.Nodes = nodes
}
