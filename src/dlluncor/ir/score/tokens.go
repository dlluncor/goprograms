package score

import (
	"strings"
)

func Tokenize(text string) []string {
	// TODO: Consider:
	//   - removing stopwords
	//   - is You and you the same?
	//   - is hello? and hello the same?
	//   - is Twilight and twilight the same?
	return strings.Split(text, " ")
}
