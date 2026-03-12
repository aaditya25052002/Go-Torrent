package bencode

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"string", "5:hello", "hello"},
		{"integer", "i52e", 52},
		{"list", "li1ei2ei3ee", []interface{}{1, 2, 3}},
		{"nested list", "lli1ei2eei3ee", []interface{}{[]interface{}{1, 2}, 3}},
		{"dictionary", "d3:foo3:bare", map[string]interface{}{"foo": "bar"}},
		{"dict with int", "d5:helloi52ee", map[string]interface{}{"hello": 52}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decoded, _, err := Decode(test.input)
			if err != nil {
				t.Fatalf("error decoding %s: %v", test.input, err)
			}
			if !reflect.DeepEqual(decoded, test.expected) {
				t.Fatalf("expected %v, got %v", test.expected, decoded)
			}
		})
	}

}
