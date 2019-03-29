package app

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"sour.is/x/toolbox/dbm/rsql"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
	"sour.is/x/toolbox/mercury"
	"sour.is/x/toolbox/mercury/dummy"
)

const (
	appDotSettings = "app.settings"
	appDotPriority = "app.priority"
)

type appDefault struct {
	dummy.IndexDummy
	dummy.ObjectsDummy
	dummy.WriteDummy
	dummy.NotifyDummy
}

func init() {
	mercury.Register("*", 0, appDefault{})
	mercury.Register("app.*", math.MaxInt64, appConfig{})
}

// Rules returns nil
func (appDefault) GetRules(u ident.Ident) (lis mercury.Rules) {

	if u.HasRole("admin") {
		lis = append(lis,
			mercury.Rule{
				Role:  "admin",
				Type:  "NS",
				Match: "*",
			},
			mercury.Rule{
				Role:  "write",
				Type:  "NS",
				Match: "*",
			},
			mercury.Rule{
				Role:  "admin",
				Type:  "GR",
				Match: "*",
			},
		)
	} else if u.HasRole("write") {
		lis = append(lis,
			mercury.Rule{
				Role:  "write",
				Type:  "NS",
				Match: "*",
			},
		)
	} else if u.HasRole("read") {
		lis = append(lis,
			mercury.Rule{
				Role:  "read",
				Type:  "NS",
				Match: "*",
			},
		)
	}

	return lis
}

type appConfig struct {
	dummy.WriteDummy
	dummy.NotifyDummy
}

// Index returns nil
func (appConfig) GetIndex(search mercury.NamespaceSearch, _ *rsql.Program) (lis mercury.ArraySpace) {

	if search.Match(appDotSettings) {
		lis = append(lis, mercury.Space{Space: appDotSettings})
	}

	if search.Match(appDotPriority) {
		lis = append(lis, mercury.Space{Space: appDotPriority})
	}

	return
}

// Objects returns nil
func (appConfig) GetObjects(search mercury.NamespaceSearch, _ *rsql.Program, _ []string) (lis mercury.ArraySpace) {

	if search.Match(appDotSettings) {
		space := mercury.Space{
			Space: appDotSettings,
		}

		keys := viper.AllKeys()
		sort.Strings(keys)

		for i, key := range keys {
			var val []string

			s := viper.GetString(key)

			if s != "" {
				val = strings.Split(s, "\n")
				log.Debug("split ", val)
			} else if viper.IsSet(key) {
				val = viper.GetStringSlice(key)
				log.Debug("slice ", val)
			} else {
				v := viper.Get(key)
				val = strings.Split(fmt.Sprintf("%#v", v), "\n")
			}

			space.List = append(space.List, mercury.Value{
				Seq:    uint64(i),
				Name:   key,
				Values: val,
			})
		}

		lis = append(lis, space)
	}

	if search.Match(appDotPriority) {
		space := mercury.Space{
			Space: appDotPriority,
		}

		for i, key := range mercury.Registry {
			space.List = append(space.List, mercury.Value{
				Seq:    uint64(i),
				Name:   key.Match,
				Values: []string{fmt.Sprint(key.Priority)},
			})
		}

		lis = append(lis, space)
	}
	return
}

// Rules returns nil
func (appConfig) GetRules(u ident.Ident) (lis mercury.Rules) {

	if u.HasRole("admin") {
		lis = append(lis,
			mercury.Rule{
				Role:  "read",
				Type:  "NS",
				Match: appDotSettings,
			},
			mercury.Rule{
				Role:  "read",
				Type:  "NS",
				Match: appDotPriority,
			},
		)
	}

	return lis
}
