package parser

import "fmt"

// Parser struct with token buffer for backtracking
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token
		n   int
	}
}

// NewParser creates a new instance of Parser
func NewParser(s *Scanner) *Parser {
	return &Parser{s: s}
}

// ParseStatement parse the input and returns the tokens for the statement
func (p *Parser) ParseStatement() {
	for {
		tok := p.nextToken()
		if tok.Type == EOF {
			break
		}
		fmt.Printf("Token: %+v\n", tok)
	}
}

// nextToken retrieves the next token, handling buffering for backtracking
func (p *Parser) nextToken() Token {
	// Return token from buffer if available
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok
	}

	// Otherwise read from the scanner
	tok := p.s.NextToken()

	// Save it to the buffer
	p.buf.tok = tok
	return tok
}

// unscan pushes the last token back onto the buffer
func (p *Parser) unscan() {
	p.buf.n = 1
}
