package bencode

import (
	"fmt"
	"strconv"
	"unicode"
)

func DecodeBencode(bencodedString string) (interface{}, int, error) {
	firstChar := bencodedString[0]

	if unicode.IsDigit(rune(firstChar)) {
		return decodeString(bencodedString)
	}

	if firstChar == 'i' {
		return decodeInteger(bencodedString)
	}

	if firstChar == 'l' {
		return decodeList(bencodedString)
	}

	return "", 0, fmt.Errorf("unsupported bencode type: %c", firstChar)
}

func decodeString(bencodedString string) (interface{}, int, error) {
	var firstColonIndex int

	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", 0, err
	}

	end := firstColonIndex + 1 + length
	return bencodedString[firstColonIndex+1 : end], end, nil
}

func decodeInteger(bencodedString string) (interface{}, int, error) {
	var firstEIndex int

	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == 'e' {
			firstEIndex = i
			break
		}
	}

	integer, err := strconv.Atoi(bencodedString[1:firstEIndex])
	if err != nil {
		return 0, 0, err
	}

	return integer, firstEIndex + 1, nil
}

func decodeList(bencodedString string) (interface{}, int, error) {
	list := []interface{}{}
	cursor := 1 // skip the 'l'

	for cursor < len(bencodedString) && bencodedString[cursor] != 'e' {
		value, consumed, err := DecodeBencode(bencodedString[cursor:])
		if err != nil {
			return nil, 0, err
		}
		list = append(list, value)
		cursor += consumed
	}

	return list, cursor + 1, nil
}
