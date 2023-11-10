package codec

import (
	"bytes"
	"encoding/gob"
)

type Codec interface {
	Register(v interface{})
	Encode(v any) ([]byte, error)
	Decode(data []byte, value any) error
}

type codec struct{}

func NewCodec() Codec {
	return &codec{}
}

func (c *codec) Register(v interface{}) {
	gob.Register(v)
}

func (c *codec) Encode(v any) ([]byte, error) {
	var out bytes.Buffer

	enc := gob.NewEncoder(&out)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (c *codec) Decode(data []byte, value any) error {
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(value)
	if err != nil {
		return err
	}
	return nil
}
