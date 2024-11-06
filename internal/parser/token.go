package parser

type TokenType string

const (
	// Special tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Literals
	IDENT = "IDENT"

	// Misc characters
	ASTERISK = "*"
	COMMA    = ","
	LBRACE   = "{"
	RBRACE   = "}"
	LPAREN   = "("
	RPAREN   = ")"

	// Keywords
	TYPE    = "TYPE"
	STRUCT  = "STRUCT"
	PACKAGE = "PACKAGE"

	// Date Types
	INT    = "INT"
	STRING = "STRING"
	BOOL   = "BOOL"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keyword = map[string]TokenType{
	"type":    TYPE,
	"struct":  STRUCT,
	"package": PACKAGE,
	"int":     INT,
	"string":  STRING,
	"bool":    BOOL,
}
