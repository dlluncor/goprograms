package types

// TInfo describes information about a token.
type TInfo struct {

}

// DocMetadata == document in index.
type DocMetadata struct {
  title string
  id string
  description string
}

func (m *DocMetadata) GetField(field string) string {
  if field == "description" {
    return m.description
  }
  if field == "id" {
    return m.id
  }
  if field == "title" {
    return m.title
  }
  panic("Unrecognized field.")
}
