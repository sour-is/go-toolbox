package squirrel

import (
	"sour.is/x/toolbox/dbm/rsql"
	"github.com/Masterminds/squirrel"
	"log"
	"sour.is/x/toolbox/dbm"
	"strings"
)

type errors []string
func (e errors) Error() string {
	return strings.Join(e, ",\n")
}

func Query(in string, db dbm.DbInfo) (interface{}, []string) {
	d := decoder{dbInfo: db}
	l := rsql.NewLexer(in)
	p := rsql.NewParser(l)
	program := p.ParseProgram()
	log.Print(program.String())
	return d.decode(program)
}

type decoder struct{
	dbInfo dbm.DbInfo
	errors []string
}

func (db *decoder) decode(in *rsql.Program) (squirrel.Sqlizer, []string) {

	switch len(in.Statements) {
	case 0:
		return nil, db.errors
	case 1:
		return  db.decodeStatement(in.Statements[0]), db.errors
	default:
		a := squirrel.And{}
		for _, stmt := range in.Statements {
			a = append(a, db.decodeStatement(stmt))
		}
		return a, db.errors
	}
}

func  (db *decoder) decodeStatement(in rsql.Statement) squirrel.Sqlizer {
	switch s := in.(type) {
	case *rsql.ExpressionStatement:
		return db.decodeExpression(s.Expression)
	}
	return nil
}

func  (db *decoder) decodeExpression(in rsql.Expression) squirrel.Sqlizer {
	switch e := in.(type) {
	case *rsql.InfixExpression:
		return db. decodeInfix(e)
	}
	return nil
}

func  (db *decoder) decodeInfix(in *rsql.InfixExpression) squirrel.Sqlizer {
	defer func(){
		if r := recover(); r != nil {

		}
	}()

	switch in.Token.Type {
	case rsql.TokAND:
		a := squirrel.And{}
		left := db.decodeExpression(in.Left)
		switch v := left.(type) {
		case squirrel.And:
			for _, el := range v {
				a = append(a, el)
			}
		default:
			a = append(a, v)
		}

		right := db.decodeExpression(in.Right)
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
			db.decodeExpression(in.Left),
			db.decodeExpression(in.Right),
		}
	case rsql.TokEQ:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			db.errors = append(db.errors, err.Error())
			return nil
		}

		return squirrel.Eq{col: db.decodeValue(in.Right)}
	case rsql.TokNEQ:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			db.errors = append(db.errors, err.Error())
			return nil
		}

		return squirrel.NotEq{col: db.decodeValue(in.Right)}
	case rsql.TokGT:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			db.errors = append(db.errors, err.Error())
			return nil
		}

		return squirrel.Gt{col: db.decodeValue(in.Right)}
	case rsql.TokGE:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			db.errors = append(db.errors, err.Error())
			return nil
		}

		return squirrel.GtOrEq{col: db.decodeValue(in.Right)}
	case rsql.TokLT:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			db.errors = append(db.errors, err.Error())
			return nil
		}

		return squirrel.Lt{col: db.decodeValue(in.Right)}
	case rsql.TokLE:
		col, err := db.dbInfo.Col(in.Left.String())
		if err != nil {
			db.errors = append(db.errors, err.Error())
			return nil
		}

		return squirrel.LtOrEq{col: db.decodeValue(in.Right)}
	default:
		return nil
	}
}
func (db *decoder) decodeValue(in rsql.Expression) interface{} {
	switch v := in.(type) {
	case *rsql.Array:
		var values []interface{}
		for _, el := range v.Elements {
			values = append(values, db.decodeValue(el))
		}

		return values
	case *rsql.InfixExpression:
		return db.decodeInfix(v)
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
