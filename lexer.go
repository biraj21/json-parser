package main

import (
	"errors"
	"fmt"
	"regexp"
	"unicode"
)

var (
	jsonTrue  = []rune("true")
	jsonFalse = []rune("false")
	jsonNull  = []rune("null")
)

func Lex(s string) ([]Token, error) {
	tokens := []Token{}

	// for tracking line & column numbers to show to the user in case of error
	lineNo, colNo := 1, 1

	for runes := []rune(s); len(runes) > 0; {
		char := runes[0]

		// ignore white spaces
		if unicode.IsSpace(char) {
			colNo++

			if char == '\n' {
				lineNo++
				colNo = 1
			}

			runes = runes[1:]
			continue
		}

		token := Token{}
		var err error

		token, runes, err = lexString(runes, lineNo, colNo)
		if err != nil {
			return []Token{}, err
		} else if token != (Token{}) {
			tokens = append(tokens, token)
			colNo += len(token.value) + 2 // +2 for quotes
			continue
		}

		token, runes, err = lexBoolean(runes, lineNo, colNo)
		if err != nil {
			return []Token{}, err
		} else if token != (Token{}) {
			tokens = append(tokens, token)
			colNo += len(token.value)
			continue
		}

		token, runes, err = lexNull(runes, lineNo, colNo)
		if err != nil {
			return []Token{}, err
		} else if token != (Token{}) {
			tokens = append(tokens, token)
			colNo += len(token.value)
			continue
		}

		token, runes, err = lexNumber(runes, lineNo, colNo)
		if err != nil {
			return []Token{}, err
		} else if token != (Token{}) {
			tokens = append(tokens, token)
			colNo += len(token.value)
			continue
		}

		_, ok := JsonSyntaxChars[char]
		if ok {
			tokens = append(tokens, Token{JsonSyntax, string(char), lineNo, colNo})
			colNo++
			runes = runes[1:]
		} else {
			return tokens, fmt.Errorf("unexpected character '%s' at line %d, column %d", string(char), lineNo, colNo)
		}
	}

	return tokens, nil
}

func lexString(runes []rune, lineNo, colNo int) (Token, []rune, error) {
	if runes[0] != '"' {
		return Token{}, runes, nil
	}

	runes = runes[1:]

	escaped := false

	for i, char := range runes {
		if escaped {
			switch char {
			case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
				escaped = false
			default:
				return Token{}, runes, fmt.Errorf("invalid escaped character '\\%s' at line %d, col %d", string(char), lineNo, i+colNo)
			}
		} else if char == '\\' {
			escaped = true
		} else if char == '"' {
			return Token{JsonString, string(runes[:i]), lineNo, colNo}, runes[i+1:], nil
		}
	}

	return Token{}, runes, errors.New("expected end-of-string quote")
}

func lexBoolean(runes []rune, lineNo, colNo int) (Token, []rune, error) {
	if CompareRuneSlices(runes, jsonTrue, len(jsonTrue)) {
		return Token{JsonBoolean, string(runes[:len(jsonTrue)]), lineNo, colNo}, runes[len(jsonTrue):], nil
	}

	if CompareRuneSlices(runes, jsonFalse, len(jsonFalse)) {
		return Token{JsonBoolean, string(runes[:len(jsonFalse)]), lineNo, colNo}, runes[len(jsonFalse):], nil
	}

	return Token{}, runes, nil
}

func lexNull(runes []rune, lineNo, colNo int) (Token, []rune, error) {
	if CompareRuneSlices(runes, jsonNull, len(jsonNull)) {
		return Token{JsonNull, string(runes[:len(jsonNull)]), lineNo, colNo}, runes[len(jsonNull):], nil
	}

	return Token{}, runes, nil
}

func lexNumber(runes []rune, lineNo, colNo int) (Token, []rune, error) {
	if !unicode.IsDigit(runes[0]) && runes[0] != '-' {
		return Token{}, runes, nil
	}

	var endsAt int = len(runes) - 1
	for i, char := range runes {
		if !unicode.IsDigit(char) && char != 'e' && char != 'E' && char != '.' && char != '-' && char != '+' {
			endsAt = i - 1
			break
		}
	}

	tokenValue := string(runes[:endsAt+1])
	if !regexp.MustCompile(`^-?\d+(?:\.\d+)?(?:(e|E)-?\+?\d+)?$`).MatchString(tokenValue) {
		return Token{}, runes, fmt.Errorf("invalid number %s", tokenValue)
	}

	return Token{JsonNumber, tokenValue, lineNo, colNo}, runes[endsAt+1:], nil
}
