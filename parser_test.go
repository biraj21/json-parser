package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseValidJson(t *testing.T) {
	const validJson = `{
    "key": "value",
    "key-n": 101,
    "key-o": {
      "inner key": "inner value"
    },
    "key-l": ["list value"]
  }`

	expectedJson := map[string]any{
		"key":   "value",
		"key-n": 101.0,
		"key-o": map[string]any{"inner key": "inner value"},
		"key-l": []any{"list value"},
	}

	tokens, err := Lex(validJson)
	assert.Equal(t, err, nil)

	json, err := Parse(tokens)
	assert.NoError(t, err)
	assert.Equal(t, expectedJson, json)
}
