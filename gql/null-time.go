package phoenix

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"io"
	"strings"
	"time"
)

// NullTime implements a nullable timestamp for sql
// swagger:type string
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface for NullTime.
func (n *NullTime) Scan(value interface{}) (err error) {
	n.Valid = true

	switch value.(type) {
	case time.Time:
		n.Time = value.(time.Time)
	default:
		n.Time, n.Valid = time.Time{}, false
	}

	return
}

// Value implements the driver Valuer interface.
func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

// UnmarshalGQL implements the graphql.Marshaler interface for NullTime
func (n *NullTime) UnmarshalGQL(v interface{}) (err error) {
	value, ok := v.(string)
	if !ok {
		return fmt.Errorf("points must be strings")
	}

	if value == "null" {
		n.Valid = false
		n.Time = time.Time{}
	} else {
		n.Valid = true
		n.Time, err = time.Parse(time.RFC3339, strings.Trim(value, ` "`))
	}

	return
}

// MarshalGQL implements the graphql.Marshaler interface for NullTime
func (n NullTime) MarshalGQL(w io.Writer) {
	var err error

	if n.Valid {
		_, err = w.Write([]byte(`"` + n.Time.Format(time.RFC3339+`"`)))
	} else {
		_, err = w.Write([]byte(`null`))
	}

	if err != nil {
		return
	}
}

// UnmarshalJSON implements the json.Marshaler interface for NullTime
func (n *NullTime) UnmarshalJSON(in []byte) (err error) {
	if bytes.Equal(in, []byte(`null`)) {
		n.Valid = false
		return
	}
	in = bytes.Trim(in, `"`)
	n.Time, err = time.Parse(time.RFC3339, string(in))
	if err != nil {
		n.Valid = false
		return
	}
	n.Valid = true

	return
}

// MarshalJSON implements the json.Marshaler interface for NullTime
func (n NullTime) MarshalJSON() (out []byte, err error) {
	if n.Valid {
		out = []byte(`"` + n.Time.Format(time.RFC3339) + `"`)
	} else {
		out = []byte("null")
	}

	return
}
