package ir

import(
  "strings"
  sc "dlluncor/ir/score"
)

type listener interface {
  Score(q *query, d *doc) score
}

type ascorer struct {
  listeners []listener
}

type score struct {
  weight float64
  value float64
}

func (a *ascorer) Score(q *query, d *doc) {
  scores := []score{}
  for _, lis := range a.listeners {
    scores = append(scores, lis.Score(q, d)) 
  }
  // Combine.
  val := 0.0
  for _, score := range scores {
    val += score.weight * score.value
  } 
  d.score = val 
}

type unigram struct {
  docField string
}

func (s *unigram) Score(q *query, d *doc) score {
  matches := nGramMatch(q.raw, d.data.GetField(s.docField), 1) 
  value := 0.0
  for i := 0; i < len(matches); i++ {
      value += 3.0
  }

  return score{
    weight: 0.2,
    value: value, 
  }
}

// finds nGram in array starting at position i of length len.
func nGram(arr []string, i, len int) string {
  // Create a word of length n
  word := ""
  for j := i; j < i + len; j++ {
    // Decide higher where to parameterize this if need be, e.g.,
    // lower-case, don't lower case, etc.
    word += "++" + strings.ToLower(arr[j])
  }
  return word
}

func nGramMatch(qText, docText string, n int) []string {
  qWords := sc.Tokenize(qText)
  docWords := sc.Tokenize(docText)
   
  docMap := make(map[string]bool)
  for i := 0; i < len(docWords) - n + 1; i++ {
    word := nGram(docWords, i, n)
    docMap[word] = true
  }
  
  // Find the hits.
  matches := []string{} 
  for i := 0; i < len(qWords) - n + 1; i++ {
    word := nGram(qWords, i, n)
    if docMap[word] {
      // Found a ngram match.
      matches = append(matches, word)
    }
  }
  return matches
}

type bigram struct {
  docField string
}

func (s *bigram) Score(q *query, d *doc) score {
  matches := nGramMatch(q.raw, d.data.GetField(s.docField), 2)
  value := 0.0
  for i := 0; i < len(matches); i++ {
    value += 4.0
  }

  return score{
    weight: 0.4,
    value: value, 
  }
}

var allListeners = []listener{
  // term.
  &unigram{"description"},
  &bigram{"description"},
  &unigram{"title"},
  &bigram{"title"},
}

func RegisterListeners(a *ascorer) {
  a.listeners = allListeners
}
