package generator

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{input: "ID", expected: "id"},
		{input: "Name", expected: "name"},
		{input: "CreatedAt", expected: "created_at"},
		{input: "UpdatedAt", expected: "updated_at"},
		{input: "IsMember", expected: "is_member"},
	}

	for _, tt := range test {
		op := toSnakeCase(tt.input)
		if op != tt.expected {
			t.Fatalf("Not Equal. expected=%s. got=%s", tt.expected, op)
		}
	}
}
