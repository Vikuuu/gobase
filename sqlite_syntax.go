package gobase

import (
	"fmt"
)

func SqLiteCreateTable(fileName string, schema Schema) (cQuery string, dQuery string) {
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
