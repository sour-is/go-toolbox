package mercury

import (
	"context"
	"strings"

	"sour.is/x/toolbox/gql"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
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

// WriteConfigText saves a config set formated in text
func (g GraphMercury) WriteConfigText(ctx context.Context, config string) (string, error) {
	r := strings.NewReader(config)
	c, err := parseText(r)
	if err != nil {
		return "ERR", err
	}
	return g.WriteConfig(ctx, c.ToArray())
}

// WriteConfig saves a space and attributes to database
func (GraphMercury) WriteConfig(ctx context.Context, config []Space) (result string, err error) {
	user := ident.GetContextIdent(ctx)
	rules := Registry.GetRules(user)

	notify, err := Registry.GetNotify("updated")
	if err != nil {
		log.Error(err)
	}

	var notifyActive = make(map[string]struct{})
	filteredConfigs := make(ArraySpace, 0)
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
