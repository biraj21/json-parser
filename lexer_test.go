package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexValidJson(t *testing.T) {
	const json = `{
  "key": "value",
  "key-n": 101,
  "key-o": {
    "inner key": "inner value"
  },
  "key-l": ["list value"]
}`

	expectedTokens := []Token{
		{JsonSyntax, "{", 1, 1},

		{JsonString, "key", 2, 3},
		{JsonSyntax, ":", 2, 8},
		{JsonString, "value", 2, 10},
		{JsonSyntax, ",", 2, 17},

		{JsonString, "key-n", 3, 3},
		{JsonSyntax, ":", 3, 10},
		{JsonNumber, "101", 3, 12},
		{JsonSyntax, ",", 3, 15},

		{JsonString, "key-o", 4, 3},
		{JsonSyntax, ":", 4, 10},
		{JsonSyntax, "{", 4, 12},

		{JsonString, "inner key", 5, 5},
		{JsonSyntax, ":", 5, 16},
		{JsonString, "inner value", 5, 18},

		{JsonSyntax, "}", 6, 3},
		{JsonSyntax, ",", 6, 4},

		{JsonString, "key-l", 7, 3},
		{JsonSyntax, ":", 7, 10},
		{JsonSyntax, "[", 7, 12},
		{JsonString, "list value", 7, 13},
		{JsonSyntax, "]", 7, 25},

		{JsonSyntax, "}", 8, 1},
	}

	tokens, err := Lex(json)
	assert.NoError(t, err)
	assert.Equal(t, expectedTokens, tokens)
}

func TestLexInvalidJson(t *testing.T) {
	const json = `{
  "key": "value",
  "key-n": 101,
  "key-o": {
    "inner key": "inner value"
  },
  "key-l": ['list value']
}`

	_, err := Lex(json)
	assert.EqualError(t, err, `unexpected character ''' at line 7, column 13`)
}
