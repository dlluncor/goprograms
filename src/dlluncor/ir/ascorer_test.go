package ir

import (
	"testing"
)

type scoreMatch interface {
	Match(score float64) bool
}

type gtMatch struct {
	expected float64
}

func (m *gtMatch) Match(score float64) bool {
	return score >= m.expected
}

var scoreTests = []struct {
	descrip string
	doc     *doc
	q       *query
	s       listener
	m       scoreMatch
}{
	{
		"Title unigram match.",
		&doc{
			data: &docMetadata{
				title: "Temple Run",
			},
		},
		&query{
			raw: "temple",
		},
		&unigram{"title"},
		&gtMatch{0.0},
	},
}

func TestScorers(t *testing.T) {
	for _, test := range scoreTests {
		as := &ascorer{}
		as.listeners = append(as.listeners, test.s)
		as.Score(test.q, test.doc)
		if !test.m.Match(test.doc.score) {
			t.Errorf("Score mismatch: %v", test.descrip)
		}
	}
}
