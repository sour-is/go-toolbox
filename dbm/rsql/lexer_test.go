package rsql

import (
	"testing"
)

func TestReservedToken(t *testing.T) {
	input := `( ) ; , == != ~ < > <= >= [ ]`
	tests := []struct{
		expectedType TokenType
		expectedLiteral string
	} {
		{TokLParen,   "("},
		{TokRParen,   ")"},
		{TokAND,      ";"},
		{TokOR,       ","},
		{TokEQ,       "=="},
		{TokNEQ,      "!="},
		{TokLIKE,     "~"},
		{TokLT,       "<"},
		{TokGT,       ">"},
		{TokLE,       "<="},
		{TokGE,       ">="},
		{TokLBracket, "["},
		{TokRBracket, "]"},
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
		{TokIdent, `director`},
		{TokEQ, `==`},
		{TokString, `name\'s`},
		{TokAND, `;`},
		{TokIdent, `actor`},
		{TokEQ, `=eq=`},
		{TokString,`name's`},
		{TokAND, `;`},
		{TokIdent, "Year"},
		{TokLE, "=le="},
		{TokInteger, "2000"},
		{TokOR, ","},
		{TokIdent, "Year"},
		{TokGE, ">="},
		{TokInteger,"2010"},
		{TokAND, ";"},
		{TokLParen, "("},
		{TokIdent, "one"},
		{TokLE, "<="},
		{TokFloat, "-1.0"},
		{TokOR, ","},
		{TokIdent, "two"},
		{TokNEQ, "!="},
		{TokTRUE, "true"},
		{TokRParen, ")"},
		{TokOR,","},
		{TokIdent,"three"},
		{TokExtend,"=in="},
		{TokLParen, "("},
		{TokInteger, "1"},
		{TokOR, ","},
		{TokInteger, "2"},
		{TokOR, ","},
		{TokInteger, "3"},
		{TokRParen, ")"},
	}

	l := NewLexer(input)

	c := 0
	for i, tt := range tests {
		c++
		tok := l.NextToken()

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