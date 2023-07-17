package main

import "log"

func Parse(tokens []Token) {
	if len(tokens) == 0 {
		log.Fatalf("empty JSON string")
	}
}
