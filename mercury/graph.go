package mercury

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"sour.is/x/toolbox/gql"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
)

// GraphMercury implements the resolvers for gqlgen
type GraphMercury struct{}

// Config returns a list of config items
func (GraphMercury) Config(ctx context.Context, search *string, query *gql.QueryInput) (lis []*Space, err error) {
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

	cfg := Registry.GetObjects(ns.String(), "", "")

	return cfg, nil
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

func fmtID(format string, args ...interface{}) string {
	s := fmt.Sprintf(format, args...)

	return base64.RawURLEncoding.EncodeToString([]byte(s))
}

func spaceFromID(id string) (string, error) {
	s, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return "", err
	}
	sp := strings.Split(string(s), ":")
	switch len(sp) {
	case 2:
		if sp[0] != "MercurySpace" {
			return "", fmt.Errorf("Invalid ID: %s", s)
		}

		return sp[1], nil

	default:
		return "", fmt.Errorf("Invalid ID: %s", s)
	}
}

func valueFromID(id string) (string, uint64, error) {
	s, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return "", 0, err
	}
	sp := strings.Split(string(s), ":")
	switch len(sp) {
	case 3:
		if sp[0] != "MercurySpace" {
			return "", 0, fmt.Errorf("Invalid ID: %s", s)
		}

		var seq uint64
		if seq, err = strconv.ParseUint(sp[2], 10, 64); err != nil {
			return "", 0, fmt.Errorf("invalid ID: %s", s)
		}

		return sp[1], seq, nil

	default:
		return "", 0, fmt.Errorf("Invalid ID: %s", s)
	}
}
