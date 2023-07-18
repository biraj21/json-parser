package main

import (
	"log"
	"regexp"
	"unicode"
)

var (
	jsonTrue  = []rune("true")
	jsonFalse = []rune("false")
	jsonNull  = []rune("null")
)

func Lex(s string) []Token {
	tokens := []Token{}

	// for tracking line & column numbers to show to the user in case of error
	lineNo, colNo := 1, 1

	runes := []rune(s)

	for len(runes) > 0 {
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

		var token Token
		var ok bool

		token, runes, ok = lexString(runes, lineNo, colNo)
		if ok {
			tokens = append(tokens, token)
			colNo += len(token.value) + 2 // +2 for quotes
			continue
		}

		token, runes, ok = lexBoolean(runes, lineNo, colNo)
		if ok {
			tokens = append(tokens, token)
			colNo += len(token.value)
			continue
		}

		token, runes, ok = lexNull(runes, lineNo, colNo)
		if ok {
			tokens = append(tokens, token)
			colNo += len(token.value)
			continue
		}

		token, runes, ok = lexNumber(runes, lineNo, colNo)
		if ok {
			tokens = append(tokens, token)
			colNo += len(token.value)
			continue
		}

		_, ok = JsonSyntaxChars[char]
		if ok {
			tokens = append(tokens, Token{JsonSyntax, string(char), lineNo, colNo})
			colNo++
			runes = runes[1:]
		} else {
			log.Fatalf("unexpected character '%s' at line %d, column %d\n", string(char), lineNo, colNo)
		}
	}

	return tokens
}

func lexString(runes []rune, lineNo, colNo int) (Token, []rune, bool) {
	if runes[0] != '"' {
		return Token{}, runes, false
	}

	runes = runes[1:]

	for i, char := range runes {
		if char == '"' {
			return Token{JsonString, string(runes[:i]), lineNo, colNo}, runes[i+1:], true
		}
	}

	log.Fatal("expected end-of-string quote")
	return Token{}, runes, true
}

func lexBoolean(runes []rune, lineNo, colNo int) (Token, []rune, bool) {
	if CompareRuneSlices(runes, jsonTrue, len(jsonTrue)) {
		return Token{JsonBoolean, string(runes[:len(jsonTrue)]), lineNo, colNo}, runes[len(jsonTrue):], true
	}

	if CompareRuneSlices(runes, jsonFalse, len(jsonFalse)) {
		return Token{JsonBoolean, string(runes[:len(jsonFalse)]), lineNo, colNo}, runes[len(jsonFalse):], true
	}

	return Token{}, runes, false
}

func lexNull(runes []rune, lineNo, colNo int) (Token, []rune, bool) {
	if CompareRuneSlices(runes, jsonNull, len(jsonNull)) {
		return Token{JsonNull, string(runes[:len(jsonNull)]), lineNo, colNo}, runes[len(jsonNull):], true
	}

	return Token{}, runes, false
}

func lexNumber(runes []rune, lineNo, colNo int) (Token, []rune, bool) {
	if !unicode.IsDigit(runes[0]) {
		return Token{}, runes, false
	}

	var endsAt int = len(runes) - 1
	for i, char := range runes {
		if !unicode.IsDigit(char) && char != 'e' && char != '.' {
			endsAt = i - 1
			break
		}
	}

	tokenValue := string(runes[:endsAt+1])
	if !regexp.MustCompile(`^\d+(?:\.\d+)?(?:e\d+)?$`).MatchString(tokenValue) {
		log.Fatalf("invalid number %s", tokenValue)
	}

	return Token{JsonNumber, tokenValue, lineNo, colNo}, runes[endsAt+1:], true
}
