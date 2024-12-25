package gobase

import (
	"reflect"
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

type TestSchema struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type AnotherTestSchema struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func TestGenericSerializationDeserialization(t *testing.T) {
	t.Run("TestSchema Tests", func(t *testing.T) {
		tests := []struct {
			name        string
			input       any
			serialized  string
			expected    any
			shouldError bool
		}{
			{
				name:        "valid struct",
				input:       TestSchema{Name: "Test", Value: 42},
				serialized:  `{"name":"Test","value":42}`,
				expected:    TestSchema{Name: "Test", Value: 42},
				shouldError: false,
			},
			{
				name:        "empty struct",
				input:       TestSchema{},
				serialized:  `{"name":"","value":0}`,
				expected:    TestSchema{},
				shouldError: false,
			},
			{
				name:        "invalid JSON for deserialization",
				input:       TestSchema{},
				serialized:  `{"name":"Test", "value":}`,
				expected:    TestSchema{},
				shouldError: true,
			},
			{
				name:        "type mismatch during deserialization",
				input:       TestSchema{},
				serialized:  `{"id":"1","type":"test"}`,
				expected:    TestSchema{},
				shouldError: true,
			},
		}

		runSerializationTests[TestSchema](t, tests)
	})

	t.Run("AnotherTestSchema Tests", func(t *testing.T) {
		tests := []struct {
			name        string
			input       any
			serialized  string
			expected    any
			shouldError bool
		}{
			{
				name:        "valid struct",
				input:       AnotherTestSchema{ID: "123", Type: "test"},
				serialized:  `{"id":"123","type":"test"}`,
				expected:    AnotherTestSchema{ID: "123", Type: "test"},
				shouldError: false,
			},
			{
				name:        "empty struct",
				input:       AnotherTestSchema{},
				serialized:  `{"id":"","type":""}`,
				expected:    AnotherTestSchema{},
				shouldError: false,
			},
			{
				name:        "invalid JSON for deserialization",
				input:       AnotherTestSchema{},
				serialized:  `{"id":123,"type":}`,
				expected:    AnotherTestSchema{},
				shouldError: true,
			},
			{
				name:        "type mismatch during deserialization",
				input:       AnotherTestSchema{},
				serialized:  `{"name":"Test","value":42}`,
				expected:    AnotherTestSchema{},
				shouldError: true,
			},
		}

		runSerializationTests[AnotherTestSchema](t, tests)
	})
}

// Helper function to run serialization tests for any schema type
func runSerializationTests[T any](t *testing.T, tests []struct {
	name        string
	input       any
	serialized  string
	expected    any
	shouldError bool
},
) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialization test
			if tt.input != "" {
				serialized, err := serializeStruct(tt.input)
				if err != nil && !tt.shouldError {
					t.Fatalf("Unexpected error during serialization: %v", err)
				}
				if !tt.shouldError && serialized != tt.serialized {
					t.Errorf("Expected serialized %s, but got %s", tt.serialized, serialized)
				}
			}

			// Deserialization test
			deserialized, err := deserializeStruct[T](tt.serialized)
			if tt.shouldError {
				if err == nil {
					t.Error("Expected error during deserialization, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error during deserialization: %v", err)
				}

				expectedStruct, ok := tt.expected.(T)
				if !ok {
					t.Fatal("Expected value is not of the correct type")
				}

				// Use reflection to compare the values
				if !reflect.DeepEqual(deserialized, expectedStruct) {
					t.Errorf("Expected deserialized %+v, but got %+v", expectedStruct, deserialized)
				}
			}
		})
	}
}

func TestSerializeChangeLog(t *testing.T) {
	// Define a populated ChangeLog
	populatedChangeLog := ChangeLog{
		Creations: []Create{
			{
				CreationType: "INSERT",
				ON:           "row",
				TableName:    "users",
				CreationData: `{"id":1,"name":"John"}`,
			},
		},
		Updates: []Update{
			{
				UpdateType: "MODIFY",
				ON:         "row",
				TableName:  "users",
				UpdateData: `{"id":1,"name":"John Doe"}`,
			},
		},
		Deletions: []Delete{
			{
				DeletionType: "DELETE",
				ON:           "row",
				TableName:    "users",
				DeletionData: `{"id":1}`,
			},
		},
	}

	// Serialize populated ChangeLog
	jsonString, err := serializeStruct(populatedChangeLog)
	if err != nil {
		t.Fatalf("Error serializing ChangeLog: %v", err)
	}

	expectedJSON := `{"creations":[{"creation_type":"INSERT","on":"row","table_name":"users","creation_data":"{\"id\":1,\"name\":\"John\"}"}],"updates":[{"update_type":"MODIFY","on":"row","table_name":"users","update_data":"{\"id\":1,\"name\":\"John Doe\"}"}],"deletions":[{"deletion_type":"DELETE","on":"row","table_name":"users","deletion_data":"{\"id\":1}"}]}`
	if jsonString != expectedJSON {
		t.Errorf("Expected JSON: %s, got: %s", expectedJSON, jsonString)
	}

	// Define an empty ChangeLog
	emptyChangeLog := ChangeLog{
		Creations: make([]Create, 0),
		Updates:   make([]Update, 0),
		Deletions: make([]Delete, 0),
	}

	// Serialize empty ChangeLog
	jsonString, err = serializeStruct(emptyChangeLog)
	if err != nil {
		t.Fatalf("Error serializing empty ChangeLog: %v", err)
	}

	expectedEmptyJSON := `{"creations":[],"updates":[],"deletions":[]}`
	if jsonString != expectedEmptyJSON {
		t.Errorf("Expected JSON: %s, got: %s", expectedEmptyJSON, jsonString)
	}
}
