//go:build !tinygo
// +build !tinygo

package json

import (
	"encoding/json"
	"errors"
)

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

type RawMessage json.RawMessage

// MarshalJSON returns m as the JSON encoding of m.
func (m RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
