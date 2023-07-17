package main

import (
	"fmt"
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

	tokens := Lex(string(contents))

	for _, token := range tokens {
		fmt.Println(token)
	}

	Parse(tokens)
}
