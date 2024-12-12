// File storing all code for migrations
package gobase

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

const (
	upMigrationTemp = `-- Up Migration

`
	downMigrationTemp = `

-- Down Migration

`
	// createMigDirName = "./migrations"
	upMigTempLine   = "-- Up Migration"
	downMigTempLine = "-- Down Migration"
)

// calls creation func accordingly
func MigrationFile(dbConn *sql.DB, fileName, createMigDirName, MigFilename string) error {
	data, err := creationMigration(dbConn, fileName)
	if err != nil {
		return err
	}

	err = os.MkdirAll(createMigDirName, 0750)
	if err != nil {
		return fmt.Errorf("Error creating dir: %s", err)
	}
	// outFileName := createMigDirName + "/001_users.sql"
	outFileName := filepath.Join(createMigDirName, MigFilename)

	f, err := os.Create(outFileName)
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err)
	}
	defer f.Close()
	return saveMigrationFile(f.Name(), data)
}

func creationMigration(dbConn *sql.DB, fileName string) ([]byte, error) {
	schema := Parse(fileName)
	upMigQuery, downMigQuery := "", ""

	// TODO: Now check it with the previous state of the database
	// If their is no previous state, then call the table creation function
	// previous State var
	_, err := getPreviousState(dbConn)
	if err != nil {
		// This case will mean that their is no previous state
		// and this is the first time the migration is being run
		if err == sql.ErrNoRows {
			upMigQuery, downMigQuery = SqLiteCreateTable(fileName, schema)
		} else {
			return []byte{}, err
		}
	}

	var buffer bytes.Buffer

	buffer.WriteString(upMigrationTemp)
	buffer.Write([]byte(upMigQuery))
	buffer.WriteString(downMigrationTemp)
	buffer.Write([]byte(downMigQuery))

	return buffer.Bytes(), nil
}

func saveMigrationFile(outName string, data []byte) error {
	return os.WriteFile(outName, data, 0644)
}

// NOTE: Read from the given migration file and
// read only btw `--Up migration` and `--Down migration`
func getUpMigration(migrationFile string) (string, error) {
	// open the file
	file, err := os.Open(migrationFile)
	if err != nil {
		return "", fmt.Errorf("Err opening file: %s", err)
	}
	defer file.Close()

	var migration string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == downMigTempLine {
			break
		} else if line == upMigTempLine {
			continue
		} else {
			migration += line + "\n"
		}
	}

	if err = scanner.Err(); err != nil {
		return "", fmt.Errorf("Err scanning file: %s", err)
	}

	return migration, nil
}

// NOTE: Read from the given migration file and
// read only btw `--Down migration` and `EOF`
func getDownMigration(migrationFile string) (string, error) {
	file, err := os.Open(migrationFile)
	if err != nil {
		return "", fmt.Errorf("Err opening file: %s", err)
	}
	defer file.Close()

	var migration []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		migration = append(migration, line)
	}

	if err = scanner.Err(); err != nil {
		return "", fmt.Errorf("Err scanning file: %s", err)
	}

	indexDown := slices.IndexFunc(migration, func(s string) bool {
		return s == "-- Down Migration"
	})

	res := ""

	for i := indexDown + 1; i < len(migration); i++ {
		res += migration[i]
	}

	return res, nil
}

// Func responsible for migrating to the database
// Takes the dbCon, and migration file as Input
func upMigrate(dbCon *sql.DB, migrationFile string) error {
	upMigration, err := getUpMigration(migrationFile)
	if err != nil {
		return err
	}

	_, err = dbCon.Exec(upMigration)
	if err != nil {
		return fmt.Errorf("Err migrating: %s", err)
	}

	return nil
}

func downMigrate(dbCon *sql.DB, migrationFile string) error {
	downMigration, err := getDownMigration(migrationFile)
	if err != nil {
		return err
	}

	_, err = dbCon.Exec(downMigration)
	if err != nil {
		return fmt.Errorf("Err migrating: %s", err)
	}

	return nil
}
