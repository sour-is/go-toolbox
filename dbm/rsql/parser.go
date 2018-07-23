package rsql

type (
	prefixParseFn func() Expression
	infixParseFn func(expression Expression) Expression
)

type Parser struct {
	l *Lexer

	curToken Token
	peekToken Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

