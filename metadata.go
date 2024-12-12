// All the metadata needed for automated changes detection
package gobase

import (
	"database/sql"
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
LIMIT 1;
    `

	md := metadata{}

	row := dbCon.QueryRow(query)
	err := row.Scan(&md)
	if err != nil {
		return md, err
	}

	return md, nil
}
