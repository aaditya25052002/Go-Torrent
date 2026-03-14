package bencode

import (
	"fmt"
)

func ExtractInfoBytes(data []byte) ([]byte, error) {
	s := string(data)

	if len(s) == 0 || s[0] != 'd' {
		return nil, fmt.Errorf("invalid bencoded data")
	}

	cursor := 1 // skip leading 'd'
	for cursor < len(s) && s[cursor] != 'e' {
		keyAny, keyConsumed, err := Decode(s[cursor:])

		if err != nil {
			return nil, fmt.Errorf("error decoding bencoded data: %w", err)
		}

		keyStr, ok := keyAny.(string)
		if !ok {
			return nil, fmt.Errorf("expected string key, got %T", keyAny)
		}

		cursor += keyConsumed
		valueStart := cursor

		_, valueConsumed, err := Decode(s[cursor:])
		if err != nil {
			return nil, fmt.Errorf("error decoding bencoded data: %w", err)
		}

		if keyStr == "info" {
			return []byte(s[valueStart : valueStart+valueConsumed]), nil
		}

		cursor += valueConsumed
	}

	return nil, fmt.Errorf("info section not found")
}
