package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

const validJson = `{
  "key": "value",
  "key-n": 101,
  "key-o": {
    "inner key": "inner value"
  },
  "key-l": ["list value"]
}`

const invalidJson = `{
  "key": "value",
  "key-n": 101,
  "key-o": {
    "inner key": "inner value"
  },
  "key-l": ['list value']
}`

func TestParseValidJson(t *testing.T) {
	Parse(Lex(validJson))
}

func TestParseInvalidJson(t *testing.T) {
	// run the crashing code when FLAG is set
	if os.Getenv("FLAG") == "1" {
		Parse(Lex(invalidJson))
		return
	}

	cmd := exec.Command("go", "test")
	cmd.Env = append(os.Environ(), "FLAG=1")

	err := cmd.Run()

	assert.NotEqual(t, nil, err)
	if err != nil {
		assert.Equal(t, "exit status 1", err.Error())
	}
}
