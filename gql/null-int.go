package prospr

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/vektah/gqlgen/graphql"
	"sour.is/x/toolbox/log"
)

// MarshalNullInt32 is a nullable int marshaller
func MarshalNullInt32(b sql.NullInt32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if b.Valid {
			_, err := w.Write([]byte(fmt.Sprintf("%v", b.Int32)))
			if err != nil {
				log.Error(err)
			}
		} else {
			_, err := w.Write([]byte("null"))
			if err != nil {
				log.Error(err)
			}
		}
	})
}

// UnmarshalNullInt32 is a nullable int unmarshaller
func UnmarshalNullInt32(v interface{}) (i sql.NullInt32, err error) {
	switch v := v.(type) {
	case string:
		if v == "<nil>" || v == "null" || v == "" {
			return
		}
		i.Int32, err = strconv.ParseInt(v, 10, 32)
		if err == nil {
			i.Valid = true
		} else {
			log.Error(err)
		}

		return

	case int:
		i.Valid = true
		i.Int32 = int32(v)
		return

	case nil:
		return

	case json.Number:
		i.Int32, err = v.Int32()
		if err == nil {
			i.Valid = true
		} else {
			log.Error(err)
		}

		return

	default:
		err = fmt.Errorf("%T is not a int32 or null", v)
		log.Error(err)

		return
	}
}

// MarshalNullInt64 is a nullable int marshaller
func MarshalNullInt64(b sql.NullInt64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if b.Valid {
			_, err := w.Write([]byte(fmt.Sprintf("%v", b.Int64)))
			if err != nil {
				log.Error(err)
			}
		} else {
			_, err := w.Write([]byte("null"))
			if err != nil {
				log.Error(err)
			}
		}
	})
}

// UnmarshalNullInt64 is a nullable int unmarshaller
func UnmarshalNullInt64(v interface{}) (i sql.NullInt64, err error) {
	switch v := v.(type) {
	case string:
		if v == "<nil>" || v == "null" || v == "" {
			return
		}
		i.Int64, err = strconv.ParseInt(v, 10, 64)
		if err == nil {
			i.Valid = true
		} else {
			log.Error(err)
		}

		return

	case int:
		i.Valid = true
		i.Int64 = int64(v)
		return

	case nil:
		return

	case json.Number:
		i.Int64, err = v.Int64()
		if err == nil {
			i.Valid = true
		} else {
			log.Error(err)
		}

		return

	default:
		err = fmt.Errorf("%T is not a int64 or null", v)
		log.Error(err)

		return
	}
}

// NullUint64 represents an int64 that may be null.
// NullUint64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullUint64 struct {
	Uint64 uint64
	Valid  bool // Valid is true if Uint64 is not NULL
}

// Scan implements the Scanner interface.
func (n *NullUint64) Scan(value interface{}) (err error) {
	if value == nil {
		n.Uint64, n.Valid = 0, false
		return nil
	}
	n.Valid = true

	switch t := value.(type) {
	case string:
		if t == "" {
			n.Valid = false
			return
		}

		n.Uint64, err = strconv.ParseUint(t, 10, 64)
	case int:
		n.Uint64, err = uint64(t), nil
	case int64:
		n.Uint64, err = uint64(t), nil
	case json.Number:
		var i int64
		i, err = t.Int64()
		n.Uint64 = uint64(i)
	case float64:
		n.Uint64, err = uint64(t), nil
	}

	return
}

// Value implements the driver Valuer interface.
func (n NullUint64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return fmt.Sprintf("%v", n.Uint64), nil
}

// MarshalNullUint overrides the default data type of NullUint for graphQL
func MarshalNullUint64(t uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, err := io.WriteString(w, strconv.FormatUint(t, 10))
		if err != nil {
			return
		}
	})
}

// UnmarshalNullUint overrides the default data type of NullUint for graphQL
func UnmarshalNullUint64(v interface{}) (n NullUint64, err error) {
	if v == nil {
		n.Uint64, n.Valid = 0, false
		return
	}
	n.Valid = true

	switch t := v.(type) {
	case string:
		if t == "" {
			n.Valid = false
			return
		}
		n.Uint64, err = strconv.ParseUint(t, 10, 64)
	case int:
		n.Uint64, err = uint64(t), nil
	case int64:
		n.Uint64, err = uint64(t), nil
	case json.Number:
		var i int64
		i, err = t.Int64()
		n.Uint64 = uint64(i)
	case float64:
		n.Uint64, err = uint64(t), nil
	}

	return
}
