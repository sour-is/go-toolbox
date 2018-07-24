package squirrel

import (
	"testing"
	"github.com/Masterminds/squirrel"
	"fmt"
)

func TestQuery(t *testing.T) {
	tests := []struct{
		input string
		expect squirrel.Sqlizer
	}{
		{"foo==[1, 2, 3]", squirrel.Eq{"foo": []interface{}{1,2,3}}},
		{"foo==1,(bar==2;baz==3)", squirrel.Or{squirrel.Eq{"foo": 1}, squirrel.And{squirrel.Eq{"bar": 2}, squirrel.Eq{"baz": 3}}}},

		{"foo==1",   squirrel.Eq{"foo": 1}},
		{"foo!=1.1",   squirrel.NotEq{"foo": 1.1}},
		{"foo==true",squirrel.Eq{"foo": true}},
		{"foo==null",squirrel.Eq{"foo": nil}},
		{"foo>2",squirrel.Gt{"foo": 2}},
		{"foo>=2.1",squirrel.GtOrEq{"foo": 2.1}},
		{"foo<3",squirrel.Lt{"foo": 3}},
		{"foo<=3.1",squirrel.LtOrEq{"foo": 3.1}},

		{"foo=eq=1",   squirrel.Eq{"foo": 1}},
		{"foo=neq=1.1",   squirrel.NotEq{"foo": 1.1}},
		{"foo=gt=2",squirrel.Gt{"foo": 2}},
		{"foo=ge=2.1",squirrel.GtOrEq{"foo": 2.1}},
		{"foo=lt=3",squirrel.Lt{"foo": 3}},
		{"foo=le=3.1",squirrel.LtOrEq{"foo": 3.1}},

		{"foo==1;bar==2", squirrel.And{squirrel.Eq{"foo": 1}, squirrel.Eq{"bar": 2}}},
		{"foo==1,bar==2", squirrel.Or{squirrel.Eq{"foo": 1}, squirrel.Eq{"bar": 2}}},
		{"foo==1,bar==2;baz=3", squirrel.And{squirrel.Or{squirrel.Eq{"foo": 1}, squirrel.Eq{"bar": 2}}, squirrel.Eq{"baz": 3}}},
		{
			input:`director=='name\'s';actor=eq="name\'s";Year=le=2000,Year>=2010;one <= -1.0, two != true`,
			expect:
				squirrel.And{
					squirrel.Eq{"director": "name's"},
					squirrel.Eq{"actor": "name's"},
					squirrel.Or{
						squirrel.LtOrEq{"Year": 2000},
						squirrel.GtOrEq{"Year": 2010},
					},
					squirrel.Or{
						squirrel.LtOrEq{"one": -1.0},
						squirrel.NotEq{"two": true},
					},
				},
		},

		{
			`genres==[sci-fi,action] ; genres==[romance,animated,horror] , director==Que*Tarantino`,
			squirrel.And{
				squirrel.Eq{"genres": []interface{}{"sci-fi", "action"}},
				squirrel.Or{
					squirrel.Eq{"genres": []interface{}{"romance","animated","horror"}},
					squirrel.Eq{"director": "Que*Tarantino"},
				},
			},
		},
		{"", nil},

	}

	for i, tt := range tests {
		actual := fmt.Sprintf("%#v", Query(tt.input))
		expect := fmt.Sprintf("%#v", tt.expect)
		if expect != actual {
			t.Errorf("test[%d]: %v\n\tinput and expected are not the same. wanted=%v got=%v", i, tt.input, expect, actual)
		}
	}

}