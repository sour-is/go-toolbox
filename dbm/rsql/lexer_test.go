package rsql

import (
	"testing"
)

func TestReservedToken(t *testing.T) {
	input := `\();,=!~<>`
	tests := []struct{
		expectedType TokenType
		expectedLiteral string
	} {
		{ESCAPE, `\`},
		{LPAREN, "("},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{COMMA, ","},
		{EQUAL, "="},
		{BANG, "!"},
		{TILDA,"~"},
		{LT, "<"},
		{GT,">"},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%v]: token type wrong. expected=%v got=%v", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%v]: token literal type wrong. expected=%v got=%v", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
func TestNextToken(t *testing.T) {
	input := `director=='name\'s';actor=eq="name's";Year=le=2000,Year>=2010;(one <= -1.0, two != true),three=in=(1,2,3)`
	tests := []struct{
		expectedType TokenType
		expectedLiteral string
	} {
		{IDENT, `director`},
		{EQ, `==`},
		{STR, `name\'s`},
		{SEMICOLON, `;`},
		{IDENT, `actor`},
		{EQ, `=eq=`},
		{STR,`name's`},
		{SEMICOLON, `;`},
		{IDENT, "Year"},
		{LE, "=le="},
		{NUMBER, "2000"},
		{COMMA, ","},
		{IDENT, "Year"},
		{GE, ">="},
		{NUMBER,"2010"},
		{SEMICOLON, ";"},
		{LPAREN, "("},
		{IDENT, "one"},
		{LE, "<="},
		{NUMBER, "-1.0"},
		{COMMA, ","},
		{IDENT, "two"},
		{NEQ, "!="},
		{TRUE, "true"},
		{RPAREN, ")"},
		{COMMA,","},
		{IDENT,"three"},
		{EXTEND,"=in="},
		{LPAREN, "("},
		{NUMBER, "1"},
		{COMMA, ","},
		{NUMBER, "2"},
		{COMMA, ","},
		{NUMBER, "3"},
		{RPAREN, ")"},
	}

	l := NewLexer(input)

	c := 0
	for i, tt := range tests {
		c++
		tok := l.NextToken()
		t.Log(tok, c)
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%v]: token type wrong. expected=%v got=%v", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%v]: token literal type wrong. expected=%v got=%v", i, tt.expectedLiteral, tok.Literal)
		}
	}
	if c != len(tests) {
		t.Errorf("tests[x]: lexer returned different number of tokens. expected=%v got=%v", len(tests), c)
		t.Error(tests)
	}


}