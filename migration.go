// File storing all code for migrations
package gobase

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"slices"
)

const (
	upMigrationTemp = `-- Up Migration

`
	downMigrationTemp = `

-- Down Migration

`
	createMigFileName = "./migrations"
	upMigTempLine     = "-- Up Migration"
	downMigTempLine   = "-- Down Migration"
)

// calls creation func accordingly
func MigrationFile(fileName, createMigFileName string) error {
	data := creationMigration(fileName)
	err := os.MkdirAll(createMigFileName, 0750)
	if err != nil {
		return fmt.Errorf("Error creating dir: %s", err)
	}
	outFileName := createMigFileName + "/001_user.sql"

	f, err := os.Create(outFileName)
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err)
	}
	defer f.Close()
	return saveMigrationFile(f.Name(), data)
}

func creationMigration(fileName string) []byte {
	upMigQuery, downMigQuery := SqLiteCreateTable(fileName)

	var buffer bytes.Buffer

	buffer.WriteString(upMigrationTemp)
	buffer.Write([]byte(upMigQuery))
	buffer.WriteString(downMigrationTemp)
	buffer.Write([]byte(downMigQuery))

	return buffer.Bytes()
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
