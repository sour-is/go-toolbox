package app

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"sour.is/x/toolbox/dbm/rsql"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/ident/session"
	"sour.is/x/toolbox/mercury"
)

func Test_appConfig_GetRules(t *testing.T) {

	type args struct {
		u ident.Ident
	}
	tests := []struct {
		name    string
		args    args
		wantLis mercury.Rules
	}{
		{"normal", args{ident.NullUser{Active: false}}, nil},
		{
			"admin",
			args{
				ident.Ident(
					session.User{
						Active: true,
						Roles:  map[string]struct{}{"admin": struct{}{}},
					},
				),
			},
			mercury.Rules{
				{
					Role:  "read",
					Type:  "NS",
					Match: "app.settings",
				},
				{
					Role:  "read",
					Type:  "NS",
					Match: "app.priority",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := appConfig{}
			if gotLis := a.GetRules(tt.args.u); !reflect.DeepEqual(gotLis, tt.wantLis) {
				t.Errorf("appConfig.GetRules() = %v, want %v", gotLis, tt.wantLis)
			}
		})
	}
}

func Test_appConfig_GetIndex(t *testing.T) {
	type args struct {
		search mercury.NamespaceSearch
		in1    *rsql.Program
	}
	tests := []struct {
		name    string
		args    args
		wantLis mercury.ArraySpace
	}{
		{"nil", args{
			nil,
			nil,
		}, nil},

		{"app.settings", args{
			mercury.ParseNamespace("app.settings"),
			nil,
		}, mercury.ArraySpace{mercury.Space{Space: "app.settings"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := appConfig{}
			if gotLis := a.GetIndex(tt.args.search, tt.args.in1); !reflect.DeepEqual(gotLis, tt.wantLis) {
				t.Errorf("appConfig.GetIndex() = %#v, want %#v", gotLis, tt.wantLis)
			}
		})
	}
}

func Test_appConfig_GetObjects(t *testing.T) {
	type args struct {
		search mercury.NamespaceSearch
		in1    *rsql.Program
		in2    []string
	}
	tests := []struct {
		name    string
		args    args
		wantLis mercury.ArraySpace
	}{
		{"nil", args{
			nil,
			nil,
			nil,
		}, nil},

		{"app.settings", args{
			mercury.ParseNamespace("app.settings"),
			nil,
			nil,
		}, mercury.ArraySpace{mercury.Space{Space: "app.settings", List: []mercury.Value{{Name: "app.setting", Values: []string{"TRUE"}}}}}},
	}
	for _, tt := range tests {
		viper.Set("app.setting", "TRUE")
		t.Run(tt.name, func(t *testing.T) {
			a := appConfig{}
			if gotLis := a.GetObjects(tt.args.search, tt.args.in1, tt.args.in2); !reflect.DeepEqual(gotLis, tt.wantLis) {
				t.Errorf("appConfig.GetIndex() = %#v, want %#v", gotLis, tt.wantLis)
			}
		})
	}
}
