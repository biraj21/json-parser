package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTokenKind(t *testing.T) {
	assert.Equal(t, "boolean", GetTokenKind(JsonBoolean))
	assert.Equal(t, "null", GetTokenKind(JsonNull))
	assert.Equal(t, "number", GetTokenKind(JsonNumber))
	assert.Equal(t, "string", GetTokenKind(JsonString))
	assert.Equal(t, "json syntax", GetTokenKind(JsonSyntax))
}
