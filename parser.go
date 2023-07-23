package main

import (
	"errors"
	"fmt"
)

func Parse(tokens []Token) (any, error) {
	if len(tokens) == 0 {
		return nil, errors.New("empty JSON string")
	}

	token := tokens[0]
	if token.kind != JsonSyntax {
		return nil, UnexpectedTokenError(token)
	}

	json := any(nil)

	var err error
	if token.value == "{" {
		tokens, json, err = parseObject(tokens[1:])
	} else if token.value == "[" {
		tokens, json, err = parseArray(tokens[1:])
	} else {
		return nil, UnexpectedTokenError(token)
	}

	if err != nil {
		return nil, err
	}

	if len(tokens) > 0 {
		return nil, UnexpectedTokenError(tokens[0])
	}

	return json, nil
}

func parseObject(tokens []Token) ([]Token, map[string]any, error) {
	if len(tokens) == 0 {
		return []Token{}, nil, errors.New("expected a key or an end-of-object brace '}'")
	}

	json := map[string]any{}

	token := tokens[0]
	if token.kind == JsonSyntax && token.value == "}" {
		return tokens[1:], json, nil
	}

	keys := map[string]struct{}{}

	const (
		checkKey   = iota
		checkColon = iota
		checkValue = iota
		checkEnd   = iota
	)
	var check = checkKey

	var currentKey string

	var err error
	for len(tokens) > 0 {
		token = tokens[0]

		switch check {
		case checkKey:
			if token.kind != JsonString {
				return []Token{}, nil, UnexpectedTokenError(token)
			}

			_, ok := keys[token.value]
			if ok {
				fmt.Printf("warning: duplicate object key '%s' at line %d, col %d\n", token.value, token.lineNo, token.colNo)
			}

			keys[token.value] = struct{}{}
			currentKey = token.value
			tokens = tokens[1:]
			check = checkColon
		case checkColon:
			if token.kind != JsonSyntax || (token.kind == JsonSyntax && token.value != ":") {
				return []Token{}, nil, UnexpectedTokenError(token)
			}

			tokens = tokens[1:]
			check = checkValue
		case checkValue:
			var value any
			if token.kind == JsonSyntax {
				if token.value == "{" {
					tokens, value, err = parseObject(tokens[1:])
					if err != nil {
						return []Token{}, nil, err
					}

					json[currentKey] = value
				} else if token.value == "[" {
					tokens, value, err = parseArray(tokens[1:])
					if err != nil {
						return []Token{}, nil, err
					}

					json[currentKey] = value
				} else {
					return []Token{}, nil, UnexpectedTokenError(token)
				}
			} else {
				value, err = ConvertTokenToType(token)
				if err != nil {
					return []Token{}, nil, err
				}

				json[currentKey] = value

				tokens = tokens[1:]

			}

			check = checkEnd
		case checkEnd:
			if token.kind != JsonSyntax {
				return []Token{}, nil, UnexpectedTokenError(token)
			}

			if token.value == "," {
				tokens = tokens[1:]
			} else if token.value == "}" {
				return tokens[1:], json, nil
			} else {
				return []Token{}, nil, UnexpectedTokenError(token)
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

	return []Token{}, nil, err
}

func parseArray(tokens []Token) ([]Token, []any, error) {
	if len(tokens) == 0 {
		return []Token{}, nil, errors.New("expected an element or an end-of-array bracket ']'")
	}

	json := []any{}

	token := tokens[0]
	if token.kind == JsonSyntax && token.value == "]" {
		return tokens[1:], json, nil
	}

	prevWasElement := false // to know if previous token was a valid array element or not

	var err error

	for len(tokens) > 0 {
		token = tokens[0]

		var value any

		if token.kind == JsonSyntax {
			if token.value == "[" && !prevWasElement {
				tokens, value, err = parseArray(tokens[1:])
				if err != nil {
					return []Token{}, nil, err
				}

				json = append(json, value)
				prevWasElement = true
			} else if token.value == "{" && !prevWasElement {
				tokens, value, err = parseObject(tokens[1:])
				if err != nil {
					return []Token{}, nil, err
				}

				json = append(json, value)
				prevWasElement = true
			} else if token.value == "]" && prevWasElement {
				return tokens[1:], json, nil
			} else if token.value == "," && prevWasElement {
				prevWasElement = false
				tokens = tokens[1:]
			} else {
				return []Token{}, nil, UnexpectedTokenError(token)
			}
		} else if prevWasElement {
			return []Token{}, nil, UnexpectedTokenError(token)
		} else {
			value, err = ConvertTokenToType(token)
			if err != nil {
				return []Token{}, nil, err
			}

			json = append(json, value)

			prevWasElement = true
			tokens = tokens[1:]
		}
	}

	return []Token{}, nil, errors.New("expected end-of-array bracket ']'")
}
