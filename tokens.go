package main

import (
	"fmt"
	"strconv"
)

type TokenKind int

const (
	JsonBoolean TokenKind = iota
	JsonNull
	JsonNumber
	JsonString
	JsonSyntax
)

type Token struct {
	kind   TokenKind
	value  string
	lineNo int
	colNo  int
}

var JsonSyntaxChars = map[rune]struct{}{
	'{': {},
	'}': {},
	':': {},
	'[': {},
	']': {},
	',': {},
}

func GetTokenKind(kind TokenKind) string {
	switch kind {
	case JsonBoolean:
		return "boolean"
	case JsonNull:
		return "null"
	case JsonNumber:
		return "number"
	case JsonString:
		return "string"
	case JsonSyntax:
		return "json syntax"
	default:
		return "invalid"
	}
}

func ConvertTokenToType(token Token) (any, error) {
	var value any
	var err error

	switch token.kind {
	case JsonBoolean:
		value = token.value == "true"
	case JsonNull:
		value = nil
	case JsonString:
		value = token.value
	case JsonNumber:
		value, err = strconv.ParseFloat(token.value, 64)
	default:
		err = UnexpectedTokenError(token)
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func UnexpectedTokenError(token Token) error {
	return fmt.Errorf("unexpected %s token '%s' at line %d, column %d", GetTokenKind(token.kind), token.value, token.lineNo, token.colNo)
}
