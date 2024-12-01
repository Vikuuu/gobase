package gobase

import (
	"testing"
)

func TestSqLiteCreateTable(t *testing.T) {
	fileName := "./testdata/create_table.go"
	expectedUpQuery := "CREATE TABLE users (\n\tid INTEGER,\n\tname TEXT,\n\tcreated_at DATETIME,\n\tupdated_at DATETIME,\n\tis_member BOOLEAN\n);"
	expectedDownQuery := "DROP TABLE users;"

	outputUpQuery, outputDownQuery := SqLiteCreateTable(fileName)

	if expectedUpQuery != outputUpQuery {
		t.Errorf("Up query err. expected=%s. got=%s", expectedUpQuery, outputUpQuery)
	}
	if expectedDownQuery != outputDownQuery {
		t.Errorf("Down query err. expected=%s. got=%s", expectedDownQuery, outputDownQuery)
	}
}
