package session // import "sour.is/x/toolbox/ident/session"

import (
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/uuid"
	"github.com/spf13/viper"
)

var store *cache.Cache

var groupRoles map[string][]string
var userRoles map[string][]string
var userGroups map[string][]string

func init() {
	store = cache.New(4*time.Hour, 30*time.Second)

	ident.Register("session", CheckSession)
}

type User struct {
	Ident  string   `json:"ident"`
	Aspect string   `json:"aspect"`
	Name   string   `json:"name"`
	Active bool     `json:"active"`
	Groups map[string]struct{} `json:"groups"`
	Roles  map[string]struct{} `json:"roles"`
	Session string `json:"session"`
}

func Config() {
	if viper.IsSet("idm.session.user-roles") {
		userRoles = viper.GetStringMapStringSlice("idm.session.user-roles")
	}
	if viper.IsSet("idm.session.user-groups") {
		userGroups = viper.GetStringMapStringSlice("idm.session.user-groups")
	}
	if viper.IsSet("idm.session.group-roles") {
		groupRoles = viper.GetStringMapStringSlice("idm.session.group-roles")
	}
}

func GetSessionId(r *http.Request) string {
	var auth string

	if auth = r.Header.Get("authorization"); auth == "" {
		return ""
	}

	f := strings.Fields(auth)
	if len(f) < 2 || f[0] != "session" {
		return ""
	}

	return f[1]

}

func CheckSession(r *http.Request) ident.Ident {

	id := GetSessionId(r)

	if id == "" {
		return User{}
	}

	if user, ok := store.Get(id); ok == true {
		u := user.(User)
		store.SetDefault(id, u)

		return u
	}

	return User{}
}

func NewSession(ident, aspect, display string, groups []string, roles []string) (id string) {
	id = uuid.V4()

	u := User{
		Ident:  ident,
		Aspect: aspect,
		Name:   display,
		Active: true,
		Groups: make(map[string]struct{}),
		Roles: make(map[string]struct{}),
		Session: id,
	}

	if g, ok := userGroups[ident]; ok {
		groups = append(groups, g...)
	}

	if r, ok := userRoles[ident]; ok {
		roles = append(roles, r...)
	}

	for _, g := range groups {
		if r, ok := groupRoles[g]; ok {
			roles = append(roles, r...)
		}
	}

	for i := range groups {
		u.Groups[groups[i]] = struct{}{}
	}

	for i := range roles {
		u.Roles[roles[i]] = struct{}{}
	}

	store.SetDefault(id, u)

	return
}
func DeleteSession(id string) {
	store.Delete(id)
}

func (u User) GetIdentity() string {
	return u.Ident
}
func (u User) GetAspect() string {
	return u.Aspect
}
func (u User) HasRole(r ...string) bool {
	for _, k := range r {
		if _, ok := u.Roles[k]; ok {
			return true
		}
	}
	return false
}
func (u User) HasGroup(g ...string) bool {
	for _, k := range g {
		if _, ok := u.Groups[k]; ok {
			return true
		}
	}
	return false
}
func (u User) IsActive() bool {
	return u.Active
}
func (u User) GetDisplay() string {
	return u.Name
}
