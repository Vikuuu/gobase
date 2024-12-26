package gobase

import (
	"testing"
)

func TestSqLiteCreateTable(t *testing.T) {
	fileName := "./testdata/create_table.go"
	expectedUpQuery := "CREATE TABLE users (\n\tid INTEGER,\n\tname TEXT,\n\tcreated_at DATETIME,\n\tupdated_at DATETIME,\n\tis_member BOOLEAN\n);"
	expectedDownQuery := "DROP TABLE users;"
	schema := Parse(fileName)

	outputUpQuery, outputDownQuery := SqLiteCreateTable(schema)

	if expectedUpQuery != outputUpQuery {
		t.Errorf("Up query err. expected=%s. got=%s", expectedUpQuery, outputUpQuery)
	}
	if expectedDownQuery != outputDownQuery {
		t.Errorf("Down query err. expected=%s. got=%s", expectedDownQuery, outputDownQuery)
	}
}

func TestSqliteMigration(t *testing.T) {
	tests := []struct {
		name    string
		change  ChangeLog
		expUp   string
		expDown string
	}{
		{
			name:    "No Change",
			change:  ChangeLog{},
			expUp:   "",
			expDown: "",
		},
		{
			name: "Field Creation",
			change: ChangeLog{
				Creations: []Create{
					{
						CreationType: FIELD,
						ON:           ONTABLE,
						TableName:    "users",
						CreationData: "id:int",
					},
				},
			},
			expUp:   "ALTER TABLE users\nADD COLUMN id INTEGER;\n",
			expDown: "ALTER TABLE users\nDROP COLUMN id;\n",
		},
		{
			name: "Field Deletion",
			change: ChangeLog{
				Deletions: []Delete{
					{
						DeletionType: FIELD,
						ON:           ONTABLE,
						TableName:    "users",
						DeletionData: "id:int",
					},
				},
			},
			expUp:   "ALTER TABLE users\nDROP COLUMN id;\n",
			expDown: "ALTER TABLE users\nADD COLUMN id INTEGER;\n",
		},
		{
			name: "Table Rename",
			change: ChangeLog{
				Updates: []Update{
					{
						UpdateType: NAMEUPDATE,
						ON:         ONTABLE,
						TableName:  "users",
						UpdateData: "accounts",
					},
				},
			},
			expUp:   "ALTER TABLE users\nRENAME TO accounts;\n",
			expDown: "ALTER TABLE accounts\nRENAME TO users;\n",
		},
		{
			name: "Field Rename",
			change: ChangeLog{
				Updates: []Update{
					{
						UpdateType: NAMEUPDATE,
						ON:         ONFIELD,
						TableName:  "users",
						UpdateData: "id:user_id",
					},
				},
			},
			expUp:   "ALTER TABLE users\nRENAME COLUMN id to user_id;\n",
			expDown: "ALTER TABLE users\nRENAME COLUMN user_id to id;\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUp, gotDown := SqliteMigration(tt.change)
			if gotUp != tt.expUp {
				t.Errorf("Up Mig Not Same\nExpected: %s\nGot: %s", tt.expUp, gotUp)
			}
			if gotDown != tt.expDown {
				t.Errorf("Down Mig Not Same\nExpected: %s\nGot: %s", tt.expDown, gotDown)
			}
		})
	}
}
