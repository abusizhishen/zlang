package lexer

import "zlang/token"

type Lexer struct {
	input        string
	readPosition int
	position     int
	ch           byte
}

func New(string2 string) *Lexer {
	lex := &Lexer{input: string2}
	lex.readChar()
	return lex
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpaces()

	switch l.ch {
	case '+':
		tok = token.Token{token.PLUS, string(l.ch)}
	case '-':
		tok = token.Token{token.MINUS, string(l.ch)}
	case '*':
		tok = token.Token{token.ASTERISK, string(l.ch)}
	case '/':
		tok = token.Token{token.SLASH, string(l.ch)}
	case '=':
		tok = token.Token{token.ASSIGN, string(l.ch)}
	case '(':
		tok = token.Token{token.LPAREN, string(l.ch)}
	case ')':
		tok = token.Token{token.RPAREN, string(l.ch)}
	case '{':
		tok = token.Token{token.LBRACE, string(l.ch)}
	case '}':
		tok = token.Token{token.RBRACE, string(l.ch)}
	case '>':
		if l.peekChar() == '=' {
			tok = token.Token{token.GE, l.input[l.position : l.readPosition+1]}
			l.readChar()
		} else {
			tok = token.Token{token.GT, string(l.ch)}
		}
	case '<':
		if l.peekChar() == '=' {
			tok = token.Token{token.LE, l.input[l.position : l.readPosition+1]}
			l.readChar()
		} else {
			tok = token.Token{token.LT, string(l.ch)}
		}

	case '"':
		tok.Type = token.String
		tok.Literal = l.readString()
	case 0:
		tok = token.Token{Type: token.EOF}
	default:
		if l.isLetter(l.ch) {
			tok.Type = token.Identity
			tok.Literal = l.readIdentify()
			if t, ok := token.Keywords[tok.Literal]; ok {
				tok.Type = t
			}
			return tok
		} else if l.isNumber(l.ch) {
			tok.Type = token.Number
			tok.Literal = l.readNum()
			return tok
		} else {
			tok.Type = token.INVALID
			tok.Literal = string(l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch > +'A' && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readString() string {
	position := l.position
	for l.ch != '"' {
		l.readChar()
	}

	if l.ch == '"' {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readIdentify() string {
	position := l.position
	for l.isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) readNum() string {
	position := l.position
	for l.isNumber(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) Tokens() []token.Token {
	var tokens []token.Token
	tok := l.NextToken()
	for tok.Type != token.EOF {
		tokens = append(tokens, tok)
		tok = l.NextToken()
	}

	return tokens
}

func (l *Lexer) skipWhiteSpaces() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}
