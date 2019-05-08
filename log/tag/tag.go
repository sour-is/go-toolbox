package tag

import (
	"fmt"
	"strings"
)

// Tags are named datum to add to a log line
type Tags map[string]Tag

// Tag is tag value
type Tag fmt.Stringer

// Value implements a tag value stringer.
type Value string

// String returns a string
func (v Value) String() string {
	return string(v)
}

// NewTag formats a value to print.
func NewTag(in interface{}) Tag {
	return Value(fmt.Sprintf("%v", in))
}

// MapToTags convert a interface map to Tags
func MapToTags(m map[string]interface{}) (tags Tags) {
	for k, v := range m {
		switch val := v.(type) {
		case string:
			tags[k] = Value(val)
		case Tag:
			tags[k] = val
		default:
			tags[k] = NewTag(val)
		}
	}
	return
}

// String return a string
func (t Tags) String() string {
	var b strings.Builder
	for k, v := range t {
		b.WriteString(k)
		b.WriteRune('=')
		b.WriteString(fmt.Sprint(v))
		b.WriteRune(' ')
	}

	return b.String()
}

// Set value to key
func (t *Tags) Set(key string, val interface{}) {
	switch v := val.(type) {
	case Tag:
		(*t)[key] = v
	default:
		(*t)[key] = NewTag(v)
	}
}

// ReadTags read in tags from a list of [string, value]
func ReadTags(tags ...interface{}) Tags {
	if len(tags) < 2 {
		return nil
	}

	m := make(Tags, len(tags)/2)
	var key string
	for i, v := range tags {
		if i%2 == 0 {
			key = NewTag(v).String()
		} else {
			m[key] = NewTag(v)
		}
	}

	return m
}
