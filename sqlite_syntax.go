package gobase

import (
	"fmt"
)

func SqLiteCreateTable(fileName string) string {
	schema := Parse(fileName)
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

	return createQuery
}
