package main

func Deserialize(data string) (any, error) {
	var json any
	var err error

	tokens, err := Lex(data)
	if err != nil {
		return nil, err
	}

	json, err = Parse(tokens)
	if err != nil {
		return nil, err
	}

	return json, nil
}