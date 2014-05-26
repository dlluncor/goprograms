package ir

import (
	"dlluncor/ir/types"
)

func docToInterface(in []*types.DocMetadata) []interface{} {
	out := make([]interface{}, len(in))
	for i, el := range in {
		out[i] = el
	}
	return out
}
