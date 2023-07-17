package main

import (
	"log"
	"regexp"
	"unicode"
)

type TokenKind int

const (
	tokenBoolean   TokenKind = iota
	typeJsonSyntax TokenKind = iota
	tokenString    TokenKind = iota
	tokenNull      TokenKind = iota
	tokenNumber    TokenKind = iota
)

type Token struct {
	kind  TokenKind
	value string
}

var (
	JsonTrue  = []rune("true")
	JsonFalse = []rune("false")
	JsonNull  = []rune("null")
)

var JsonSyntaxTokens = map[rune]struct{}{
	'{': {},
	'}': {},
	':': {},
	'[': {},
	']': {},
	',': {},
}

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

		token, runes, ok = lexString(runes)
		if ok {
			tokens = append(tokens, token)
			colNo += len(token.value) + 2 // +2 for quotes
			continue
		}

		token, runes, ok = lexBoolean(runes)
		if ok {
			tokens = append(tokens, token)
			colNo += len(token.value)
			continue
		}

		token, runes, ok = lexNull(runes)
		if ok {
			tokens = append(tokens, token)
			colNo += len(token.value)
			continue
		}

		token, runes, ok = lexNumber(runes)
		if ok {
			tokens = append(tokens, token)
			colNo += len(token.value)
			continue
		}

		_, ok = JsonSyntaxTokens[char]
		if ok {
			tokens = append(tokens, Token{typeJsonSyntax, string(char)})
			colNo++
			runes = runes[1:]
		} else {
			log.Fatalf("invalid character '%s' at line %d, column %d\n", string(char), lineNo, colNo)
		}
	}

	return tokens
}

func lexString(runes []rune) (Token, []rune, bool) {
	if runes[0] != '"' {
		return Token{}, runes, false
	}

	runes = runes[1:]

	for i, char := range runes {
		if char == '"' {
			return Token{tokenString, string(runes[:i])}, runes[i+1:], true
		}
	}

	log.Fatal("expected end-of-string quote")
	return Token{}, runes, true
}

func lexBoolean(runes []rune) (Token, []rune, bool) {
	if CompareRuneSlices(runes, JsonTrue, len(JsonTrue)) {
		return Token{tokenBoolean, string(runes[:len(JsonTrue)])}, runes[len(JsonTrue):], true
	}

	if CompareRuneSlices(runes, JsonFalse, len(JsonFalse)) {
		return Token{tokenBoolean, string(runes[:len(JsonFalse)])}, runes[len(JsonFalse):], true
	}

	return Token{}, runes, false
}

func lexNull(runes []rune) (Token, []rune, bool) {
	if CompareRuneSlices(runes, JsonNull, len(JsonNull)) {
		return Token{tokenNull, string(runes[:len(JsonNull)])}, runes[len(JsonNull):], true
	}
	return Token{}, runes, false
}

func lexNumber(runes []rune) (Token, []rune, bool) {
	if !unicode.IsDigit(runes[0]) {
		return Token{}, runes, false
	}

	var endsAt int = len(runes)
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

	return Token{tokenNumber, tokenValue}, runes[endsAt+1:], true
}
