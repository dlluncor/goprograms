// MapReduce
// []V0 -> Map -> (K, V1)
// {Key: []V1} -> Reduce -> Output(V2)
package mr

import (
	"fmt"
	"reflect"
)

var (
	intType = reflect.TypeOf(1)
	keyType = reflect.TypeOf(Key(""))
)

func ToValue(in interface{}) reflect.Value {
	return reflect.ValueOf(in)
}

type Key string

type kv struct {
	k Key
	v interface{}
}

// Mapper and saving intermediate data.
type buffer struct {
	vals []kv
}

func (b *buffer) Emit(k Key, v interface{}) {
	b.vals = append(b.vals, kv{k, v})
}

type EmitFn interface {
	Emit(k Key, v interface{})
}

type mapper interface {
	Map(v interface{}, fn EmitFn)
}

// Reducer and shuffling.
type reducer interface {
	Reduce(k Key, vals []interface{}) reflect.Value
}

type Output struct {
	Kind string       // e.g., "map"
	V0   reflect.Type // e.g., Value produced by reducer.
}

// Full controller.

type Spec struct {
	Input   []interface{}
	Mapper  mapper // provide functions which return mappers and reducers?
	Reducer reducer
	Output  Output
}

func Run(spec *Spec) *reflect.Value {
	// Run mapper "in parallel".
	b := &buffer{} // where to store the shuffle data locally.
	for _, in := range spec.Input {
		mpr := spec.Mapper
		mpr.Map(in, b)
	}
	// Shuffle.
	shuffled := make(map[Key][]interface{})
	for _, kv := range b.vals {
		_, ok := shuffled[kv.k]
		if !ok {
			shuffled[kv.k] = make([]interface{}, 0)
		}
		vals := shuffled[kv.k]
		vals = append(vals, kv.v)
		shuffled[kv.k] = vals
	}

	// Reduce.
	output := reflect.MakeMap(reflect.MapOf(keyType, spec.Output.V0))
	for k, values := range shuffled {
		reducer := spec.Reducer
		out := reducer.Reduce(k, values)
		output.SetMapIndex(ToValue(k), out)
	}
	switch spec.Output.Kind {
	case "map":
		return &output
	default:
		panic(fmt.Sprintf("Unsupported output to MR: %v", spec.Output))
	}
	return nil
}
