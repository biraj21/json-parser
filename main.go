package main

import (
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	if len(os.Args) < 2 {
		log.Fatalf("usage: %s FILE...", os.Args[0])
	}

	contents, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	tokens, err := Lex(string(contents))
	if err != nil {
		log.Fatal(err)
	}

	_, err = Parse(tokens)
	if err != nil {
		log.Fatal(err)
	}
}
