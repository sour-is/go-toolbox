package ident // import "sour.is/x/toolbox/ident"

/*
Include the desired drivers for ident in your main package.

import (
	_ "sour.is/x/toolbox/ident/header"
	_ "sour.is/x/toolbox/ident/mock"
	_ "sour.is/x/toolbox/ident/rubicon"
	_ "sour.is/x/toolbox/ident/session"
)

*/

import (
	"net/http"
	"strings"

	"sour.is/x/toolbox/log"
)

// Ident interface for a logged in user
type Ident interface {
	GetIdentity() string
	GetAspect() string

	HasRole(r ...string) bool
	HasGroup(g ...string) bool

	IsActive() bool
	GetDisplay() string
}

// Anonymous is a logged out user
var Anonymous = NewNullUser("anon", "none", "Guest User", false)

// IdentConfig key values to pass to an ident handler
type IdentConfig map[string]string

// IdentConfigs configs for each handler
type IdentConfigs map[string]IdentConfig

// IdentHandler handler function to read ident from HTTP request
type IdentHandler func(r *http.Request) Ident

// IdentSet set of handlers for ident
var IdentSet = make(map[string]IdentHandler)

// Config configs for handlers
var Config = IdentConfigs{}

// Register a ident handler
func Register(name string, fn IdentHandler) {
	name = strings.ToLower(name)
	IdentSet[name] = fn
}

// RegisterConfig for an ident handler
func RegisterConfig(name string, config map[string]string) {
	if _, ok := IdentSet[name]; !ok {
		log.Fatalf("IDENT: No handler registered for %s", name)
	}

	log.Infof("IDENT: Registered config for %s.", name)

	name = strings.ToLower(name)
	Config[name] = config
}

// GetIdent read ident from a list of ident handlers
func GetIdent(authList string, r *http.Request) Ident {
	for _, name := range strings.Fields(authList) {
		var i IdentHandler
		var ok bool

		if i, ok = IdentSet[name]; !ok {
			log.Errorf("GetIdentity Plugin [%s] does not exist!", name)
			panic("GetIdentity Plugin does not exist!")
		}

		u := i(r)

		if u.IsActive() {
			return u
		}
	}

	return Anonymous
}
