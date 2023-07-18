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
	token := tokens[0]
	if token.kind == JsonSyntax && token.value == "}" {
		return tokens[1:]
	}

	for len(tokens) > 0 {
		// key
		token = tokens[0]
		if token.kind != JsonString {
			raiseTokenError(token)
		}

		tokens = tokens[1:]

		// colon
		token = tokens[0]
		if token.kind != JsonSyntax || (token.kind == JsonSyntax && token.value != ":") {
			raiseTokenError(token)
		}

		tokens = tokens[1:]

		// value
		token = tokens[0]
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

		// , or }
		token = tokens[0]
		if token.kind == JsonSyntax {
			if token.value == "," {
				tokens = tokens[1:]
			} else if token.value == "}" {
				return tokens[1:]
			}
		}
	}

	log.Fatal("expected end-of-object brace '}'")
	return []Token{}
}

func parseArray(tokens []Token) []Token {
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
