package gobase

import (
	"testing"
)

func TestSqLiteCreateTable(t *testing.T) {
	test := []struct {
		name              string
		fileName          string
		expectedUpQuery   string
		expectedDownQuery string
	}{
		{
			name:              "Single Table",
			fileName:          "./testdata/create_table.go",
			expectedUpQuery:   "CREATE TABLE users (\n\tid INTEGER,\n\tname TEXT,\n\tcreated_at DATETIME,\n\tupdated_at DATETIME,\n\tis_member BOOLEAN\n);\n",
			expectedDownQuery: "DROP TABLE users;\n",
		},
		{
			name:              "Two Table",
			fileName:          "./testdata/multi_table_test.go",
			expectedUpQuery:   "CREATE TABLE users (\n\tid INTEGER,\n\tname TEXT,\n\tcreated_at DATETIME,\n\tupdated_at DATETIME,\n\tis_member BOOLEAN\n);\nCREATE TABLE image (\n\tname TEXT,\n\ttype TEXT,\n\tsize INTEGER,\n\thidden BOOLEAN,\n\tcreated_at DATETIME,\n\tupdated_at DATETIME\n);\n",
			expectedDownQuery: "DROP TABLE image;\nDROP TABLE users;\n",
		},
	}

	for _, tt := range test {
		schema := Parse(tt.fileName)

		outputUpQuery, outputDownQuery := SqLiteCreateTable(schema)

		if tt.expectedUpQuery != outputUpQuery {
			t.Errorf("Up query err. expected=%s. got=%s", tt.expectedUpQuery, outputUpQuery)
		}
		if tt.expectedDownQuery != outputDownQuery {
			t.Errorf("Down query err. expected=%s. got=%s", tt.expectedDownQuery, outputDownQuery)
		}
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
			expUp:   "ALTER TABLE users\nADD COLUMN id INTEGER;\n\n",
			expDown: "ALTER TABLE users DROP COLUMN id;\n\n",
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
			expUp:   "ALTER TABLE users\nDROP COLUMN id;\n\n",
			expDown: "ALTER TABLE users ADD COLUMN id INTEGER;\n\n",
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
			expUp:   "ALTER TABLE users\nRENAME TO accounts;\n\n",
			expDown: "ALTER TABLE accounts RENAME TO users;\n\n",
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
			expUp:   "ALTER TABLE users\nRENAME COLUMN id to user_id;\n\n",
			expDown: "ALTER TABLE users RENAME COLUMN user_id to id;\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUp, gotDown := SqliteMigration(tt.change)
			if gotUp != tt.expUp {
				t.Errorf("Up Mig Not Same\nExpected: %q\nGot: %q", tt.expUp, gotUp)
			}
			if gotDown != tt.expDown {
				t.Errorf("Down Mig Not Same\nExpected: %q\nGot: %q", tt.expDown, gotDown)
			}
		})
	}
}
