# JSON Parser in Go

My implementaion of a JSON parser in Go. It returns exit code 0 for a valid JSON file, otherwise exits with code 1 and prints an error message with line & column number.

Example: for the following input file

```
{
  "key1": true,
  "key2": False,
  "key3": null,
  "key4": "value",
  "key5": 101
}
```

it would print

```
error: unexpected character 'F' at line 3, column 11
```

## Usage

```bash
json-parser FILE
```

## Build & Run

```bash
go build
```

```bash
./json-parser FILE
```

## Or just run

```bash
go run . FILE
```

### Why implement a JSON parser?

- It's fun.
- I was learning Go & it's better to learn it this way.
- I was curious on how these langauge parsers worked under the hood.
- And I found this [amazing website](https://codingchallenges.fyi/) by [John Crickett](https://www.linkedin.com/in/johncrickett/) where he publish challenges which are divided into steps. Challenge - [Write Your Own JSON Parser](https://codingchallenges.fyi/challenges/challenge-json-parser).

### Acknowledgements

I didn't know how to get started with this cuz I had no idea on how to implement a Lexer & Parser. [@eatonphil](https://github.com/eatonphil)'s [blog post](https://notes.eatonphil.com/writing-a-simple-json-parser.html) on the topic helped me a lot. Do check it out!
