package data

import (
	"encoding/json"
	"io"
)

func ToJson(v interface{}, w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}

func FromJson(v interface{}, r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(v)
}
