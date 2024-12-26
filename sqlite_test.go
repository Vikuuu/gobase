package gobase

import (
	"testing"
)

func TestSqliteConn(t *testing.T) {
	db, err := SqliteConn(DBFILENAME)
	if err != nil {
		t.Fatalf("Error creating conn. error=%s", err)
	}

	db.Close()
}
