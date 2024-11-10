package parser

// Lex struct with token buffer for backtracking
type Lex struct {
	s   *Scanner
	buf struct {
		tok Token
		n   int
	}
}

// NewLexer creates a new instance of Parser
func NewLexer(s *Scanner) *Lex {
	return &Lex{s: s}
}

// ParseStatement parse the input and returns the tokens for the statement
func (p *Lex) Lexer() {
	for {
		tok := p.nextToken()
		if tok.Type == EOF {
			break
		}
	}
}

// nextToken retrieves the next token, handling buffering for backtracking
func (p *Lex) nextToken() Token {
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
func (p *Lex) unscan() {
	p.buf.n = 1
}
