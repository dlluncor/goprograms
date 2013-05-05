package spoj

import (
  "dlluncor/spoj"
  "testing"
)

type FakeReader struct {
  lines []string
  i int
}

func (r *FakeReader) Read() string {
  line := r.lines[r.i]
  r.i++
  return line
}

func TestBitmapper(t *testing.T) {
  bm := &spoj.Bitmapper{}
  r := &FakeReader{
    []string{"3 4", "0001", "0010", "0110"},
    0,
  }
  bm.ReadInput(r)
  expectedAnswer := "3 2 1 0\n2 1 0 1\n1 0 0 1"
  actualAnswer := bm.Solve()
  if actualAnswer != expectedAnswer {
    t.Errorf("Bitmapper solve method does not work. %v != %v", expectedAnswer, actualAnswer)
  }
}
