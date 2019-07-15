package app

import (
	"fmt"
	"math"
	"os"
	"os/user"
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
	appDotHost     = "app.host"
)

func init() {
	mercury.Register("app.*", math.MaxInt64, appConfig{})
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

	if search.Match(appDotHost) {
		lis = append(lis, mercury.Space{Space: appDotHost})
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
				log.NilDebug("split ", val)
			} else if viper.IsSet(key) {
				val = viper.GetStringSlice(key)
				log.NilDebug("slice ", val)
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

	if search.Match(appDotHost) {
		if usr, err := user.Current(); err == nil {
			space := mercury.Space{
				Space: appDotHost,
			}

			hostname, _ := os.Hostname()
			wd, _ := os.Getwd()
			grp, _ := usr.GroupIds()
			space.List = []mercury.Value{
				{
					Seq:    1,
					Name:   "hostname",
					Values: []string{hostname},
				},
				{
					Seq:    2,
					Name:   "username",
					Values: []string{usr.Username},
				},
				{
					Seq:    3,
					Name:   "uid",
					Values: []string{usr.Uid},
				},
				{
					Seq:    4,
					Name:   "gid",
					Values: []string{usr.Gid},
				},
				{
					Seq:    5,
					Name:   "display",
					Values: []string{usr.Name},
				},
				{
					Seq:    6,
					Name:   "home",
					Values: []string{usr.HomeDir},
				},
				{
					Seq:    7,
					Name:   "groups",
					Values: grp,
				},
				{
					Seq:    8,
					Name:   "pid",
					Values: []string{fmt.Sprintf("%v", os.Getpid())},
				},
				{
					Seq:    9,
					Name:   "wd",
					Values: []string{wd},
				},
				{
					Seq:    10,
					Name:   "environ",
					Values: os.Environ(),
				},
			}

			lis = append(lis, space)
		}
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
			mercury.Rule{
				Role:  "read",
				Type:  "NS",
				Match: appDotHost,
			},
		)
	}

	return lis
}
