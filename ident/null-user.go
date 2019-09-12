package ident

import "net/http"

// NullUser implements a null ident
type NullUser struct {
	Ident  string `json:"ident"`
	Aspect string `json:"aspect"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

// NewNullUser creates a null user ident
func NewNullUser(ident, aspect, name string, active bool) NullUser {
	return NullUser{ident, aspect, name, active}
}

// GetIdentity returns identity
func (m NullUser) GetIdentity() string {
	return m.Ident
}

// GetAspect returns aspect
func (m NullUser) GetAspect() string {
	return m.Aspect
}

// HasRole returns true if matches role
func (m NullUser) HasRole(r ...string) bool {
	return m.Active
}

// HasGroup returns true if matches group
func (m NullUser) HasGroup(g ...string) bool {
	return m.Active
}

// GetGroups returns empty list
func (m NullUser) GetGroups() []string {
	return []string{}
}

// GetRoles returns empty list
func (m NullUser) GetRoles() []string {
	return []string{}
}

// GetMeta returns empty list
func (m NullUser) GetMeta() map[string]string {
	return make(map[string]string)
}

// IsActive returns true if active
func (m NullUser) IsActive() bool {
	return m.Active
}

// GetDisplay returns display name
func (m NullUser) GetDisplay() string {
	return m.Name
}

// MakeHandlerFunc returns handler func
func (m NullUser) MakeHandlerFunc() func(r *http.Request) Ident {
	return func(r *http.Request) Ident {
		return m
	}
}
