package main

import (
  "dlluncor/myio"
  "testing"
)

type FakeReader struct {
}

func (r *FakeReader) ReadString(delim string) string {
  return ""
}

func TestBitmapper(t *testing.T) {
  bm := Bitmapper{}
  r := &FakeReader{}
  bm.ReadInput(r)
  t.Errorf("ssss")
}
