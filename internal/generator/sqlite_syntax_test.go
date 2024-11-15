package generator

import (
	"testing"
)

func TestSqLiteCreateTable(t *testing.T) {
	fileName := "../test/create_table.go"
	expectedQuery := "CREATE TABLE users (\n\tid INTEGER,\n\tname TEXT,\n\tcreated_at DATETIME,\n\tupdated_at DATETIME,\n\tis_member BOOLEAN\n);"

	outputQuery := sqLiteCreateTable(fileName)

	if expectedQuery != outputQuery {
		t.Fatalf("Not Equal. expected=%s. got=%s", expectedQuery, outputQuery)
	}
}
