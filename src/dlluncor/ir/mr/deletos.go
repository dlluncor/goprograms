// delete this when I understand type conversion better.
package mr

func arr(in []string) []interface{} {
	out := make([]interface{}, len(in))
	for i, el := range in {
		out[i] = el
	}
	return out
}
