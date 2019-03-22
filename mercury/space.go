package mercury

import (
	"fmt"
	"strings"

	"sour.is/x/toolbox/ident"
)

// Space stores a registry of spaces
type Space struct {
	ID    uint64
	Space string
	Tags  []string
	Notes []string
	List  []Value
}

// SpaceMap generic map of space values
type SpaceMap map[string]Space

// ArraySpace is a list of config spaces
type ArraySpace []Space

// Len implements Len for sort.interface
func (lis ArraySpace) Len() int {
	return len(lis)
}

// Less implements Less for sort.interface
func (lis ArraySpace) Less(i, j int) bool {
	return lis[i].Space < lis[j].Space
}

// Swap implements Swap for sort.interface
func (lis ArraySpace) Swap(i, j int) { lis[i], lis[j] = lis[j], lis[i] }

// StringList returns the space names as a list
func (lis ArraySpace) StringList() string {
	var buf strings.Builder
	for _, o := range lis {
		buf.WriteString(o.Space)
		buf.WriteRune('\n')
	}
	return buf.String()
}

// String format config as string
func (lis ArraySpace) String() string {
	nsLen := 0
	attLen := 0
	tagLen := 0

	for _, o := range lis {
		if nsLen <= len(o.Space) {
			nsLen = len(o.Space)
		}
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
			buf.WriteString("# @")
			buf.WriteString(o.Space)
			if len(o.Tags) > 0 {
				buf.WriteRune(' ')
				buf.WriteString(strings.Join(o.Tags, " "))
			}
		} else if len(o.Tags) > 0 {
			buf.WriteString("# @")
			buf.WriteString(o.Space)
			if len(o.Tags) > 0 {
				buf.WriteRune(' ')
				buf.WriteString(strings.Join(o.Tags, " "))
			}
		}

		if len(o.Notes) > 0 || len(o.Tags) > 0 {
			buf.WriteRune('\n')
			buf.WriteRune('\n')
		}

		for _, v := range o.List {
			if len(v.Notes) > 0 {
				buf.WriteString("# ")
				buf.WriteString(strings.Join(v.Notes, "\n# "))
				buf.WriteString("\n")
			}
			buf.WriteString(o.Space)
			buf.WriteString(strings.Repeat(" ", nsLen-len(o.Space)))

			if len(v.Tags) > 0 {
				buf.WriteRune(' ')
				t := strings.Join(v.Tags, " ")

				buf.WriteString(t)
				buf.WriteString(strings.Repeat(" ", tagLen-len(t)))
			} else {
				buf.WriteString(strings.Repeat(" ", tagLen+1))
			}

			buf.WriteRune(' ')
			buf.WriteString(v.Name)
			buf.WriteString(strings.Repeat(" ", attLen-len(v.Name)))

			switch len(v.Values) {
			case 0:
				buf.WriteString(" :")
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
					buf.WriteString(strings.Repeat(" ", nsLen+attLen+tagLen+3))
					buf.WriteString(":")
					buf.WriteString(s)
					buf.WriteString("\n")
				}
			}
		}
	}

	return buf.String()
}

// EnvString format config as environ
func (lis ArraySpace) EnvString() string {
	var buf strings.Builder
	for _, o := range lis {
		for _, v := range o.List {
			buf.WriteString(o.Space)
			buf.WriteRune(':')
			buf.WriteString(o.Space)
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
func (lis ArraySpace) IniString() string {
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

func (lis ArraySpace) accessFilter(id ident.Ident) (out ArraySpace, err error) {
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

func (rules Rules) filterSpace(lis ArraySpace) (out ArraySpace, err error) {
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

func (lis ArraySpace) stringArray() []string {
	var out []string
	for _, o := range lis {
		out = append(out, o.Space)
	}
	return out
}

// ToSpaceMap formats as SpaceMap
func (lis ArraySpace) ToSpaceMap() SpaceMap {
	out := make(SpaceMap)
	for _, c := range lis {
		out[c.Space] = c
	}
	return out
}

// ToArray converts SpaceMap to ArraySpace
func (m SpaceMap) ToArray() ArraySpace {
	var a ArraySpace
	for _, s := range m {
		a = append(a, s)
	}
	return a
}

// Value stores the attributes for space values
type Value struct {
	Seq    uint64
	Name   string
	Values []string
	Notes  []string
	Tags   []string
}
