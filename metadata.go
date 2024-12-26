package gobase

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// This is to be used to create the query when the tags
// usable is added, currently not supported.
type GobaseMetadata struct {
	ID           int       `gobase:"primary_key"`
	CurrentState string    `gobase:""`
	Changes_made string    `gobase:""`
	CreatedAt    time.Time `gobase:"default:current_timestamp"`
}

func CreateMetaTable(dbCon *sql.DB) error {
	query := `CREATE TABLE gobase_metadata (
    id INTEGER PRIMARY KEY,
    current_state TEXT,
    changes_made TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	_, err := dbCon.Exec(query)
	if err != nil {
		return fmt.Errorf("Err migrating: %s", err)
	}

	return nil
}

type metadata struct {
	ID           int
	CurrentState string
	ChangesMade  string
	CreatedAt    time.Time
}

func getPreviousState(dbCon *sql.DB) (metadata, error) {
	query := `
SELECT id, current_state, changes_made, created_at
FROM gobase_metadata
ORDER BY id DESC
LIMIT 1;
    `
	md := metadata{}
	row := dbCon.QueryRow(query)
	err := row.Scan(&md.ID, &md.CurrentState, &md.ChangesMade, &md.CreatedAt)
	if err != nil {
		return md, err
	}

	return md, nil
}

func GetLatestId(dbCon *sql.DB) (int, error) {
	query := `SELECT id
FROM gobase_metadata
ORDER BY id DESC
LIMIT 1;`
	var lastestVersion int
	row := dbCon.QueryRow(query)
	err := row.Scan(&lastestVersion)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		} else if err.Error() == "no such table: gobase_metadata" {
			return 0, nil
		}
		return 0, err
	}
	return lastestVersion, nil
}

func updateMetadata(dbCon *sql.DB, newJSON, changeJSON string) error {
	query := `INSERT INTO gobase_metadata(current_state, changes_made)
VALUES (?, ?);
`
	_, err := dbCon.Exec(query, newJSON, changeJSON)
	if err != nil {
		return err
	}
	return nil
}

func serializeSchema(schema Schema) (string, error) {
	jsonData, err := json.Marshal(schema)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func deserializeSchema(schema string) (Schema, error) {
	var s Schema
	err := json.Unmarshal([]byte(schema), &s)
	if err != nil {
		return Schema{}, err
	}
	return s, nil
}
