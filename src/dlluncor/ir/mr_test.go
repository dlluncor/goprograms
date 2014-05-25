package ir

import(
  "testing"
  "reflect"
  "strings"
  "fmt"
)

type Output interface {
  Equals(i interface{}) (bool, error)
}

/*
type wrapper struct {
  i interface{}
}

func (w *wrapper) Equals(i interface{}) (bool, error) {
  ok := reflect.DeepEqual(w.i, i)
  if !ok {
    return false, errors.New("Els not equal.")
  }
  return true, nil
}
*/

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

func (r *sumReducer) Reduce(k Key, vals []interface{}) interface{} {
  sum := int(0)
  for _, val := range vals {
    switch val.(type) {
      case int:
        sum += val.(int)
      default:
        panic("Cannot sum non-int.")
    }
  }
  return sum 
}

func arr(in []string) []interface{} {
  out := make([]interface{}, len(in))
  for i, el := range in {
    out[i] = el
  }
  return out
}

var mrTests = []struct{
  mrSpec *Spec 
  output interface{}
}{
  {
    &Spec{
      Input: arr([]string{"hi there", "hi", "momma"}),
      Mapper: &tokenMapper{},
      Reducer: &sumReducer{},
      Output: "map",
    },
    map[string]int{
      "hi": 2,
      "there": 1,
      "momma": 1,
    },
  },
}

// Annoying with types!!!
func toMine(in interface{}) map[string]int {
  out := map[string]int{}
  for k, v := range in.(map[Key]interface{}) {
    out[string(k)] = v.(int)
  }
  return out
}

func TestMr(t *testing.T) {
  for _, test := range mrTests {
    c := &mrCtrl{}
    c.Spec = test.mrSpec
    out := c.Run()
    expected := test.output
    actual := toMine(out)
    fmt.Printf("%v", reflect.TypeOf(out))
    if !reflect.DeepEqual(actual, expected) {
      t.Errorf("Output mismatch: Expected %v. Actual %v.", test.output, out)
    }
  }
}
