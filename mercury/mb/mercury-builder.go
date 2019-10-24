package mb

import "sour.is/x/toolbox/mercury"

type option interface {
	applySpace(*mercury.Space)
	applyValue(*mercury.Value)
}

type fragment struct {
	tags   []string
	notes  []string
	values []string
	list   []mercury.Value
}

func (o fragment) applySpace(m *mercury.Space) {
	if o.tags != nil {
		m.Tags = append(m.Tags, o.tags...)
	}
	if o.notes != nil {
		m.Notes = append(m.Notes, o.notes...)
	}
	if o.list != nil {
		m.List = append(m.List, o.list...)
	}
}

func (o fragment) applyValue(m *mercury.Value) {
	if o.tags != nil {
		m.Tags = append(m.Tags, o.tags...)
	}
	if o.notes != nil {
		m.Notes = append(m.Notes, o.notes...)
	}
	if o.values != nil {
		m.Values = append(m.Values, o.values...)
	}
}

// Config builds a config from spaces
func Config(spaces ...mercury.Space) mercury.ArraySpace {
	return spaces
}

func NewSpace(space string, opts ...option) mercury.Space {
	o := mercury.Space{}

	o.Space = space
	for i := range opts {
		opts[i].applySpace(&o)
	}

	return o
}

func WithNotes(n ...string) option {
	return fragment{notes: n}
}

func WithTags(n ...string) option {
	return fragment{tags: n}
}

func WithValue(n ...string) option {
	return fragment{values: n}
}

func NewKey(name string, opts ...option) option {
	o := mercury.Value{}
	o.Name = name
	for i := range opts {
		opts[i].applyValue(&o)
	}

	return fragment{list: []mercury.Value{o}}
}

