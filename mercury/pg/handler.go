package pg

import (
	"sour.is/x/toolbox/dbm"
	"sour.is/x/toolbox/dbm/rsql"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
	"sour.is/x/toolbox/mercury"
	"sour.is/x/toolbox/mercury/dummy"
)

type postgresHandler struct {
	dummy.GroupsDummy
}

func init() {
	mercury.Register("*", 1, postgresHandler{})
}

func (postgresHandler) GetIndex(search mercury.NamespaceSearch, pgm *rsql.Program) (lis mercury.Config) {
	where := getWhere(search, dbm.GetDbInfo(Space{}))
	spaces, err := ListSpace(where, 0, 0, []string{"space asc"})
	if err != nil {
		log.Error(err)
		return nil
	}

	for _, s := range spaces {
		lis = append(lis, &mercury.Space{
			Space: s.Space,
			Tags:  s.Tags,
			Notes: s.Notes,
		})
	}

	return
}

func (p postgresHandler) GetObjects(search mercury.NamespaceSearch, pgm *rsql.Program, fields []string) mercury.Config {
	idx := p.GetIndex(search, pgm)
	spaceMap := make(map[string]int, len(idx))
	for u, s := range idx {
		spaceMap[s.Space] = u
	}

	where := getWhere(search, dbm.GetDbInfo(Config{}))
	values, err := ListConfig(where, 0, 0, []string{"space asc", "name asc"})
	if err != nil {
		log.Error(err)
		return nil
	}

	for _, v := range values {
		if u, ok := spaceMap[v.Space]; ok {
			idx[u].List = append(idx[u].List, mercury.Value{
				Space:  v.Space,
				Name:   v.Name,
				Seq:    v.Seq,
				Notes:  v.Notes,
				Tags:   v.Tags,
				Values: v.Values,
			})
		}
	}

	return idx
}

func (postgresHandler) GetRules(user ident.Ident) (rules mercury.Rules) {
	rules, err := GetRules(user)
	if err != nil {
		return nil
	}

	return
}

func (postgresHandler) GetNotify(event string) mercury.ListNotify {
	return GetNotify(event)
}

func (postgresHandler) WriteObjects(lis mercury.Config) error {
	err := dbm.Transaction(func(tx *dbm.Tx) error {
		return WriteConfig(tx, lis)
	})

	return err
}
