package spoj

import (
  "dlluncor/spoj"
  "testing"
  "fmt"
)

var (
  index = 0
)

type FakeReader struct {
  lines []string
}

func (r FakeReader) Read() string {
  fmt.Printf("hi there %d", index)
  line := r.lines[index]
  index++
  return line
}

func TestBitmapper(t *testing.T) {
  bm := &spoj.Bitmapper{}
  r := FakeReader{
    []string{"3 4", "0001", "0010", "0110"},
  }
  bm.ReadInput(r)
  expectedAnswer := "3 2 1 0\n2 1 0 1\n1 0 0 1"
  actualAnswer := bm.Solve()
  if actualAnswer != expectedAnswer {
    t.Errorf("Bitmapper solve method does not work. %v != %v", expectedAnswer, actualAnswer)
  }
}
