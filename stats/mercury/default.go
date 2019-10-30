package stats

import (
	"math"

	"sour.is/x/toolbox/dbm/rsql"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/mercury"
	"sour.is/x/toolbox/mercury/dummy"
	"sour.is/x/toolbox/stats"
)

func init() {
	mercury.Register("stats.*", math.MaxInt64, config{})
}

type config struct {
	dummy.WriteDummy
	dummy.NotifyDummy
}

func (config) GetIndex(search mercury.NamespaceSearch, _ *rsql.Program) (lis mercury.Config) {
	registry := stats.GetRegistry()
	for k := range registry {
		name := "stats." + k
		if search.Match(name) {
			space := mercury.Space{}
			space.Space = name
			lis = append(lis, &space)
		}
	}

	return
}

func (config) GetObjects(search mercury.NamespaceSearch, _ *rsql.Program, _ []string) (lis mercury.Config) {
	registry := stats.GetRegistry()

	for k, fn := range registry {
		name := "stats." + k
		if search.Match(name) {
			exp := fn()
			space := mercury.Space{}
			space.Space = name
			space.List = exp.Exposition().ToSpaceValues()
			lis = append(lis, &space)
		}
	}

	return
}

// Rules returns nil
func (config) GetRules(u ident.Ident) (lis mercury.Rules) {
	if u.HasRole("admin") {
		registry := stats.GetRegistry()

		for k := range registry {
			lis = append(lis, mercury.Rule{
				Role:  "read",
				Type:  "NS",
				Match: "stats." + k,
			})
		}
	}

	return lis
}
