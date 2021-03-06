package token

type Token struct {
	Type    TokenType
	Literal string
}

type TokenType string

const (
	ASSIGN   TokenType = "="
	PLUS               = "+"
	MINUS              = "-"
	BANG               = "!"
	ASTERISK           = "*"
	SLASH              = "/"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	COMMA     = ","
	SEMICOLON = ";"

	LT     = "<"
	GT     = ">"
	LE     = "<="
	GE     = ">="
	EQ     = "=="
	NOT_EQ = "!="

	String     = "STRING"
	Integer    = "INTEGER"
	Identifier = "IDENTIFIER"
	Let        = "LET"
	IF         = "IF"
	Else       = "ELSE"
	Return     = "RETURN"
	Switch     = "SWITCH"
	FUN        = "FUN"
	True       = "True"
	False      = "False"

	EOF     = "EOF"
	INVALID = "INVALID"
)

var Keywords = map[string]TokenType{
	"let":    Let,
	"if":     IF,
	"else":   Else,
	"return": Return,
	"switch": Switch,
	"fun":    FUN,
	"true":   True,
	"false":  False,
}

func NewToken(t TokenType, ch byte) Token {
	return Token{Type: t, Literal: string(ch)}
}
