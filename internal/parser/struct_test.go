package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
	testFileName := "../test/create_table.go"
	expectedSchema := Schema{
		SchemaName: "users",
		SchemaFields: []struct {
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

	if expectedSchema.SchemaName != outputSchema.SchemaName {
		t.Fatalf(
			"Struct name not Equal. expected=%s. got=%s",
			expectedSchema.SchemaName,
			outputSchema.SchemaName,
		)
	}

	if len(expectedSchema.SchemaFields) != len(outputSchema.SchemaFields) {
		t.Fatalf(
			"Number of fields not equal. expected=%d, got=%d",
			len(expectedSchema.SchemaFields),
			len(outputSchema.SchemaFields),
		)
	}

	for i := range expectedSchema.SchemaFields {
		expectedField := expectedSchema.SchemaFields[i]
		outputField := outputSchema.SchemaFields[i]

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
