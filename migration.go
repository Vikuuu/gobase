// File storing all code for migrations
package gobase

import (
	"bufio"
	"bytes"
	"database/sql"
	"errors"
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

var ErrNoChange = errors.New("No changes detected")

func MigrationFile(
	dbConn *sql.DB,
	fileName, createMigDirName, MigFilename string,
) (newJSON, changeJSON string, err error) {
	data, newJSON, changeJSON, err := creationMigration(dbConn, fileName)
	if err != nil {
		return "", "", err
	}

	err = os.MkdirAll(createMigDirName, 0750)
	if err != nil {
		return "", "", fmt.Errorf("Error creating dir: %s", err)
	}
	// outFileName := createMigDirName + "/001_users.sql"
	outFileName := filepath.Join(createMigDirName, MigFilename)

	f, err := os.Create(outFileName)
	if err != nil {
		return "", "", fmt.Errorf("Error creating file: %s", err)
	}
	defer f.Close()
	return newJSON, changeJSON, saveMigrationFile(f.Name(), data)
}

func creationMigration(
	dbConn *sql.DB,
	fileName string,
) (migrationSyntax []byte, newState, changes string, err error) {
	newSchema := Parse(fileName)
	newState, err = serializeSchema(newSchema)
	if err != nil {
		return nil, "", "", nil
	}

	prevState, err := getPreviousState(dbConn)
	if err != nil {
		if err == sql.ErrNoRows {
			upMigQuery, downMigQuery := SqLiteCreateTable(newSchema)
			migrationSyntax = writeBuffer(upMigQuery, downMigQuery)
			return migrationSyntax, newState, changes, nil
		} else if err.Error() == "no such table: gobase_metadata" {
			err := CreateMetaTable(dbConn)
			if err != nil {
				return nil, "", "", err
			}
			return creationMigration(dbConn, fileName)
		}
		return nil, "", "", err
	}
	upMigQuery, downMigQuery, changes, err := generateMigrationQueries(
		newSchema,
		prevState.CurrentState,
	)
	migrationSyntax = writeBuffer(upMigQuery, downMigQuery)
	return migrationSyntax, newState, changes, nil
}

func generateMigrationQueries(
	newSchema Schema,
	prevState string,
) (upQuery, downQuery, changeJSON string, err error) {
	prevSchema, err := deserializeSchema(prevState)
	if err != nil {
		return "", "", "", err
	}
	changes := categorizeSchemaChanges(newSchema, prevSchema)
	if len(changes.Creations) == 0 && len(changes.Deletions) == 0 && len(changes.Updates) == 0 {
		return "", "", "", ErrNoChange
	}
	upQuery, downQuery = SqliteMigration(changes)
	changesJSON, err := serializeStruct[ChangeLog](changes)
	if err != nil {
		return "", "", "", err
	}
	return upQuery, downQuery, changesJSON, nil
}

func writeBuffer(up, down string) []byte {
	var buffer bytes.Buffer
	buffer.WriteString(upMigrationTemp)
	buffer.Write([]byte(up))
	buffer.WriteString(downMigrationTemp)
	buffer.Write([]byte(down))
	return buffer.Bytes()
}

func saveMigrationFile(outName string, data []byte) error {
	return os.WriteFile(outName, data, 0644)
}

// NOTE: Read from the given migration file and
// read only btw `--Up migration` and `--Down migration`
func getUpMigration(migrationFile string) (string, error) {
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
func UpMigrate(dbCon *sql.DB, migrationFile, newJSON, changeJSON string) error {
	upMigration, err := getUpMigration(migrationFile)
	if err != nil {
		return err
	}

	_, err = dbCon.Exec(upMigration)
	if err != nil {
		return fmt.Errorf("Err migrating: %s", err)
	}

	// Update the metadata table to add the current state.
	err = updateMetadata(dbCon, newJSON, changeJSON)
	if err != nil {
		return err
	}
	return nil
}

func DownMigrate(dbCon *sql.DB, migrationFile string) error {
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
