package main

import (
	"log"
)

func Parse(tokens []Token) {
	if len(tokens) == 0 {
		log.Fatalf("empty JSON string")
	}

	token := tokens[0]
	if token.kind != JsonSyntax {
		raiseTokenError(token)
	}

	if token.value == "{" {
		tokens = parseObject(tokens[1:])
	} else if token.value == "[" {
		tokens = parseArray(tokens[1:])
	} else {
		raiseTokenError(token)
	}

	if len(tokens) > 0 {
		raiseTokenError(tokens[0])
	}
}

func parseObject(tokens []Token) []Token {
	if len(tokens) == 0 {
		log.Fatalf("expected a key or an end-of-object brace '}'")
	}

	token := tokens[0]
	if token.kind == JsonSyntax && token.value == "}" {
		return tokens[1:]
	}

	const (
		checkKey   = iota
		checkColon = iota
		checkValue = iota
		checkEnd   = iota
	)

	var check = checkKey

	for len(tokens) > 0 {
		token = tokens[0]

		switch check {
		case checkKey:
			if token.kind != JsonString {
				raiseTokenError(token)
			}

			tokens = tokens[1:]
			check = checkColon
		case checkColon:
			if token.kind != JsonSyntax || (token.kind == JsonSyntax && token.value != ":") {
				raiseTokenError(token)
			}

			tokens = tokens[1:]
			check = checkValue
		case checkValue:
			if token.kind == JsonSyntax {
				if token.value == "{" {
					tokens = parseObject(tokens[1:])
				} else if token.value == "[" {
					tokens = parseArray(tokens[1:])
				} else {
					raiseTokenError(token)
				}
			} else {
				tokens = tokens[1:]
			}

			check = checkEnd
		case checkEnd:
			if token.kind != JsonSyntax {
				raiseTokenError(token)
			}

			if token.value == "," {
				tokens = tokens[1:]
			} else if token.value == "}" {
				return tokens[1:]
			} else {
				raiseTokenError(token)
			}

			check = checkKey
		}
	}

	switch check {
	case checkKey:
		log.Fatal("expected a key string")
	case checkColon:
		log.Fatal("expected a colon ':'")
	case checkValue:
		log.Fatal("expexted a value")
	default:
		log.Fatal("expected end-of-object brace '}'")
	}

	return []Token{}
}

func parseArray(tokens []Token) []Token {
	if len(tokens) == 0 {
		log.Fatalf("expected an element or an end-of-array bracket ']'")
	}

	token := tokens[0]
	if token.kind == JsonSyntax && token.value == "]" {
		return tokens[1:]
	}

	prevWasElement := false // to know if previous token was a valid array element or not

	for len(tokens) > 0 {
		token = tokens[0]

		if token.kind == JsonSyntax {
			if token.value == "[" && !prevWasElement {
				tokens = parseArray(tokens[1:])
				prevWasElement = true
			} else if token.value == "{" && !prevWasElement {
				tokens = parseObject(tokens[1:])
				prevWasElement = true
			} else if token.value == "]" && prevWasElement {
				return tokens[1:]
			} else if token.value == "," && prevWasElement {
				prevWasElement = false
				tokens = tokens[1:]
			} else {
				raiseTokenError(token)
			}
		} else if prevWasElement {
			raiseTokenError(token)
		} else {
			prevWasElement = true
			tokens = tokens[1:]
		}
	}

	log.Fatal("expected end-of-array bracket ']'")
	return []Token{}
}

func raiseTokenError(token Token) {
	log.Fatalf("unexpected %s token '%s' at line %d, column %d", GetTokenKind(token.kind), token.value, token.lineNo, token.colNo)
}
