//go:build tinygo
// +build tinygo

package json

import "errors"

func Marshal(v any) ([]byte, error) {
	return nil, errors.New("Not implemented")
}

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return nil, errors.New("Not implemented")
}

type RawMessage []byte

func Unmarshal(data []byte, v any) error {
	return errors.New("Not implemented")
}
