package mercury

import (
	"context"
	"strings"

	"sour.is/x/toolbox/gql"
	"sour.is/x/toolbox/ident"
)

// GraphMercury implements the resolvers for gqlgen
type GraphMercury struct{}

// Config returns a list of config items
func (GraphMercury) Config(ctx context.Context, search *string, query *gql.QueryInput) ([]Space, error) {
	user := ident.GetContextIdent(ctx)

	rules := Registry.GetRules(user)
	space := ""
	if search != nil {
		space = *search
	}
	if space == "" {
		space = "*"
	}

	ns := ParseNamespace(space)
	ns = rules.ReduceSearch(ns)

	lis := Registry.GetObjects(ns.String(), "", "")

	return lis, nil
}

// RegistryStore saves a space and attributes to database
func (GraphMercury) RegistryStore(ctx context.Context, space string, attributes []Value) (o Space, err error) {
	config := make(SpaceMap)
	if c, ok := config[space]; !ok {
		config[space] = Space{Space: space, List: attributes}
	} else {
		c.List = attributes
		config[space] = c
	}

	var lis ArraySpace
	for _, c := range config {
		lis = append(lis, c)
	}

	err = Registry.WriteObjects(lis)
	if err != nil {
		return
	}

	return
}

// Value returns a joined value
func (GraphMercury) Value(ctx context.Context, value *Value) (string, error) {
	if value == nil {
		return "", nil
	}
	return strings.Join(value.Values, "\n"), nil
}
