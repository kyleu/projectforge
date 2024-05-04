package util

import (
	"bytes"
	"io"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

func ToTOML(x any) string {
	return string(ToTOMLBytes(x))
}

func ToTOMLBytes(x any) []byte {
	b, err := toml.Marshal(x) //nolint:errchktoml // no chance of error
	if err != nil {
		if err.Error() == "Only a struct or map can be marshaled to TOML" {
			err = errors.Wrapf(err, "can't serialize object of type [%T] to TOML", x)
		}
		panic(err)
	}
	return b
}

func FromTOML(msg []byte, tgt any) error {
	return toml.Unmarshal(msg, tgt)
}

func FromTOMLString(msg []byte) (string, error) {
	var tgt string
	err := toml.Unmarshal(msg, &tgt)
	return tgt, err
}

func FromTOMLMap(msg []byte) (ValueMap, error) {
	var tgt ValueMap
	err := toml.Unmarshal(msg, &tgt)
	return tgt, err
}

func FromTOMLAny(msg []byte) (any, error) {
	var tgt any
	err := FromTOML(msg, &tgt)
	return tgt, err
}

func FromTOMLReader(r io.Reader, tgt any) error {
	return toml.NewDecoder(r).Decode(tgt)
}

func FromTOMLStrict(msg []byte, tgt any) error {
	dec := toml.NewDecoder(bytes.NewReader(msg))
	dec.Strict(true)
	return dec.Decode(tgt)
}
