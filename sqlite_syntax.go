package gobase

import (
	"fmt"
	"strings"
)

func SqLiteCreateTable(schema Schema) (upQuery string, downQuery string) {
	createQuery := fmt.Sprintf("CREATE TABLE %s (\n\t", schema.SchemaName)
	for i, field := range schema.SchemaFields {
		if i == len(schema.SchemaFields)-1 {
			createQuery += fmt.Sprintf(
				"%s %s\n",
				toSnakeCase(field.Name),
				sqliteMapping[field.DataType],
			)
		} else {
			createQuery += fmt.Sprintf("%s %s,\n\t", toSnakeCase(field.Name), sqliteMapping[field.DataType])
		}
	}
	createQuery += ");"

	dropQuery := fmt.Sprintf("DROP TABLE %s;", schema.SchemaName)

	return createQuery, dropQuery
}

func SqliteMigration(changes ChangeLog) (upQuery string, downQuery string) {
	uQuery, dQuery := "", ""

	for _, cChange := range changes.Creations {
		switch cChange.CreationType {
		case FIELD:
			cData := strings.Split(cChange.CreationData, ":")
			uQuery += fmt.Sprintf(
				"ALTER TABLE %s\nADD COLUMN %s %s;\n",
				cChange.TableName,
				cData[0],
				sqliteMapping[cData[1]],
			)
			dQuery += fmt.Sprintf(
				"ALTER TABLE %s\nDROP COLUMN %s;\n",
				cChange.TableName,
				cData[0],
			)
		}
	}

	for _, dChange := range changes.Deletions {
		switch dChange.DeletionType {
		case FIELD:
			dData := strings.Split(dChange.DeletionData, ":")
			uQuery += fmt.Sprintf(
				"ALTER TABLE %s\nDROP COLUMN %s;\n",
				dChange.TableName,
				dData[0],
			)
			dQuery += fmt.Sprintf(
				"ALTER TABLE %s\nADD COLUMN %s %s;\n",
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
				uQuery += fmt.Sprintf(
					"ALTER TABLE %s\nRENAME TO %s;\n",
					uChange.TableName,
					uChange.UpdateData,
				)
				dQuery += fmt.Sprintf(
					"ALTER TABLE %s\nRENAME TO %s;\n",
					uChange.UpdateData,
					uChange.TableName,
				)
			case ONFIELD:
				colNameArr := strings.Split(uChange.UpdateData, ":")
				uQuery += fmt.Sprintf(
					"ALTER TABLE %s\nRENAME COLUMN %s to %s;\n",
					uChange.TableName,
					colNameArr[0],
					colNameArr[1],
				)
				dQuery += fmt.Sprintf(
					"ALTER TABLE %s\nRENAME COLUMN %s to %s;\n",
					uChange.TableName,
					colNameArr[1],
					colNameArr[0],
				)
			}
		}
	}

	return uQuery, dQuery
}
