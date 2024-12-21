package gobase

import (
	"os"
	"testing"
)

func TestCreateMetaTable(t *testing.T) {
	dbConn, err := SqliteConn(DBFILENAME)
	if err != nil {
		t.Fatalf("Error getting conn: %s", err)
	}
	defer dbConn.Close()

	err = CreateMetaTable(dbConn)
	if err != nil {
		t.Fatalf("Error exec query: %s", err)
	}

	expectedTable := "gobase_metadata"
	row, err := dbConn.Query(
		`SELECT name FROM sqlite_schema WHERE type='table' AND name NOT LIKE 'sqlite_%'`,
	)
	if err != nil {
		t.Fatalf("Err getting table rows: %s", err)
	}

	for row.Next() {
		var resultTable string
		if err := row.Scan(&resultTable); err != nil {
			t.Fatalf("Err scanning value: %s", err)
		}
		if resultTable != expectedTable {
			t.Errorf("Table not match. expected: %s. got: %s", expectedTable, resultTable)
		}
	}

	os.Remove(DBFILENAME)
}

func TestSerializeSchema(t *testing.T) {
	tests := []struct {
		name         string
		input        Schema
		expectedJSON string
		expectError  bool
	}{
		{
			name: "Valid schema with fields",
			input: Schema{
				SchemaName: "users",
				SchemaFields: []SchemaField{
					{"id", "INTEGER"},
					{"name", "TEXT"},
					{"email", "TEXT"},
				},
			},
			expectedJSON: `{"schema_name":"users","schema_fields":[{"name":"id","data_type":"INTEGER"},{"name":"name","data_type":"TEXT"},{"name":"email","data_type":"TEXT"}]}`,
			expectError:  false,
		},
		{
			name: "Empty schema fields",
			input: Schema{
				SchemaName:   "empty_table",
				SchemaFields: []SchemaField{},
			},
			expectedJSON: `{"schema_name":"empty_table","schema_fields":[]}`,
			expectError:  false,
		},
		{
			name: "Empty schema name and fields",
			input: Schema{
				SchemaName:   "",
				SchemaFields: []SchemaField{},
			},
			expectedJSON: `{"schema_name":"","schema_fields":[]}`,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serializeSchema(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("serializeSchema() error = %v, expectError = %v", err, tt.expectError)
			}
			if got != tt.expectedJSON {
				t.Errorf("serializeSchema() got = %s, want = %s", got, tt.expectedJSON)
			}
		})
	}
}
