package ir

import(
  "testing"
  "reflect"
  "strings"
  "fmt"
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
  return val(sum) 
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
      Output: Output{"map", intType},
    },
    map[Key]int{
      "hi": 2,
      "there": 1,
      "momma": 1,
    },
  },
}

// Annoying with types!!!
func toMine(in interface{}) map[Key]int {
  out := map[Key]int{}
  for k, v := range in.(map[Key]interface{}) {
    out[k] = v.(int)
  }
  return out
}

func TestMr(t *testing.T) {
  for _, test := range mrTests {
    c := &mrCtrl{}
    c.Spec = test.mrSpec
    out := c.Run()
    expected := test.output
    actual := ((*out).Interface()).(map[Key]int)
    fmt.Printf("%v", reflect.TypeOf(out))
    if !reflect.DeepEqual(actual, expected) {
      t.Errorf("Output mismatch: Expected %v. Actual %v.", test.output, out)
    }
  }
}
