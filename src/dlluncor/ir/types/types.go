package types

// - Indexing.

// Info about a doc when indexing.
type DocInfo struct {
	Terms map[string]*TInfo
}

// TInfo describes information about a token in a document.
type TInfo struct {
	Num int // num occurences in one doc
}

// Info about a term across many docs.

type TF struct {
	Num  int    // num occurences in many documents.
	Term string // the term itself when used for sorting.
}

// DocMetadata == document in index.
type DocMetadata struct {
	Title       string
	Id          string
	Description string
}

func (m *DocMetadata) GetField(field string) string {
	if field == "description" {
		return m.Description
	}
	if field == "id" {
		return m.Id
	}
	if field == "title" {
		return m.Title
	}
	panic("Unrecognized field.")
}

// - Querying 

type Query struct {
	Raw string
	Num int
}
