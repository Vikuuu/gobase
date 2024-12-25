package gobase

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
)

const (
	testFile    = "./testdata/create_table.go"
	testMigDir  = "./testdata/migrations"
	testMigFile = "./testdata/migrations/001_users.sql"
)

var expected = fmt.Sprintf(
	"-- Up Migration\n\nCREATE TABLE users (\n\tid INTEGER,\n\tname TEXT,\n\tcreated_at DATETIME,\n\tupdated_at DATETIME,\n\tis_member BOOLEAN\n);\n\n-- Down Migration\n\nDROP TABLE users;",
)

func testPrerequisite(dbConn *sql.DB) error {
	err := CreateMetaTable(dbConn)
	if err != nil {
		return err
	}
	return nil
}

func TestCreationMigrationFile(t *testing.T) {
	dbCon, err := SqliteConn(DBFILENAME)
	if err != nil {
		t.Fatalf("Err connecting: %s", err)
	}
	defer dbCon.Close()

	// err = testPrerequisite(dbCon)
	// if err != nil {
	// 	t.Fatalf("Err creating meta table: %s", err)
	// }

	// Call the function under test
	got, _, _, err := creationMigration(dbCon, testFile)
	if err != nil {
		t.Fatalf("Error creating mig: %s", err)
	}

	// Compare results
	if expected != string(got) {
		t.Errorf("Mismatch in migration SQL:\nExpected:\n%s\n\nGot:\n%s", expected, got)
	}

	os.Remove(DBFILENAME)
}

func TestMigrationFile(t *testing.T) {
	dbConn, err := SqliteConn(DBFILENAME)
	if err != nil {
		t.Fatalf("Err connecting to db: %s", err)
	}
	defer dbConn.Close()

	err = testPrerequisite(dbConn)
	if err != nil {
		t.Fatalf("Err creating meta table: %s", err)
	}

	_, _, err = MigrationFile(dbConn, testFile, testMigDir, "001_users.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat("./testdata/migrations/001_users.sql")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}

	got, err := os.ReadFile("./testdata/migrations/001_users.sql")
	if err != nil {
		t.Fatal(err)
	}

	if expected != string(got) {
		t.Errorf("Mismatch in migration SQL:\nExpected:\n%s\n\nGot:\n%s", expected, got)
	}

	os.Remove(DBFILENAME)
}

func TestGetUpMigration(t *testing.T) {
	expUpMig := fmt.Sprintf(
		"\nCREATE TABLE users (\n\tid INTEGER,\n\tname TEXT,\n\tcreated_at DATETIME,\n\tupdated_at DATETIME,\n\tis_member BOOLEAN\n);\n\n",
	)
	gotUpMig, err := getUpMigration("./testdata/migrations/001_users.sql")
	if err != nil {
		t.Fatalf("Err getting Up mig: %s", err)
	}

	if expUpMig != gotUpMig {
		t.Errorf("Migration not equal. expected: \n%s\n. got: \n%s\n", expUpMig, gotUpMig)
	}
}

func TestGetDownMigration(t *testing.T) {
	expDownMig := fmt.Sprintf("DROP TABLE users;")
	gotDownMig, err := getDownMigration("./testdata/migrations/001_users.sql")
	if err != nil {
		t.Fatalf("Err getting Down mig: %s", err)
	}

	if expDownMig != gotDownMig {
		t.Errorf("Migration not equal. expected: \n%s\n. got: \n%s\n", expDownMig, gotDownMig)
	}
}

func TestMigration(t *testing.T) {
	dbCon, err := SqliteConn(":memory:")
	if err != nil {
		t.Fatalf("Err conn db: %s", err)
	}
	testPrerequisite(dbCon)

	testCase := []struct {
		name   string
		testFn func(*sql.DB, string) error
		result string
	}{
		{name: "TestUpMigration", testFn: func(db *sql.DB, mFile string) error {
			return UpMigrate(db, mFile, "", "")
		}, result: "users"},
		{name: "TestDownMigration", testFn: DownMigrate, result: ""},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			// Test both UP and DOWN migrations
			// First the UP mig
			err := tc.testFn(dbCon, testMigFile)
			if err != nil {
				t.Fatalf("Err migrating: %s", err)
			}

			row, err := dbCon.Query(
				`SELECT name FROM sqlite_schema WHERE type='table' AND name NOT LIKE 'sqlite_%' AND 'gobase_metadata'`,
			)
			if err != nil {
				t.Fatalf("Err getting table rows: %s", err)
			}

			for row.Next() {
				var name string
				if err := row.Scan(&name); err != nil {
					t.Fatal(err)
				}

				if name != tc.result {
					t.Errorf("Table not match. expected: %s. got: %s", tc.result, name)
				}
			}
		})
	}
}
