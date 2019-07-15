package mercury

import (
	"io"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseText(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		wantConfig SpaceMap
		wantErr    bool
	}{
		{
			"space",
			`@space`,
			SpaceMap{
				"space": Space{Space: "space", Tags: []string{}, Notes: []string{}},
			},
			false,
		},
		{
			"space tag",
			`@space tag`,
			SpaceMap{
				"space": Space{Space: "space", Tags: []string{"tag"}, Notes: []string{}},
			},
			false,
		},
		{
			"space tag note",
			"# note\n@space tag",
			SpaceMap{
				"space": Space{Space: "space", Tags: []string{"tag"}, Notes: []string{"note"}},
			},
			false,
		},
		{
			"space tag note value",
			"# note\n@space tag\n\nname :value",
			SpaceMap{
				"space": Space{
					Space: "space",
					Tags:  []string{"tag"},
					Notes: []string{"note"},
					List: []Value{
						Value{Seq: 1, Name: "name", Values: []string{"value"}, Tags: []string{}, Notes: []string{}},
					},
				},
			},
			false,
		},
		{
			"space tag note value tag note",
			"# note1\n@space tag1\n\n# note2\nname tag2 :value",
			SpaceMap{
				"space": Space{
					Space: "space",
					Tags:  []string{"tag1"},
					Notes: []string{"note1"},
					List: []Value{
						Value{Seq: 1, Name: "name", Values: []string{"value"}, Tags: []string{"tag2"}, Notes: []string{"note2"}},
					},
				},
			},
			false,
		},
		{
			"space tag note value tag note x2",
			"# note1\n@space tag1\n\n# note2\nname tag2 :value\nname2 :value1\n:value2\n    :value3",
			SpaceMap{
				"space": Space{
					Space: "space",
					Tags:  []string{"tag1"},
					Notes: []string{"note1"},
					List: []Value{
						Value{Seq: 1, Name: "name", Values: []string{"value"}, Tags: []string{"tag2"}, Notes: []string{"note2"}},
						Value{Seq: 2, Name: "name2", Values: []string{"value1", "value2", "value3"}, Tags: []string{}, Notes: []string{}},
					},
				},
			},
			false,
		},
	}

	Convey("Parse Text to SpaceMap", t, func() {
		for _, tt := range tests {
			gotConfig, err := parseText(reader(tt.text))
			if tt.wantErr {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
			}
			So(gotConfig, ShouldResemble, tt.wantConfig)
		}
	})
}

func reader(s string) io.Reader {
	return strings.NewReader(s)
}
