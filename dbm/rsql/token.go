package rsql

// Tokens for RSQL FIQL
const (
	TokIllegal = "TokIllegal"
	TokEOF     = "TokEOF"

	TokIdent   = "TokIdent"
	TokInteger = "TokInteger"
	TokString  = "TokString"
	TokFloat   = "TokFloat"
	TokExtend  = "TokExtend"

	TokLParen = "("
	TokRParen = ")"

	TokLBracket = "["
	TokRBracket = "]"

	TokLIKE = "~"
	TokNOT  = "!"
	TokLT   = "<"
	TokGT   = ">"
	TokLE   = "<="
	TokGE   = ">="
	TokEQ   = "=="
	TokNEQ  = "!="
	TokAND  = ";"
	TokOR   = ","

	TokTRUE  =  "true"
	TokFALSE = "false"
	TokNULL  = "null"
)

var keywords = map[string]TokenType {
	"true":  TokTRUE,
	"false": TokFALSE,
	"null":  TokNULL,
	"and":   TokAND,
	"or":    TokOR,
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
	return TokIdent
}