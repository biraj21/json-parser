package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareRuneSlices(t *testing.T) {
	assert.True(t, CompareRuneSlices([]rune("hello"), []rune("hello"), 5))

	assert.True(t, CompareRuneSlices([]rune("hello"), []rune("hey"), 2))

	assert.False(t, CompareRuneSlices([]rune("helper"), []rune("help"), 7))

	assert.False(t, CompareRuneSlices([]rune("helper"), []rune("help"), 5))
}
