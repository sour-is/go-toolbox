package mercury

import (
	"fmt"
	"strings"

	"sour.is/x/toolbox/ident"
)

// Space stores a registry of spaces
type Space struct {
	ID    uint64   `json:"-"`
	Space string   `json:"space"`
	Tags  []string `json:"tags"`
	Notes []string `json:"notes"`
	List  []Value  `json:"list"`
}

// HasTag returns true if needle is found
// If the needle ends with a / it will be treated
// as a prefix for tag meta data.
func (s Space) HasTag(needle string) bool {
	isPrefix := strings.HasSuffix(needle, "/")
	for i := range s.Tags {
		switch isPrefix {
		case true:
			if strings.HasPrefix(s.Tags[i], needle) {
				return true
			}
		case false:
			if s.Tags[i] == needle {
				return true
			}
		}
	}
	return false
}

// GetTagMeta retuns the value after a '/' in a tag.
// Tags are in the format 'name' or 'name/value'
// This function returns the value.
func (s Space) GetTagMeta(needle string, offset int) string {
	if !strings.HasSuffix(needle, "/") {
		needle += "/"
	}

	for i := range s.Tags {
		if strings.HasPrefix(s.Tags[i], needle) {
			if offset > 0 {
				offset--
				continue
			}

			return strings.TrimPrefix(s.Tags[i], needle)
		}
	}

	return ""
}

// FirstTagMeta returns the first meta tag value.
func (s Space) FirstTagMeta(needle string) string {
	return s.GetTagMeta(needle, 0)
}

// GetValues that match name
func (s Space) GetValues(name string) (lis []Value) {
	for i := range s.List {
		if s.List[i].Name == name {
			lis = append(lis, s.List[i])
		}
	}
	return
}

// FirstValue that matches name
func (s Space) FirstValue(name string) Value {
	for i := range s.List {
		if s.List[i].Name == name {
			return s.List[i]
		}
	}
	return Value{}
}

type Config []*Space

func NewConfig(spaces ...*Space) Config {
	return spaces
}
func (c *Config) AddSpace(spaces ...*Space) *Config {
	*c = append(*c, spaces...)
	return c
}
func NewSpace(space string) *Space {
	return &Space{Space: space}
}
func (s *Space) SetTags(tags ...string) *Space {
	s.Tags = tags
	return s
}
func (s *Space) AddTags(tags ...string) *Space {
	s.Tags = append(s.Tags, tags...)
	return s
}
func (s *Space) SetNotes(notes ...string) *Space {
	s.Notes = notes
	return s
}
func (s *Space) AddNotes(notes ...string) *Space {
	s.Notes = append(s.Notes, notes...)
	return s
}
func (s *Space) SetKeys(keys ...*Value) *Space {
	for i := range keys {
		k := *keys[i]
		k.Seq = uint64(i)
		s.List = append(s.List, k)
	}

	return s
}
func (s *Space) AddKeys(keys ...*Value) *Space {
	l := uint64(len(s.List))
	for i := range keys {
		k := *keys[i]
		k.Seq = uint64(i) + l
		s.List = append(s.List, k)
	}
	return s
}

// SpaceMap generic map of space values
type SpaceMap map[string]*Space

// Len implements Len for sort.interface
func (lis Config) Len() int {
	return len(lis)
}

// Less implements Less for sort.interface
func (lis Config) Less(i, j int) bool {
	return lis[i].Space < lis[j].Space
}

// Swap implements Swap for sort.interface
func (lis Config) Swap(i, j int) { lis[i], lis[j] = lis[j], lis[i] }

// StringList returns the space names as a list
func (lis Config) StringList() string {
	var buf strings.Builder
	for _, o := range lis {
		buf.WriteString(o.Space)
		buf.WriteRune('\n')
	}
	return buf.String()
}

// String format config as string
func (lis Config) String() string {
	attLen := 0
	tagLen := 0

	for _, o := range lis {
		for _, v := range o.List {
			l := len(v.Name)
			if attLen <= l {
				attLen = l
			}

			t := len(strings.Join(v.Tags, " "))
			if tagLen <= t {
				tagLen = t
			}
		}
	}

	var buf strings.Builder
	for _, o := range lis {
		if len(o.Notes) > 0 {
			buf.WriteString("# ")
			buf.WriteString(strings.Join(o.Notes, "\n# "))
			buf.WriteRune('\n')
		}

		buf.WriteString("@")
		buf.WriteString(o.Space)
		if len(o.Tags) > 0 {
			buf.WriteRune(' ')
			buf.WriteString(strings.Join(o.Tags, " "))
		}
		buf.WriteRune('\n')

		for _, v := range o.List {
			if len(v.Notes) > 0 {
				buf.WriteString("# ")
				buf.WriteString(strings.Join(v.Notes, "\n# "))
				buf.WriteString("\n")
			}

			buf.WriteString(v.Name)
			buf.WriteString(strings.Repeat(" ", attLen-len(v.Name)+1))

			if len(v.Tags) > 0 {
				t := strings.Join(v.Tags, " ")
				buf.WriteString(t)
				buf.WriteString(strings.Repeat(" ", tagLen-len(t)+1))
			} else {
				buf.WriteString(strings.Repeat(" ", tagLen+1))
			}

			switch len(v.Values) {
			case 0:
				buf.WriteString("\n")
			case 1:
				buf.WriteString(" :")
				buf.WriteString(v.Values[0])
				buf.WriteString("\n")
			default:
				buf.WriteString(" :")
				buf.WriteString(v.Values[0])
				buf.WriteString("\n")
				for _, s := range v.Values[1:] {
					buf.WriteString(strings.Repeat(" ", attLen+tagLen+3))
					buf.WriteString(":")
					buf.WriteString(s)
					buf.WriteString("\n")
				}
			}
		}

		buf.WriteRune('\n')
	}

	return buf.String()
}

// EnvString format config as environ
func (lis Config) EnvString() string {
	var buf strings.Builder
	for _, o := range lis {
		for _, v := range o.List {
			buf.WriteString(o.Space)
			for _, t := range o.Tags {
				buf.WriteRune(' ')
				buf.WriteString(t)
			}
			buf.WriteRune(':')
			buf.WriteString(v.Name)
			for _, t := range v.Tags {
				buf.WriteRune(' ')
				buf.WriteString(t)
			}
			switch len(v.Values) {
			case 0:
				buf.WriteRune('=')
				buf.WriteRune('\n')
			case 1:
				buf.WriteRune('=')
				buf.WriteString(v.Values[0])
				buf.WriteRune('\n')
			default:
				buf.WriteRune('+')
				buf.WriteRune('=')
				buf.WriteString(v.Values[0])
				buf.WriteRune('\n')
				for _, s := range v.Values[1:] {
					buf.WriteString(o.Space)
					buf.WriteRune(':')
					buf.WriteString(v.Name)
					buf.WriteRune('+')
					buf.WriteRune('=')
					buf.WriteString(s)
					buf.WriteRune('\n')
				}
			}
		}
	}

	return buf.String()
}

// IniString format config as ini
func (lis Config) IniString() string {
	var buf strings.Builder
	for _, o := range lis {
		buf.WriteRune('[')
		buf.WriteString(o.Space)
		buf.WriteRune(']')
		buf.WriteRune('\n')
		for _, v := range o.List {
			buf.WriteString(v.Name)
			switch len(v.Values) {
			case 0:
				buf.WriteRune('=')
				buf.WriteRune('\n')
			case 1:
				buf.WriteRune('=')
				buf.WriteString(v.Values[0])
				buf.WriteRune('\n')
			default:
				buf.WriteRune('[')
				buf.WriteRune('0')
				buf.WriteRune(']')

				buf.WriteRune('=')
				buf.WriteString(v.Values[0])
				buf.WriteRune('\n')
				for i, s := range v.Values[1:] {
					buf.WriteString(v.Name)
					buf.WriteRune('[')
					buf.WriteString(fmt.Sprintf("%d", i))
					buf.WriteRune(']')
					buf.WriteRune('=')
					buf.WriteString(s)
					buf.WriteRune('\n')
				}
			}
		}
	}

	return buf.String()
}

func (lis Config) accessFilter(id ident.Ident) (out Config, err error) {
	rules := Registry.GetRules(id)

	accessList := make(map[string]struct{})
	for _, o := range lis {
		if _, ok := accessList[o.Space]; ok {
			out = append(out, o)
			continue
		}

		if role := rules.GetRoles("NS", o.Space); role.HasRole("read", "write") && !role.HasRole("deny") {
			accessList[o.Space] = struct{}{}
			out = append(out, o)
		}
	}

	return
}

func (rules Rules) filterSpace(lis Config) (out Config, err error) {
	accessList := make(map[string]struct{})
	for _, o := range lis {
		if _, ok := accessList[o.Space]; ok {
			out = append(out, o)
			continue
		}

		if role := rules.GetRoles("NS", o.Space); role.HasRole("read", "write") && !role.HasRole("deny") {
			accessList[o.Space] = struct{}{}
			out = append(out, o)
		}
	}

	return
}

func (lis Config) stringArray() []string {
	out := make([]string, 0, len(lis))
	for _, o := range lis {
		out = append(out, o.Space)
	}
	return out
}

// ToSpaceMap formats as SpaceMap
func (lis Config) ToSpaceMap() SpaceMap {
	out := make(SpaceMap)
	for _, c := range lis {
		out[c.Space] = c
	}
	return out
}

// ToArray converts SpaceMap to ArraySpace
func (m SpaceMap) ToArray() Config {
	a := make(Config, 0, len(m))
	for _, s := range m {
		a = append(a, s)
	}
	return a
}

// Value stores the attributes for space values
type Value struct {
	Seq    uint64   `json:"-"`
	Name   string   `json:"name"`
	Values []string `json:"values"`
	Notes  []string `json:"notes"`
	Tags   []string `json:"tags"`
}

// HasTag returns true if needle is found
// If the needle ends with a / it will be treated
// as a prefix for tag meta data.
func (v Value) HasTag(needle string) bool {
	isPrefix := strings.HasSuffix(needle, "/")
	for i := range v.Tags {
		switch isPrefix {
		case true:
			if strings.HasPrefix(v.Tags[i], needle) {
				return true
			}
		case false:
			if v.Tags[i] == needle {
				return true
			}
		}
	}
	return false
}

// GetTagMeta retuns the value after a '/' in a tag.
// Tags are in the format 'name' or 'name/value'
// This function returns the value.
func (v Value) GetTagMeta(needle string, offset int) string {
	if !strings.HasSuffix(needle, "/") {
		needle += "/"
	}

	for i := range v.Tags {
		if strings.HasPrefix(v.Tags[i], needle) {
			if offset > 0 {
				offset--
				continue
			}

			return strings.TrimPrefix(v.Tags[i], needle)
		}
	}

	return ""
}

// FirstTagMeta returns the first meta tag value.
func (v Value) FirstTagMeta(needle string) string {
	return v.GetTagMeta(needle, 0)
}

// First value in array.
func (v Value) First() string {
	if len(v.Values) < 1 {
		return ""
	}

	return v.Values[0]
}

// Join values with newlines.
func (v Value) Join() string {
	return strings.Join(v.Values, "\n")
}

func NewValue(name string) *Value {
	return &Value{Name: name}
}
func (v *Value) SetTags(tags ...string) *Value {
	v.Tags = tags
	return v
}
func (v *Value) AddTags(tags ...string) *Value {
	v.Tags = append(v.Tags, tags...)
	return v
}
func (v *Value) SetNotes(notes ...string) *Value {
	v.Notes = notes
	return v
}
func (v *Value) AddNotes(notes ...string) *Value {
	v.Notes = append(v.Notes, notes...)
	return v
}
func (v *Value) SetValues(values ...string) *Value {
	v.Values = values
	return v
}
func (v *Value) AddValues(values ...string) *Value {
	v.Values = append(v.Values, values...)
	return v
}
