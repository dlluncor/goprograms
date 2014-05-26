// package ir. idf describes building an index counting
// up terms in all documents and such.
package ir

import (
        "bytes"
        "encoding/gob"
	"fmt"
        "os"
        "log"
	"reflect"
	"sort"

	"dlluncor/ir/mappers"
	"dlluncor/ir/mr"
	"dlluncor/ir/types"
)

type indCounter struct {
	docs []*types.DocMetadata
}

var docInfoType = reflect.TypeOf(types.DocInfo{})
var tfType = reflect.TypeOf(types.TF{})

type sortT []*types.TF

func (s sortT) Len() int           { return len(s) }
func (s sortT) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortT) Less(i, j int) bool { return s[i].Num < s[j].Num }

// intermediate to write to disk
type indexInfo struct {
 TF map[mr.Key]types.TF
}

func (i *indCounter) Count() *indexInfo {
	// Map doc to words within doc.
	spec := &mr.Spec{
		Input:   docToInterface(i.docs),
		Mapper:  &mappers.DocMapper{},
		Reducer: &mappers.DocReducer{},
		Output:  mr.Output{"map", docInfoType},
	}
	out := (mr.Run(spec).Interface()).(map[mr.Key]types.DocInfo)
	for id, inf := range out {
		fmt.Printf("\n*******\nDoc: %v\n", id)
		for t, tInf := range inf.Terms {
			fmt.Printf("%v: %v, ", t, tInf)
		}
		fmt.Printf("\n")
	}

        var outTF = make(map[mr.Key]types.TF)
	{
		// Reduce word counts to idf scores.
		spec := &mr.Spec{
			Input:   docToInterface(i.docs),
			Mapper:  &mappers.TermMapper{},
			Reducer: &mappers.TermReducer{},
			Output:  mr.Output{"map", tfType},
		}
		outTF = (mr.Run(spec).Interface()).(map[mr.Key]types.TF)
		ts := []*types.TF{}
		for t, tf := range outTF {
			newT := &types.TF{
				Term: string(t),
				Num:  tf.Num,
			}
			ts = append(ts, newT)
		}
		sort.Sort(sortT(ts))
		fmt.Printf("\n\n*********Terms:\n")
		for _, t := range ts {
			fmt.Printf("%v\n", t)
		}
	}
  return &indexInfo{
    TF: outTF,
  }
}

func check(err error) {
  if err != nil {
    log.Fatalf("%v\n", err)
  }
}

// Write writes out the index to disk.
func (i *indCounter) Write(inf *indexInfo) {
  var network bytes.Buffer
  enc := gob.NewEncoder(&network)

  // Term metadata for QRewrite.
  check(enc.Encode(inf))
  f, err := os.Create("qReWrite.dat")
  check(err)
  defer f.Close()
  f.Write(network.Bytes())

  
  var in indexInfo
  dec := gob.NewDecoder(&network)
  check(dec.Decode(&in))
  fmt.Println("***------------******") 
  fmt.Printf("%v\n", in) 
}

func BuildIndex() {
	c := &indCounter{
		docs: allDocs,
	}
	inf := c.Count()
        c.Write(inf)
}
