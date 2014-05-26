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

func (i *indCounter) Count() {
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

	{
		// Reduce word counts to idf scores.
		spec := &mr.Spec{
			Input:   docToInterface(i.docs),
			Mapper:  &mappers.TermMapper{},
			Reducer: &mappers.TermReducer{},
			Output:  mr.Output{"map", tfType},
		}
		out := (mr.Run(spec).Interface()).(map[mr.Key]types.TF)
		ts := []*types.TF{}
		for t, tf := range out {
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

}

func BuildIndex() {
	c := &indCounter{
		docs: allDocs,
	}
	c.Count()
}
