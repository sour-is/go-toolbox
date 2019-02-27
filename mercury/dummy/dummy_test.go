package dummy

import (
	"reflect"
	"testing"

	"sour.is/x/toolbox/dbm/rsql"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/mercury"
)

func TestIndexDummy_GetIndex(t *testing.T) {
	type args struct {
		in0 mercury.NamespaceSearch
		in1 *rsql.Program
	}
	tests := []struct {
		name string
		args args
		want mercury.ArraySpace
	}{
		{"default", args{nil, nil}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := IndexDummy{}
			if got := i.GetIndex(tt.args.in0, tt.args.in1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexDummy.GetIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectsDummy_GetObjects(t *testing.T) {
	type args struct {
		in0 mercury.NamespaceSearch
		in1 *rsql.Program
		in2 []string
	}
	tests := []struct {
		name string
		args args
		want mercury.ArraySpace
	}{
		{"default", args{nil, nil, nil}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := ObjectsDummy{}
			if got := o.GetObjects(tt.args.in0, tt.args.in1, tt.args.in2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectsDummy.GetObjects() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteDummy_WriteObjects(t *testing.T) {
	type args struct {
		in0 mercury.ArraySpace
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"default", args{nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WriteDummy{}
			if err := w.WriteObjects(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("WriteDummy.WriteObjects() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRulesDummy_GetRules(t *testing.T) {
	type args struct {
		in0 ident.Ident
	}
	tests := []struct {
		name string
		args args
		want mercury.Rules
	}{
		{"default", args{ident.NullUser{}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RulesDummy{}
			if got := r.GetRules(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RulesDummy.GetRules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupsDummy_GetGroups(t *testing.T) {
	type args struct {
		in0 ident.Ident
	}
	tests := []struct {
		name string
		args args
		want mercury.Rules
	}{
		{"default", args{ident.NullUser{}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GroupsDummy{}
			if got := g.GetGroups(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupsDummy.GetGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotifyDummy_GetNotify(t *testing.T) {
	type args struct {
		in0 string
	}
	tests := []struct {
		name string
		args args
		want mercury.ListNotify
	}{
		{"default", args{""}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NotifyDummy{}
			if got := n.GetNotify(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotifyDummy.GetNotify() = %v, want %v", got, tt.want)
			}
		})
	}
}
