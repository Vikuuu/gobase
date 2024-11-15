package generator

var sqliteMapping = map[string]string{
	"int":    "INTEGER",
	"float":  "REAL",
	"string": "TEXT",
	"Time":   "DATETIME",
	"bool":   "BOOLEAN",
}
