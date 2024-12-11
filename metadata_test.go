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
