package rsql

// Tokens for RSQL FIQL
const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	IDENT = "IDENT"
	NUMBER = "NUMBER"
	STR = "STR"
	FLOAT = "FLOAT"
	BOOL = "BOOL"
	EXTEND = "EXTEND"

	AND = "AND"
	SEMICOLON = ";"

	OR = "OR"
	COMMA = ","

	LPAREN = "("
	RPAREN = ")"

	TILDA = "~"
	BANG = "!"
	EQUAL = "="
	LT = "<"
	GT = ">"
	LE = "<="
	GE = ">="
	EQ = "=="
	NEQ = "!="
    MINUS = "-"

	SQUOT = "'"
	DQUOT = `"`

	ESCAPE = `\`

	TRUE =  "true"
	FALSE = "false"
	NULL = "null"
)

var keywords = map[string]TokenType {
	"true": TRUE,
	"false": FALSE,
	"null": NULL,
}

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

func newToken(tokenType TokenType, ch rune) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}