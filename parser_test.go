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
	tokens, err := Lex(validJson)
	assert.Equal(t, err, nil)

	err = Parse(tokens)
	assert.Equal(t, err, nil)
}
