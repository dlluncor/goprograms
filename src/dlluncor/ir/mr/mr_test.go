package mr

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type tokenMapper struct {
}

func (m *tokenMapper) Map(i interface{}, emitFn EmitFn) {
	switch i.(type) {
	case string:
		s := i.(string)
		for _, w := range strings.Split(s, " ") {
			emitFn.Emit(Key(w), 1)
		}
	default:
		panic("Cant tokenize non string.")
	}
}

type sumReducer struct {
}

func (r *sumReducer) Reduce(k Key, vals []interface{}) reflect.Value {
	sum := int(0)
	for _, val := range vals {
		switch val.(type) {
		case int:
			sum += val.(int)
		default:
			panic("Cannot sum non-int.")
		}
	}
	return ToValue(sum)
}

var mrTests = []struct {
	mrSpec *Spec
	output interface{}
}{
	{
		&Spec{
			Input:   arr([]string{"hi there", "hi", "momma"}),
			Mapper:  &tokenMapper{},
			Reducer: &sumReducer{},
			Output:  Output{"map", intType},
		},
		map[Key]int{
			"hi":    2,
			"there": 1,
			"momma": 1,
		},
	},
}

// Convert output so the test can work.
func forTestToType(in *reflect.Value) interface{} {
	inter := (*in).Interface()
	switch inter.(type) {
	case map[Key]int:
		return inter.(map[Key]int)
	default:
		panic("Need to add this type manually.")
	}
}

func TestMr(t *testing.T) {
	for _, test := range mrTests {
		out := Run(test.mrSpec)
		expected := test.output
		actual := forTestToType(out)
		fmt.Printf("%v, %v\n", reflect.TypeOf(actual), reflect.TypeOf(expected))
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Output mismatch: Expected %v. Actual %v.", test.output, out)
		}
	}
}
