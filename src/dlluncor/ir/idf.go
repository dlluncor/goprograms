// package ir. idf describes building an index counting
// up terms in all documents and such.
package ir

import (
	"fmt"
	"reflect"
	"sort"

	"dlluncor/ir/mappers"
	"dlluncor/ir/mr"
	"dlluncor/ir/types"
        "dlluncor/ir/util"
)

type indCounter struct {
	docs []*types.DocMetadata
}

var docInfoType = reflect.TypeOf(types.DocInfo{})
var tfType = reflect.TypeOf(types.DF{})

type sortT []types.DF

func (s sortT) Len() int           { return len(s) }
func (s sortT) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortT) Less(i, j int) bool { return s[i].Num < s[j].Num }

// intermediate to write to disk
type indexInfo struct {
 DF map[mr.Key]types.DF
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

        var outDF = make(map[mr.Key]types.DF)
	{
		// Reduce word counts to idf scores.
		spec := &mr.Spec{
			Input:   docToInterface(i.docs),
			Mapper:  &mappers.TermMapper{},
			Reducer: &mappers.TermReducer{},
			Output:  mr.Output{"map", tfType},
		}
		outDF = (mr.Run(spec).Interface()).(map[mr.Key]types.DF)
		ts := []types.DF{}
		for _, tf := range outDF {
			ts = append(ts, tf) 
		}
		sort.Sort(sortT(ts))
		fmt.Printf("\n\n*********Terms:\n")
		for _, t := range ts {
			fmt.Printf("%v\n", t)
		}
	}
  return &indexInfo{
    DF: outDF,
  }
}

// Write writes out the index to disk.
func (i *indCounter) Write(inf *indexInfo) {
  util.EncodeToFile(inf.DF, types.DFFile) 
  
  var in map[mr.Key]types.DF
  util.DecodeFile(&in, types.DFFile)
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
