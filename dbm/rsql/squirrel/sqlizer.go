package squirrel

import (
	"sour.is/x/toolbox/dbm/rsql"
	"github.com/Masterminds/squirrel"
	"sour.is/x/toolbox/dbm"
	"fmt"
	"strings"
)

func Query(in string, db dbm.DbInfo) (interface{}, error) {
	d := decoder{dbInfo: db}
	l := rsql.NewLexer(in)
	p := rsql.NewParser(l)
	program := p.ParseProgram()
	return d.decode(program)
}

type decoder struct{
	dbInfo dbm.DbInfo
}
func (db *decoder) decode(in *rsql.Program) (squirrel.Sqlizer, error) {

	switch len(in.Statements) {
	case 0:
		return nil, nil
	case 1:
		return db.decodeStatement(in.Statements[0])
	default:
		a := squirrel.And{}
		for _, stmt := range in.Statements {
			d, err := db.decodeStatement(stmt)
			if d == nil {
				return nil, err
			}

			a = append(a, d)
		}
		return a, nil
	}
}
func (db *decoder) decodeStatement(in rsql.Statement) (squirrel.Sqlizer, error) {
	switch s := in.(type) {
	case *rsql.ExpressionStatement:
		return db.decodeExpression(s.Expression)
	}
	return nil, nil
}
func (db *decoder) decodeExpression(in rsql.Expression) (squirrel.Sqlizer, error) {
	switch e := in.(type) {
	case *rsql.InfixExpression:
		return db. decodeInfix(e)
	}
	return nil, nil
}
func (db *decoder) decodeInfix(in *rsql.InfixExpression) (squirrel.Sqlizer, error) {

	switch in.Token.Type {
	case rsql.TokAND:
		a := squirrel.And{}
		left, err := db.decodeExpression(in.Left)
		if err != nil {
			return nil, err
		}

		switch v := left.(type) {
		case squirrel.And:
			for _, el := range v {
				if el != nil {
					a = append(a, el)
				}
			}

		default:
			if v != nil {
				a = append(a, v)
			}
		}

		right, err := db.decodeExpression(in.Right)
		if err != nil {
			return nil, err
		}

		switch v := right.(type) {
		case squirrel.And:
			for _, el := range v {
				if el != nil {
					a = append(a, el)
				}
			}

		default:
			if v != nil {
				a = append(a, v)
			}
		}

		return a, nil
	case rsql.TokOR:
		left, err := db.decodeExpression(in.Left)
		if err != nil {
			return nil, err
		}
		right, err := db.decodeExpression(in.Right)
		if err != nil {
			return nil, err
		}

		return squirrel.Or{left, right}, nil
	case rsql.TokEQ:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			return nil, err
		}
		v, err := db.decodeValue(in.Right)
		if err != nil {
			return nil, err
		}

		return squirrel.Eq{col: v}, nil
	case rsql.TokLIKE:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			return nil, err
		}

		v, err := db.decodeValue(in.Right)
		if err != nil {
			return nil, err
		}

		switch value := v.(type) {
		case string:
			return Like{col, strings.Replace(value, "*", "%", -1)}, nil
		default:
			return nil, fmt.Errorf("LIKE requires a string value")
		}

	case rsql.TokNEQ:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			return nil, err
		}
		v, err := db.decodeValue(in.Right)
		if err != nil {
			return nil, err
		}

		return squirrel.NotEq{col: v}, nil
	case rsql.TokGT:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			return nil, err
		}
		v, err := db.decodeValue(in.Right)
		if err != nil {
			return nil, err
		}

		return squirrel.Gt{col: v}, nil
	case rsql.TokGE:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			return nil, err
		}
		v, err := db.decodeValue(in.Right)
		if err != nil {
			return nil, err
		}

		return squirrel.GtOrEq{col: v}, nil
	case rsql.TokLT:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			return nil, err
		}
		v, err := db.decodeValue(in.Right)
		if err != nil {
			return nil, err
		}

		return squirrel.Lt{col: v}, nil
	case rsql.TokLE:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			return nil, err
		}
		v, err := db.decodeValue(in.Right)
		if err != nil {
			return nil, err
		}

		return squirrel.LtOrEq{col: v}, nil
	default:
		return nil, nil
	}
}
func (db *decoder) decodeValue(in rsql.Expression) (interface{}, error) {
	switch v := in.(type) {
	case *rsql.Array:
		var values []interface{}
		for _, el := range v.Elements {
			v, err := db.decodeValue(el)
			if err != nil {
				return nil, err
			}

			values = append(values, v)
		}

		return values, nil
	case *rsql.InfixExpression:
		return db.decodeInfix(v)
	case *rsql.Identifier:
		return v.Value, nil
	case *rsql.Integer:
		return v.Value, nil
	case *rsql.Float:
		return v.Value, nil
	case *rsql.String:
		return v.Value, nil
	case *rsql.Bool:
		return v.Value, nil
	case *rsql.Null:
		return nil, nil
	}

	return nil, nil
}

type Like struct{
	column string
	value string
}
func (l Like) ToSql() (sql string, args []interface{}, err error) {
	sql = l.column
	args = append(args, l.value)
	return
}
