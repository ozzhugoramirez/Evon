package compiler

type Lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune
	line         int
	column       int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  []rune(input),
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) newToken(tokenType TokenType, ch rune) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: l.line, Column: l.column}
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()

	currentLine := l.line
	currentCol := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: currentLine, Column: currentCol}
		} else {
			tok = l.newToken(ASSIGN, l.ch)
		}
	case ':':
		tok = l.newToken(COLON, l.ch)
	case '>':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: REDIRECT, Literal: string(ch) + string(l.ch), Line: currentLine, Column: currentCol}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GTE, Literal: string(ch) + string(l.ch), Line: currentLine, Column: currentCol}
		} else {
			tok = l.newToken(GT, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LTE, Literal: string(ch) + string(l.ch), Line: currentLine, Column: currentCol}
		} else {
			tok = l.newToken(LT, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch), Line: currentLine, Column: currentCol}
		} else {
			tok = l.newToken(ILLEGAL, l.ch)
		}
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: ARROW, Literal: string(ch) + string(l.ch), Line: currentLine, Column: currentCol}
		} else {
			tok = l.newToken(MINUS, l.ch)
		}
	case '+':
		tok = l.newToken(PLUS, l.ch)
	case '*':
		tok = l.newToken(ASTERISK, l.ch)
	case '/':
		tok = l.newToken(SLASH, l.ch)
	case '.':
		tok = l.newToken(DOT, l.ch)
	case '(':
		tok = l.newToken(LPAREN, l.ch)
	case ')':
		tok = l.newToken(RPAREN, l.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		tok.Line = currentLine
		tok.Column = currentCol
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = currentLine
		tok.Column = currentCol
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			tok.Line = currentLine
			tok.Column = currentCol
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			tok.Line = currentLine
			tok.Column = currentCol
			return tok
		} else {
			tok = l.newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= 0x00C0
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= 0x00C0
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
