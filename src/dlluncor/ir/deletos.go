package ir

func docToInterface(in []*docMetadata) []interface{} {
  out := make([]interface{}, len(in))
  for i, el := range in {
    out[i] = el
  }
  return out
}
