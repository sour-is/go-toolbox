package header // import "sour.is/x/toolbox/ident/header"

import (
	"net/http"

	"sour.is/x/toolbox/ident"
)

func init() {
	ident.Register("header", NewUser)
}

// User implements ident.Ident
type User struct {
	ident  string
	aspect string
	name   string
	active bool
}

// NewUser reads user info from http request
func NewUser(r *http.Request) ident.Ident {

	config := ident.Config["header"]

	hIdent := config["ident"]
	if hIdent == "" {
		hIdent = "user_ident"
	}
	hAspect := config["aspect"]
	if hAspect == "" {
		hAspect = "user_aspect"
	}
	hName := config["name"]
	if hName == "" {
		hName = "user_name"
	}

	loggedIn := true

	ident := r.Header.Get(hIdent)
	if ident == "" {
		ident = "anon"
	}

	aspect := r.Header.Get(hAspect)
	if aspect == "" {
		aspect = "default"
	}

	name := r.Header.Get(hName)
	if name == "" {
		name = ident
	}

	if ident == "anon" {
		name = "Guest User"
		loggedIn = false
	}

	return User{
		ident,
		aspect,
		name,
		loggedIn,
	}
}

// GetIdentity returns username
func (m User) GetIdentity() string {
	return m.ident
}

// GetAspect returns aspect
func (m User) GetAspect() string {
	return m.aspect
}

// HasRole returns bool for tested roles
func (m User) HasRole(r ...string) bool {
	return m.active
}

// HasGroup returns bool for tested groups
func (m User) HasGroup(g ...string) bool {
	return m.active
}

// IsActive returns bool for logged in state
func (m User) IsActive() bool {
	return m.active
}

// GetDisplay returns human friendly name
func (m User) GetDisplay() string {
	return m.name
}
