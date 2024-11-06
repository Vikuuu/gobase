/*
* The aim of this file is to parse the given go file that stores
* the schema data into parsed file that will be changed into the
* SQL equivalent syntax.
 */

package parser

import (
	"unicode"
)

func LookupIdent(ident string) TokenType {
	if tok, ok := keyword[ident]; ok {
		return tok
	}
	return IDENT
}

// Utility functions to check character type
func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}
