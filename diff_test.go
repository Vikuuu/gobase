package gobase

import (
	"reflect"
	"testing"
)

func TestCategorizeSchemaChanges(t *testing.T) {
	tests := []struct {
		name      string
		oldSchema Schema
		newSchema Schema
		expected  ChangeLog
	}{
		{
			name: "No changes",
			oldSchema: Schema{
				SchemaName: "users",
				SchemaFields: []SchemaField{
					{Name: "id", DataType: "int"},
					{Name: "name", DataType: "string"},
				},
			},
			newSchema: Schema{
				SchemaName: "users",
				SchemaFields: []SchemaField{
					{Name: "id", DataType: "int"},
					{Name: "name", DataType: "string"},
				},
			},
			expected: ChangeLog{},
		},
		{
			name: "Schema name updated",
			oldSchema: Schema{
				SchemaName: "users",
				SchemaFields: []SchemaField{
					{Name: "id", DataType: "int"},
				},
			},
			newSchema: Schema{
				SchemaName: "accounts",
				SchemaFields: []SchemaField{
					{Name: "id", DataType: "int"},
				},
			},
			expected: ChangeLog{
				Updates: []Update{
					{
						UpdateType: NAMEUPDATE,
						ON:         ONTABLE,
						TableName:  "users",
						UpdateData: "accounts",
					},
				},
			},
		},
		{
			name: "Field created",
			oldSchema: Schema{
				SchemaName: "users",
				SchemaFields: []SchemaField{
					{Name: "id", DataType: "int"},
				},
			},
			newSchema: Schema{
				SchemaName: "users",
				SchemaFields: []SchemaField{
					{Name: "id", DataType: "int"},
					{Name: "email", DataType: "string"},
				},
			},
			expected: ChangeLog{
				Creations: []Create{
					{
						CreationType: FIELD,
						ON:           ONTABLE,
						TableName:    "users",
						CreationData: "email:string",
					},
				},
			},
		},
		{
			name: "Field deleted",
			oldSchema: Schema{
				SchemaName: "users",
				SchemaFields: []SchemaField{
					{Name: "id", DataType: "int"},
					{Name: "email", DataType: "string"},
				},
			},
			newSchema: Schema{
				SchemaName: "users",
				SchemaFields: []SchemaField{
					{Name: "id", DataType: "int"},
				},
			},
			expected: ChangeLog{
				Deletions: []Delete{
					{
						DeletionType: FIELD,
						ON:           ONTABLE,
						TableName:    "users",
						DeletionData: "email:string",
					},
				},
			},
		},
		// NOTE: This feature is not supported in SQLite3
		// {
		// 	name: "Field type updated",
		// 	oldSchema: Schema{
		// 		SchemaName: "users",
		// 		SchemaFields: []SchemaField{
		// 			{Name: "id", DataType: "int"},
		// 		},
		// 	},
		// 	newSchema: Schema{
		// 		SchemaName: "users",
		// 		SchemaFields: []SchemaField{
		// 			{Name: "id", DataType: "string"},
		// 		},
		// 	},
		// 	expected: ChangeLog{
		// 		Updates: []string{"Field type updated: id (int -> string)"},
		// 	},
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := categorizeSchemaChanges(test.oldSchema, test.newSchema)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf(
					"Test %s failed.\nExpected: %+v\nGot: %+v",
					test.name,
					test.expected,
					result,
				)
			}
		})
	}
}
