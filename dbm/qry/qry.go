package qry

import (
	"fmt"
	"strings"

	"sour.is/x/toolbox/dbm"
	"sour.is/x/toolbox/dbm/rsql/squirrel"
	"sour.is/x/toolbox/gql"
	"sour.is/x/toolbox/log"
)

// Input is a query for a given table
type Input struct {
	dbm.DbInfo
	Search interface{}
	Limit  uint64
	Offset uint64
	Sort   []string
}

// Qry builds a query for a table based on input
func Qry(db dbm.DbInfo, input *gql.QueryInput) (out Input, err error) {

	out.DbInfo = db
	out.Search = nil
	out.Limit = 0
	out.Offset = 0

	if input == nil {
		return
	}

	if input.Search != nil {
		out.Search, err = squirrel.Query(*input.Search, db)
	}

	if input.Offset != nil {
		out.Offset = *input.Offset
	}

	if input.Limit != nil {
		out.Limit = *input.Limit
	}

	for _, sort := range input.Sort {
		s := strings.Fields(sort)
		log.Debug(s)
		ord := "asc"
		if len(s) < 1 {
			continue
		}
		if len(s) > 1 {
			s[1] = strings.ToLower(s[1])
			if s[1] == "dec" || s[1] == "desc" {
				ord = "desc"
			}
		}
		log.Debug(db)
		if i, ok := db.Index(s[0]); ok {
			out.Sort = append(out.Sort, fmt.Sprintf("%v %v", db.Cols[i], ord))
		}
	}

	return
}
