package parser

type Scanner struct {
	input        string
	position     int  // current position in input
	readPosition int  // current reading position
	ch           byte // current character
}

// NewScanner creates a new instance of Scanner
func NewScanner(input string) *Scanner {
	s := &Scanner{input: input}
	s.readChar() // Initialize first character
	return s
}

// readChar read the next character and advances the position
func (s *Scanner) readChar() {
	if s.readPosition >= len(s.input) {
		s.ch = 0
	} else {
		s.ch = s.input[s.readPosition]
	}
	s.position = s.readPosition
	s.readPosition++
}

// peekChar allows lookahead by checking the next character without advancing
func (s *Scanner) peekChar() byte {
	if s.readPosition >= len(s.input) {
		return 0
	}
	return s.input[s.readPosition]
}

// skipWhitespace skips over spaces, tabs, and newlines
func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.readChar()
	}
}

// NextToken generates the next token from input
func (s *Scanner) NextToken() Token {
	var tok Token
	s.skipWhitespace()

	switch s.ch {
	case '*':
		tok = Token{Type: ASTERISK, Literal: string(s.ch)}
	case ',':
		tok = Token{Type: COMMA, Literal: string(s.ch)}
	case '{':
		tok = Token{Type: LBRACE, Literal: string(s.ch)}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(s.ch)}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(s.ch)}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(s.ch)}
	case 0:
		tok = Token{Type: EOF, Literal: ""}
	default:
		if isLetter(s.ch) {
			literal := s.readIdentifier()
			tokType := LookupIdent(literal)
			tok = Token{Type: tokType, Literal: literal}
			return tok
		} else if isDigit(s.ch) {
			literal := s.readNumber()
			tok = Token{Type: INT, Literal: literal}
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(s.ch)}
		}
	}
	s.readChar()
	return tok
}

// readIdentifier reads a sequence of letters as an identifier
func (s *Scanner) readIdentifier() string {
	position := s.position
	for isLetter(s.ch) {
		s.readChar()
	}
	return s.input[position:s.position]
}

// readNumber reads a sequence of digits as a number
func (s *Scanner) readNumber() string {
	position := s.position
	for isDigit(s.ch) {
		s.readChar()
	}
	return s.input[position:s.position]
}
