package squirrel

import (
	"sour.is/x/toolbox/dbm/rsql"
	"github.com/Masterminds/squirrel"
	"log"
)

func Query(in string) interface{} {
	l := rsql.NewLexer(in)
	p := rsql.NewParser(l)
	program := p.ParseProgram()
	log.Print(program.String())
	return decode(program)
}

func decode(in *rsql.Program) squirrel.Sqlizer {
	switch len(in.Statements) {
	case 0:
		return nil
	case 1:
		return decodeStatement(in.Statements[0])
	default:
		a := squirrel.And{}
		for _, stmt := range in.Statements {
			a = append(a, decodeStatement(stmt))
		}
		return a
	}
}

func decodeStatement(in rsql.Statement) squirrel.Sqlizer {
	switch s := in.(type) {
	case *rsql.ExpressionStatement:
		return decodeExpression(s.Expression)
	}
	return nil
}

func decodeExpression(in rsql.Expression) squirrel.Sqlizer {
	switch e := in.(type) {
	case *rsql.InfixExpression:
		return decodeInfix(e)
	}
	return nil
}

func decodeInfix(in *rsql.InfixExpression) squirrel.Sqlizer {
	switch in.Token.Type {
	case rsql.TokAND:
		a := squirrel.And{}
		left := decodeExpression(in.Left)
		switch v := left.(type) {
		case squirrel.And:
			for _, el := range v {
				a = append(a, el)
			}
		default:
			a = append(a, v)
		}

		right := decodeExpression(in.Right)
		switch v := right.(type) {
		case squirrel.And:
			for _, el := range v {
				a = append(a, el)
			}
		default:
			a = append(a, v)
		}

		return a
	case rsql.TokOR:
		return squirrel.Or{
			decodeExpression(in.Left),
			decodeExpression(in.Right),
		}
	case rsql.TokEQ:
		return squirrel.Eq{in.Left.String(): decodeValue(in.Right)}
	case rsql.TokNEQ:
		return squirrel.NotEq{in.Left.String(): decodeValue(in.Right)}
	case rsql.TokGT:
		return squirrel.Gt{in.Left.String(): decodeValue(in.Right)}
	case rsql.TokGE:
		return squirrel.GtOrEq{in.Left.String(): decodeValue(in.Right)}
	case rsql.TokLT:
		return squirrel.Lt{in.Left.String(): decodeValue(in.Right)}
	case rsql.TokLE:
		return squirrel.LtOrEq{in.Left.String(): decodeValue(in.Right)}
	default:
		return nil
	}
}
func decodeValue(in rsql.Expression) interface{} {
	switch v := in.(type) {
	case *rsql.Array:
		var values []interface{}
		for _, el := range v.Elements {
			values = append(values, decodeValue(el))
		}

		return values
	case *rsql.InfixExpression:
		return decodeInfix(v)
	case *rsql.Identifier:
		return v.Value
	case *rsql.Integer:
		return v.Value
	case *rsql.Float:
		return v.Value
	case *rsql.String:
		return v.Value
	case *rsql.Bool:
		return v.Value
	case *rsql.Null:
		return nil
	}

	return nil
}