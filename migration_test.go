package gobase

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

const (
	testFile   = "./testdata/create_table.go"
	testMigDir = "./testdata/migrations"
)

var expected = fmt.Sprintf(
	"-- Up Migration\n\nCREATE TABLE users (\n\tid INTEGER,\n\tname TEXT,\n\tcreated_at DATETIME,\n\tupdated_at DATETIME,\n\tis_member BOOLEAN\n);\n\n-- Down Migration\n\nDROP TABLE users;",
)

func TestCreationMigrationFile(t *testing.T) {
	// Call the function under test
	got := creationMigration(testFile)

	// Compare results
	if expected != string(got) {
		t.Errorf("Mismatch in migration SQL:\nExpected:\n%s\n\nGot:\n%s", expected, got)
	}
}

func TestMigrationFile(t *testing.T) {
	err := MigrationFile(testFile, testMigDir)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat("./testdata/migrations/001_user.sql")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}

	got, err := os.ReadFile("./testdata/migrations/001_user.sql")
	if err != nil {
		t.Fatal(err)
	}

	if expected != string(got) {
		t.Errorf("Mismatch in migration SQL:\nExpected:\n%s\n\nGot:\n%s", expected, got)
	}
}

func TestGetUpMigration(t *testing.T) {
	expUpMig := fmt.Sprintf(
		"\nCREATE TABLE users (\n\tid INTEGER,\n\tname TEXT,\n\tcreated_at DATETIME,\n\tupdated_at DATETIME,\n\tis_member BOOLEAN\n);\n\n",
	)
	gotUpMig, err := getUpMigration("./testdata/migrations/001_user.sql")
	if err != nil {
		t.Fatalf("Err getting Up mig: %s", err)
	}

	if expUpMig != gotUpMig {
		t.Errorf("Migration not equal. expected: \n%s\n. got: \n%s\n", expUpMig, gotUpMig)
	}
}

func TestGetDownMigration(t *testing.T) {
	expDownMig := fmt.Sprintf("DROP TABLE users;")
	gotDownMig, err := getDownMigration("./testdata/migrations/001_user.sql")
	if err != nil {
		t.Fatalf("Err getting Down mig: %s", err)
	}

	if expDownMig != gotDownMig {
		t.Errorf("Migration not equal. expected: \n%s\n. got: \n%s\n", expDownMig, gotDownMig)
	}
}
