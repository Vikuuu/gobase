package generator

import (
	"fmt"

	"github.com/Vikuuu/gobase/internal/parser"
)

func sqLiteCreateTable(fileName string) string {
	schema := parser.Parse(fileName)
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
