package gobase

import (
	"fmt"
	"strings"
)

func SqLiteCreateTable(schema Schema) (upQuery, downQuery string) {
	upQuery = fmt.Sprintf("CREATE TABLE %s (\n\t", schema.SchemaName)
	for i, field := range schema.SchemaFields {
		if i == len(schema.SchemaFields)-1 {
			upQuery += fmt.Sprintf(
				"%s %s\n",
				toSnakeCase(field.Name),
				sqliteMapping[field.DataType],
			)
		} else {
			upQuery += fmt.Sprintf("%s %s,\n\t", toSnakeCase(field.Name), sqliteMapping[field.DataType])
		}
	}
	upQuery += ");"

	downQuery = fmt.Sprintf("DROP TABLE %s;", schema.SchemaName)

	return upQuery, downQuery
}

func SqliteMigration(changes ChangeLog) (upQuery, downQuery string) {
	for _, cChange := range changes.Creations {
		switch cChange.CreationType {
		case FIELD:
			cData := strings.Split(cChange.CreationData, ":")
			upQuery += fmt.Sprintf(
				"ALTER TABLE %s\nADD COLUMN %s %s;\n\n",
				cChange.TableName,
				cData[0],
				sqliteMapping[cData[1]],
			)
			downQuery += fmt.Sprintf(
				"ALTER TABLE %s\nDROP COLUMN %s;\n\n",
				cChange.TableName,
				cData[0],
			)
		}
	}

	for _, dChange := range changes.Deletions {
		switch dChange.DeletionType {
		case FIELD:
			dData := strings.Split(dChange.DeletionData, ":")
			upQuery += fmt.Sprintf(
				"ALTER TABLE %s\nDROP COLUMN %s;\n\n",
				dChange.TableName,
				dData[0],
			)
			downQuery += fmt.Sprintf(
				"ALTER TABLE %s\nADD COLUMN %s %s;\n\n",
				dChange.TableName,
				dData[0],
				sqliteMapping[dData[1]],
			)
		}
	}

	for _, uChange := range changes.Updates {
		switch uChange.UpdateType {
		case NAMEUPDATE:
			switch uChange.ON {
			case ONTABLE:
				upQuery += fmt.Sprintf(
					"ALTER TABLE %s\nRENAME TO %s;\n\n",
					uChange.TableName,
					uChange.UpdateData,
				)
				downQuery += fmt.Sprintf(
					"ALTER TABLE %s\nRENAME TO %s;\n\n",
					uChange.UpdateData,
					uChange.TableName,
				)
			case ONFIELD:
				colNameArr := strings.Split(uChange.UpdateData, ":")
				upQuery += fmt.Sprintf(
					"ALTER TABLE %s\nRENAME COLUMN %s to %s;\n\n",
					uChange.TableName,
					colNameArr[0],
					colNameArr[1],
				)
				downQuery += fmt.Sprintf(
					"ALTER TABLE %s\nRENAME COLUMN %s to %s;\n\n",
					uChange.TableName,
					colNameArr[1],
					colNameArr[0],
				)
			}
		}
	}

	return upQuery, downQuery
}
