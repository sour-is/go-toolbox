package gql

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ListIDs is a list of ID values.
type ListIDs []uint64

// Scan implements the Scanner interface for ListIDs.
func (e *ListIDs) Scan(value interface{}) (err error) {
	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	case []rune:
		str = string(v)

	default:
		return fmt.Errorf("array must be uint64, got: %T", value)
	}

	if e == nil {
		*e = ListIDs{}
	}

	for _, s := range splitComma(string(str)) {
		an, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return fmt.Errorf("%d is not a valid uint64", an)
		}
		*e = append(*e, an)
	}

	return nil
}

// Value implements the driver Valuer interface for ListIDs.
func (e ListIDs) Value() (v driver.Value, err error) {
	var b strings.Builder

	_, err = b.WriteString("{")
	if err != nil {
		return
	}

	var arr []string
	for _, s := range e {
		arr = append(arr, strconv.FormatUint(s, 10))
	}
	_, err = b.WriteString(strings.Join(arr, ",") + "}")
	if err != nil {
		return
	}

	return b.String(), nil
}

// UnmarshalJSON implements the JSON interface for ListIDs.
func (e *ListIDs) UnmarshalJSON(in []byte) error {
	if e == nil {
		*e = ListIDs{}
	}

	for _, s := range splitComma(string(in)) {
		an, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return fmt.Errorf("%d is not a valid uint64", an)
		}
		*e = append(*e, an)
	}

	return nil
}

// MarshalJSON implements the JSON interface for ListIDs.
func (e ListIDs) MarshalJSON() (out []byte, err error) {
	var b bytes.Buffer
	_, err = b.WriteString("[")
	if err != nil {
		return
	}

	var arr []string
	for _, s := range e {
		arr = append(arr, strconv.FormatUint(s, 10))
	}
	_, err = b.WriteString(strings.Join(arr, ","))
	if err != nil {
		return
	}

	_, err = b.WriteString("]")
	if err != nil {
		return
	}

	return b.Bytes(), nil
}

// ListStrings is a list of String values.
type ListStrings []string

// Scan implements the Scanner interface for ListIDs.
func (e *ListStrings) Scan(value interface{}) (err error) {

	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	case []rune:
		str = string(v)
	default:
		return fmt.Errorf("array must be uint64, got: %T", value)
	}

	if e == nil {
		*e = ListStrings{}
	}

	str = trim(str, '{', '}')
	if len(str) == 0 {
		return nil
	}

	for _, s := range splitComma(string(str)) {
		*e = append(*e, s)
	}

	return nil
}

// Value implements the driver Valuer interface for ListStrings.
func (e ListStrings) Value() (v driver.Value, err error) {
	var b strings.Builder

	_, err = b.WriteString("{")
	if err != nil {
		return
	}

	var arr []string
	for _, s := range e {
		arr = append(arr, `"`+s+`"`)
	}
	_, err = b.WriteString(strings.Join(arr, ",") + "}")
	if err != nil {
		return
	}

	return b.String(), nil
}

// UnmarshalJSON implements the JSON interface for ListStrings.
func (e *ListStrings) UnmarshalJSON(in []byte) error {
	if e == nil {
		*e = ListStrings{}
	}
	s := string(in)

	s = trim(s, '[', ']')
	if len(s) == 0 {
		return nil
	}

	for _, p := range splitComma(s) {
		*e = append(*e, p)
	}

	return nil
}

// MarshalJSON implements the JSON interface for ListIDs.
func (e ListStrings) MarshalJSON() (out []byte, err error) {
	var b bytes.Buffer
	_, err = b.WriteString("[")
	if err != nil {
		return
	}

	var arr []string
	for _, s := range e {
		arr = append(arr, fmt.Sprintf("'%s'", strings.Replace("'", "''", s, -1)))
	}
	_, err = b.WriteString(strings.Join(arr, ","))
	if err != nil {
		return
	}

	_, err = b.WriteString("]")
	if err != nil {
		return
	}

	return b.Bytes(), nil
}

func splitComma(s string) []string {
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return c == ','
		}
	}
	lis := strings.FieldsFunc(s, f)

	var out []string
	for _, s := range lis {
		s = trim(s, '"', '"')
		out = append(out, s)
	}

	return out
}

func trim(s string, start, end rune) string {
	r0, size0 := utf8.DecodeRuneInString(s)
	if size0 == 0 {
		return s
	}
	if r0 != start {
		return s
	}

	r1, size1 := utf8.DecodeLastRuneInString(s)
	if size1 == 0 {
		return s
	}
	if r1 != end {
		return s
	}

	return s[size0 : len(s)-size1]
}
