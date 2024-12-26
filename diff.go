package gobase

import (
	"fmt"
)

type ChangeLog struct {
	Creations []Create `json:"creations"`
	Updates   []Update `json:"updates"`
	Deletions []Delete `json:"deletions"`
}

type Create struct {
	CreationType string `json:"creation_type"`
	ON           string `json:"on"`
	TableName    string `json:"table_name"`
	CreationData string `json:"creation_data"`
}

type Update struct {
	UpdateType string `json:"update_type"`
	ON         string `json:"on"`
	TableName  string `json:"table_name"`
	UpdateData string `json:"update_data"`
}

type Delete struct {
	DeletionType string `json:"deletion_type"`
	ON           string `json:"on"`
	TableName    string `json:"table_name"`
	DeletionData string `json:"deletion_data"`
}

const (
	// Update Types
	NAMEUPDATE     = "Name"
	DATATYPEUPDATE = "DataType"
	// Update Action On
	ONTABLE = "Table"
	ONFIELD = "Field"

	// Creation Types
	FIELD = "Field"
)

func categorizeSchemaChanges(oldSchema, newSchema Schema) ChangeLog {
	changes := ChangeLog{}

	// 1. compare schema name
	if oldSchema.SchemaName != newSchema.SchemaName {
		changes.Updates = append(
			changes.Updates,
			Update{
				UpdateType: NAMEUPDATE,
				ON:         ONTABLE,
				TableName:  oldSchema.SchemaName,
				UpdateData: newSchema.SchemaName,
			},
		)
	}

	// 2. Convert SchemaFields to maps for easy lookup
	oldFieldsMap := make(map[string]string)
	for _, field := range oldSchema.SchemaFields {
		oldFieldsMap[field.Name] = field.DataType
	}

	newFieldsMap := make(map[string]string)
	for _, field := range newSchema.SchemaFields {
		newFieldsMap[field.Name] = field.DataType
	}

	// 3. Detect Field Creation and Updates
	for name, newType := range newFieldsMap {
		_, exists := oldFieldsMap[name]
		if !exists {
			changes.Creations = append(
				changes.Creations,
				Create{
					CreationType: FIELD,
					ON:           ONTABLE,
					TableName:    oldSchema.SchemaName,
					CreationData: fmt.Sprintf("%s:%s", name, newType),
				},
			)
		}
		// NOTE: The feature of changing the data type of table column is not supported, so we can skip this for now,
		// will need to found a way to do so in other dbs.
		// else if oldType != newType {
		//	changes.Updates = append(changes.Updates, fmt.Sprintf("Field type updated: %s (%s -> %s)", name, oldType, newType))
		// }
	}

	// 4. Detect field deletions
	for name := range oldFieldsMap {
		if _, exists := newFieldsMap[name]; !exists {
			changes.Deletions = append(
				changes.Deletions,
				Delete{
					DeletionType: FIELD,
					ON:           ONTABLE,
					TableName:    oldSchema.SchemaName,
					DeletionData: fmt.Sprintf("%s:%s", name, oldFieldsMap[name]),
				})
		}
	}

	return changes
}
