package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
	testFileName := "./test_example.go"
	expectedSchema := Schema{
		schemaName: "users",
		schemaFields: []struct {
			Name     string
			DataType string
		}{
			{Name: "ID", DataType: "int"},
			{Name: "Name", DataType: "string"},
			{Name: "CreatedAt", DataType: "Time"},
			{Name: "UpdatedAt", DataType: "Time"},
			{Name: "IsMember", DataType: "bool"},
		},
	}

	outputSchema := Parse(testFileName)

	if expectedSchema.schemaName != outputSchema.schemaName {
		t.Fatalf(
			"Struct name not Equal. expected=%s. got=%s",
			expectedSchema.schemaName,
			outputSchema.schemaName,
		)
	}

	if len(expectedSchema.schemaFields) != len(outputSchema.schemaFields) {
		t.Fatalf(
			"Number of fields not equal. expected=%d, got=%d",
			len(expectedSchema.schemaFields),
			len(outputSchema.schemaFields),
		)
	}

	for i := range expectedSchema.schemaFields {
		expectedField := expectedSchema.schemaFields[i]
		outputField := outputSchema.schemaFields[i]

		if expectedField.Name != outputField.Name {
			t.Fatalf(
				"Field name not equal at index %d. expected=%s, got=%s",
				i,
				expectedField.Name,
				outputField.Name,
			)
		}

		if expectedField.DataType != outputField.DataType {
			t.Fatalf(
				"Field data type not equal at index %d. expected=%s, got=%s",
				i,
				expectedField.DataType,
				outputField.DataType,
			)
		}
	}
}
