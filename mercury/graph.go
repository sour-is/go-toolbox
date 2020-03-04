package mercury

import (
	"context"
	"fmt"
	"strings"

	"sour.is/x/toolbox/gql"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
)

// GraphMercury implements the resolvers for gqlgen
type GraphMercury struct{}

func doConfig(user ident.Ident, space string) ([]*Space, error) {
	rules := Registry.GetRules(user)

	ns := ParseNamespace(space)
	ns = rules.ReduceSearch(ns)

	cfg := Registry.GetObjects(ns.String(), "", "")
	return cfg, nil
}

// Config returns a list of config items
func (GraphMercury) Config(ctx context.Context, search *string, query *gql.QueryInput) (lis []*Space, err error) {
	user := ident.GetContextIdent(ctx)

	space := ""
	if search != nil {
		space = *search
	}
	if space == "" {
		space = "*"
	}

	return doConfig(user, space)
}

// WriteConfigText saves a config set formated in text
func (g GraphMercury) WriteConfigText(ctx context.Context, config string) (result string, err error) {
	r := strings.NewReader(config)
	c, err := parseText(r)
	if err != nil {
		return "ERR", err
	}
	var arr []*Space
	lis := c.ToArray()
	for i := range lis {
		arr = append(arr, lis[i])
	}
	return g.WriteConfig(ctx, arr)
}

// WriteConfig saves a space and attributes to database
func (GraphMercury) WriteConfig(ctx context.Context, config []*Space) (result string, err error) {
	user := ident.GetContextIdent(ctx)
	rules := Registry.GetRules(user)

	notify, err := Registry.GetNotify("updated")
	if err != nil {
		log.Error(err)
	}

	var notifyActive = make(map[string]struct{})
	var filteredConfigs Config
	for _, c := range config {
		log.Debug("CHECK ", c.Space, rules)
		if !rules.GetRoles("NS", c.Space).HasRole("write") {
			log.Debug("SKIP ", c.Space)
			continue
		}

		log.Debug("SAVE", c.Space)
		for _, n := range notify.Find(c.Space) {
			notifyActive[n.Name] = struct{}{}
		}
		filteredConfigs = append(filteredConfigs, c)
	}

	err = Registry.WriteObjects(filteredConfigs)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debug("SEND NOTIFYS ", notifyActive)

	for _, n := range notify {
		if _, ok := notifyActive[n.Name]; ok {
			err = n.sendNotify()
			if err != nil {
				log.Debug(err)
			}
		}
	}

	log.Debug("DONE!")

	return "OK", nil
}

// Value returns a joined value
func (GraphMercury) Value(ctx context.Context, value *Value) (string, error) {
	if value == nil {
		return "", nil
	}
	return strings.Join(value.Values, "\n"), nil
}

func NodeMercury(ctx context.Context, id []string) (gql.Node, error) {
	switch id[0] {
	case "MercurySpace":
		if len(id) != 2 {
			return nil, fmt.Errorf("ID missing space: %v", id)
		}

		user := ident.GetContextIdent(ctx)
		c, err := doConfig(user, id[1])
		if err != nil {
			return nil, err
		}

		if len(c) < 1 {
			return nil, fmt.Errorf("Not Found: %v", id)
		}

		return c[0], nil
	default:
		return nil, fmt.Errorf("Unsupported Node Type: %v", id[0])
	}
}
