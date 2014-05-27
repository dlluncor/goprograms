package types

// - Indexing.

// Info about a doc when indexing.
type DocInfo struct {
	Terms map[string]*TInfo
}

// TInfo describes information about a token in a document.
type TInfo struct {
	Num int // "tf": num occurences in one doc
}

// Info about a term across many docs.

type DF struct {
	Num  int    // "df", num occurences in many documents.
	Term string // the term itself when used for sorting.
}

// DocMetadata == document in index.
type DocMetadata struct {
	Title       string
	Id          string
	Description string
        // Came from indexer MR.
        Inf *DocInfo
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

var (
  // Data files produced by index.
  base = "dlluncor/ir/data/"
  DFFile = base + "df.dat"  // Map of DF data keyed on term
  DocInfFile = base + "docInf.dat"  // Map of DocInfo data keyed on docid
)

// - Querying 
type QNode struct {
  Token string // e.g., "angry"
  DF DF 

  // Eventually can add children and have nested queries.
}

type Query struct {
	Raw string  // e.g., "angry birds"
	Num int  // e.g., number to return
        Nodes []QNode
}
