package util

import(
  "log"
  "bytes"
  "encoding/gob"
  "os"
  "io/ioutil"
)

func Check(err error) {
  if err != nil {
    log.Fatalf("%v\n", err)
  }
}

// EncodeToFile encodes an object and saves to a file.
func EncodeToFile(in interface{}, fname string) bytes.Buffer {
  var network bytes.Buffer
  enc := gob.NewEncoder(&network)

  // Term metadata for QRewrite.
  Check(enc.Encode(in))
  f, err := os.Create(fname)
  Check(err)
  defer f.Close()
  f.Write(network.Bytes())
  return network
}

// DecodeFile reads from a file and saves the data to out.
func DecodeFile(out interface{}, fname string) {
  byts, err := ioutil.ReadFile(fname)
  network := bytes.NewBuffer(byts)
  Check(err)
  dec := gob.NewDecoder(network)
  Check(dec.Decode(out))
}
