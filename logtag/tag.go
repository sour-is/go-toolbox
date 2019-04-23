package logtag

import (
	"fmt"
	"strings"
)

// Tags are named datum to add to a log line
type Tags map[string]Tag

// Tag is tag value
type Tag fmt.Stringer
type value string

// String returns a string
func (v value) String() string {
	return string(v)
}

// NewTag formats a value to print.
func NewTag(in interface{}) Tag {
	return value(fmt.Sprintf("%v", in))
}

// MapToTags convert a interface map to Tags
func MapToTags(m map[string]interface{}) (tags Tags) {
	for k, v := range m {
		switch val := v.(type) {
		case string:
			tags[k] = value(val)
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

func readTags(tags ...interface{}) Tags {
	m := make(Tags, len(tags)/2)
	var key string
	for _, v := range tags {
		switch s := v.(type) {
		case string:
			key = s
		default:
			m[key] = NewTag(v)
		}
	}
	return m
}
