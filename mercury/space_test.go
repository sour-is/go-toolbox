package mercury

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArraySpace_String(t *testing.T) {
	tests := []struct {
		value ArraySpace
		want  string
	}{
		{
			[]Space{
				{
					Space: "space",
					Tags:  []string{},
					Notes: []string{},
					List: []Value{
						{
							Seq:    1,
							Name:   "name",
							Tags:   []string{"tag"},
							Values: []string{"value"},
							Notes:  []string{},
						},
					},
				},
			},
			"@space\nname tag  :value\n\n",
		},

		{
			[]Space{
				{
					Space: "space",
					Tags:  []string{},
					Notes: []string{},
					List: []Value{
						{
							Seq:    1,
							Name:   "name",
							Tags:   []string{},
							Values: []string{"value"},
							Notes:  []string{},
						},
					},
				},
			},
			"@space\nname   :value\n\n",
		},

		{
			[]Space{
				{
					Space: "space",
					Notes: []string{"notes"},
					Tags:  []string{},
				},
			},
			"# notes\n@space\n\n",
		},

		{
			[]Space{
				{
					Space: "space",
					Notes: []string{"notes"},
					Tags:  []string{"tag"},
				},
			},
			"# notes\n@space tag\n\n",
		},
	}
	Convey("ArraySpace to string", t, func() {
		for _, tt := range tests {
			result := tt.value.String()
			So(result, ShouldResemble, tt.want)

			in := strings.NewReader(result)
			reverse, err := parseText(in)
			So(err, ShouldBeNil)
			So(reverse.ToArray(), ShouldResemble, tt.value)
		}
	})
}
