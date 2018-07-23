package rsql

import (
	"unicode/utf8"
	"unicode"
)

type Lexer struct {
	input string
	position int
	readPosition int
	rune rune
}

func NewLexer(in string) *Lexer {
	l := &Lexer{input: in}
	l.readRune()
	return l
}

func (l *Lexer) readRune() {
	var size int
	if l.readPosition >= len(l.input) {
		l.rune = 0
	} else {
		l.rune, size = utf8.DecodeRuneInString(l.input[l.readPosition:])
	}

	l.position = l.readPosition
	l.readPosition += size
}

func (l *Lexer) peekRune() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
		return r
	}
}

func (l *Lexer) skipSpace() {
	for unicode.IsSpace(l.rune) {
		l.readRune()
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipSpace()

	switch l.rune {
	case '-':
		l.readRune()
		if isNumber(l.rune) {
			tok.Type = NUMBER
			tok.Literal = "-" + l.readNumber()

			return tok
		} else if isLetter(l.rune) {
			tok.Literal = l.readIdentifier()
			tok.Type = "-" + lookupIdent(tok.Literal)

			return tok
		} else {
			tok = newToken(ILLEGAL, l.rune)
		}
	case '=':
		r := l.peekRune()
		if r == '=' {
			r := l.rune
			l.readRune()
			tok.Type, tok.Literal = EQ, string(r)+string(l.rune)
		} else if isLetter(r) {
			tok = l.readFIQL()

			return tok
		} else {
			tok = newToken(EQUAL, l.rune)
		}
	case ';': tok = newToken(SEMICOLON, l.rune)
	case ',': tok = newToken(COMMA, l.rune)
	case ')': tok = newToken(RPAREN, l.rune)
	case '(': tok = newToken(LPAREN, l.rune)
	case '~': tok = newToken(TILDA, l.rune)
	case '!':
		if l.peekRune() == '=' {
			r := l.rune
			l.readRune()
			tok.Type, tok.Literal = NEQ, string(r) + string(l.rune)
		} else {
			tok = newToken(BANG, l.rune)
		}
	case '<':
		if l.peekRune() == '=' {
			r := l.rune
			l.readRune()
			tok.Type, tok.Literal = LE, string(r) + string(l.rune)
		} else {
			tok = newToken(LT, l.rune)
		}
	case '>':
		if l.peekRune() == '=' {
			r := l.rune
			l.readRune()
			tok.Type, tok.Literal = GE, string(r) + string(l.rune)
		} else {
			tok = newToken(GT, l.rune)
		}
	case '"', '\'':
		tok.Type = STR
		tok.Literal = l.readString(l.rune)
	case '\\': tok = newToken(ESCAPE, l.rune)
	case 0: tok.Type, tok.Literal = EOF, ""
	default:
		if isNumber(l.rune) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()

			return tok
		} else if isLetter(l.rune) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)

			return tok
		} else {
			tok = newToken(ILLEGAL, l.rune)
		}
	}

	l.readRune()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.rune) {
		l.readRune()
	}

	return l.input[position:l.position]
}
func (l *Lexer) readNumber() string {
	position := l.position
	for isNumber(l.rune) {
		l.readRune()
	}

	return l.input[position:l.position]
}
func (l *Lexer) readString(st rune) string {
	position := l.position + 1
	escape := false
	for {
		l.readRune()

		if l.rune == '\\' {
			escape = true

			continue
		}
		if escape && l.rune == st {
			escape = false
			continue
		}
		if l.rune == st || l.rune == 0 {
			break
		}
	}

	return l.input[position:l.position]

}
func (l *Lexer) readFIQL() Token {
	l.readRune()
	s := l.readIdentifier()
	if l.rune != '=' {
		return Token{ILLEGAL, "=" + s }
	}
	l.readRune()

	switch s {
	case "eq": return Token{EQ, "=" + s + "="}
	case "neq": return Token{NEQ, "=" + s + "="}
	case "gt": return Token{GT, "=" + s + "="}
	case "ge": return Token{GE, "=" + s + "="}
	case "lt": return Token{LT, "=" + s + "="}
	case "le": return Token{LE, "=" + s + "="}
	default:
		return Token{EXTEND, "=" + s + "=" }
	}
}

func isLetter(r rune) bool {
	if unicode.IsSpace(r) {
		return false
	}
	switch r {
	case '"', '\'', '(', ')', ';', ',', '=', '!', '~', '<', '>':
		return false
	}
	if '0' < r && r < '9' || r == '.' {
		return false
	}

	return unicode.IsPrint(r)
}

func isNumber(r rune) bool {
	if '0' <= r && r <= '9' || r == '.' {
		return true
	}
	return false
}