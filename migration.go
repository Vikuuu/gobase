// File storing all code for migrations
package gobase

import (
	"bytes"
	"fmt"
	"os"
)

const (
	upMigrationTemp = `-- Up Migration

`
	downMigrationTemp = `

-- Down Migration

`
	createMigFileName = "./migrations"
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
