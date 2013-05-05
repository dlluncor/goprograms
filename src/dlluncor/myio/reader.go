package myio

import(
  "bufio"
  "os"
  "strings"
)

// Create a reader interface so I test how someone reads lines.
type Reader interface {
  Read() string
}

type myReader struct {
  reader *bufio.Reader
}

func (r myReader) Read() string {
  return r.rawInput()
}

// Create a new reader which reads from standard in.
func NewReader() Reader {
  in :=  bufio.NewReader(os.Stdin)
  return myReader{in}
}

func (r myReader) rawInput() string {
  line, _ := r.reader.ReadString('\n')
  line = strings.Replace(line, "\n", "", -1)
  return line
}
