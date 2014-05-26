// Annoying with types!!!
package main

func toType(in interface{}, t1 reflect.Type, t2 reflect.Type) reflect.Value {
	out := reflect.MakeMap(reflect.MapOf(t1, t2))
	inVal := reflect.ValueOf(in)
	for _, k := range inVal.MapKeys() {
		//k1 := reflect.ValueOf(k)
		k1 := reflect.Value(k).Convert(t1)
		v1 := inVal.MapIndex(k1)
		vi1 := v1.Interface()
		switch vi1.(type) {
		case int:
			v2 := vi1.(int)
			out.SetMapIndex(k1, v2)
		default:
			panic("Update toType")
		}
	}
	return out
}
