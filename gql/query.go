package gql

// QueryInput allows you to filter and page the search.
type QueryInput struct {
	Search *string  `json:"search"`
	Limit  *uint64  `json:"limit"`
	Offset *uint64  `json:"offset"`
	Sort   []string `json:"sort"`
}
