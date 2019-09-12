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
	GetDisplay() string

	GetGroups() []string
	GetRoles() []string
	GetMeta() map[string]string

	HasRole(r ...string) bool
	HasGroup(g ...string) bool

	IsActive() bool
}

// Anonymous is a logged out user
var Anonymous = NewNullUser("anon", "none", "Guest User", false)

// Config key values to pass to an ident handler
type Config map[string]string

// Configs configs for each handler
type Configs map[string]Config

// Handler handler function to read ident from HTTP request
type Handler func(r *http.Request) Ident

var handlers = make(map[string]Handler)
var configs = Configs{}

// Register a ident handler
func Register(name string, fn Handler) {
	name = strings.ToLower(name)
	handlers[name] = fn
}

// GetConfigs that have been registered
func GetConfigs() Configs {
	return configs
}

// GetConfig that matches name
func GetConfig(name string) Config {
	return configs[name]
}

// GetHandlers get handlers that are registered
func GetHandlers() map[string]Handler {
	return handlers
}

// RegisterConfig for an ident handler
func RegisterConfig(name string, config map[string]string) {
	if _, ok := handlers[name]; !ok {
		log.Fatals("IDENT: No handler registered", "name", name)
	}

	log.Infos("IDENT: Registered config", "name", name)

	name = strings.ToLower(name)
	configs[name] = config
}

// GetIdent read ident from a list of ident handlers
func GetIdent(authList string, r *http.Request) Ident {
	for _, name := range strings.Fields(authList) {
		var i Handler
		var ok bool

		if i, ok = handlers[name]; !ok {
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
