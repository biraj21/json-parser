package main

import (
	"errors"
	"fmt"
)

func Parse(tokens []Token) error {
	if len(tokens) == 0 {
		return errors.New("empty JSON string")
	}

	token := tokens[0]
	if token.kind != JsonSyntax {
		return unexpectedTokenError(token)
	}

	var err error
	if token.value == "{" {
		tokens, err = parseObject(tokens[1:])
	} else if token.value == "[" {
		tokens, err = parseArray(tokens[1:])
	} else {
		return unexpectedTokenError(token)
	}

	if err != nil {
		return err
	}

	if len(tokens) > 0 {
		return unexpectedTokenError(tokens[0])
	}

	return nil
}

func parseObject(tokens []Token) ([]Token, error) {
	if len(tokens) == 0 {
		return []Token{}, errors.New("expected a key or an end-of-object brace '}'")
	}

	token := tokens[0]
	if token.kind == JsonSyntax && token.value == "}" {
		return tokens[1:], nil
	}

	keys := map[string]struct{}{}

	const (
		checkKey   = iota
		checkColon = iota
		checkValue = iota
		checkEnd   = iota
	)
	var check = checkKey

	var err error
	for len(tokens) > 0 {
		token = tokens[0]

		switch check {
		case checkKey:
			if token.kind != JsonString {
				return []Token{}, unexpectedTokenError(token)
			}

			_, ok := keys[token.value]
			if ok {
				fmt.Printf("warning: duplicate object key '%s' at line %d, col %d\n", token.value, token.lineNo, token.colNo)
			}

			keys[token.value] = struct{}{}

			tokens = tokens[1:]
			check = checkColon
		case checkColon:
			if token.kind != JsonSyntax || (token.kind == JsonSyntax && token.value != ":") {
				return []Token{}, unexpectedTokenError(token)
			}

			tokens = tokens[1:]
			check = checkValue
		case checkValue:
			if token.kind == JsonSyntax {
				if token.value == "{" {
					tokens, err = parseObject(tokens[1:])
					if err != nil {
						return []Token{}, err
					}
				} else if token.value == "[" {
					tokens, err = parseArray(tokens[1:])
					if err != nil {
						return []Token{}, err
					}
				} else {
					return []Token{}, unexpectedTokenError(token)
				}
			} else {
				tokens = tokens[1:]
			}

			check = checkEnd
		case checkEnd:
			if token.kind != JsonSyntax {
				return []Token{}, unexpectedTokenError(token)
			}

			if token.value == "," {
				tokens = tokens[1:]
			} else if token.value == "}" {
				return tokens[1:], nil
			} else {
				return []Token{}, unexpectedTokenError(token)
			}

			check = checkKey
		}
	}

	switch check {
	case checkKey:
		err = errors.New("expected a key string")
	case checkColon:
		err = errors.New("expected a colon ':'")
	case checkValue:
		err = errors.New("expexted a value")
	default:
		err = errors.New("expected end-of-object brace '}'")
	}

	return []Token{}, err
}

func parseArray(tokens []Token) ([]Token, error) {
	if len(tokens) == 0 {
		return []Token{}, errors.New("expected an element or an end-of-array bracket ']'")
	}

	token := tokens[0]
	if token.kind == JsonSyntax && token.value == "]" {
		return tokens[1:], nil
	}

	prevWasElement := false // to know if previous token was a valid array element or not

	var err error

	for len(tokens) > 0 {
		token = tokens[0]

		if token.kind == JsonSyntax {
			if token.value == "[" && !prevWasElement {
				tokens, err = parseArray(tokens[1:])
				if err != nil {
					return []Token{}, err
				}

				prevWasElement = true
			} else if token.value == "{" && !prevWasElement {
				tokens, err = parseObject(tokens[1:])
				if err != nil {
					return []Token{}, err
				}

				prevWasElement = true
			} else if token.value == "]" && prevWasElement {
				return tokens[1:], nil
			} else if token.value == "," && prevWasElement {
				prevWasElement = false
				tokens = tokens[1:]
			} else {
				return []Token{}, unexpectedTokenError(token)
			}
		} else if prevWasElement {
			return []Token{}, unexpectedTokenError(token)
		} else {
			prevWasElement = true
			tokens = tokens[1:]
		}
	}

	return []Token{}, errors.New("expected end-of-array bracket ']'")
}

func unexpectedTokenError(token Token) error {
	return fmt.Errorf("unexpected %s token '%s' at line %d, column %d", GetTokenKind(token.kind), token.value, token.lineNo, token.colNo)
}
