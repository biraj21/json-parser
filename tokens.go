package main

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
		return ""
	}
}
