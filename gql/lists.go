package gql

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
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

	*e = ListIDs{}
	str = strings.Trim(str, `{}`)
	if len(str) == 0 {
		return nil
	}

	for _, s := range strings.Split(str, ",") {
		an, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return fmt.Errorf("%s is not a valid uint64", s)
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
	in = bytes.Trim(in, `[]`)
	str := string(in)
	if str == "" {
		return nil
	}
	for _, s := range strings.Split(str, ",") {
		s = strings.Trim(s, `"`)
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

	*e = ListStrings{}
	str = strings.Trim(str, `{}`)
	if len(str) == 0 {
		return nil
	}

	for _, s := range strings.Split(str, ",") {
		*e = append(*e, strings.Trim(s, `"`))
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
		arr = append(arr, s)
	}
	_, err = b.WriteString(strings.Join(arr, ",") + "}")
	if err != nil {
		return
	}

	return b.String(), nil
}

// UnmarshalJSON implements the JSON interface for ListStrings.
func (e *ListStrings) UnmarshalJSON(in []byte) error {
	in = bytes.Trim(in, `[]`)
	str := string(in)
	if str == "" {
		return nil
	}
	for _, s := range strings.Split(str, ",") {
		s = strings.Trim(s, `"`)
		*e = append(*e, s)
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
