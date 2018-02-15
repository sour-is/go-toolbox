package ident // import "sour.is/x/ident"

import (
	"net/http"
	"sour.is/x/log"
	"strings"
)

type Ident interface {
	GetIdentity() string
	GetAspect() string
	HasRole(r ...string) bool
	HasGroup(g ...string) bool
	IsActive() bool
	GetDisplay() string
}

type IdentConfig map[string]string
type IdentConfigs map[string]IdentConfig
type IdentHandler func(r *http.Request) Ident

var IdentSet = make(map[string]IdentHandler)
var Config = IdentConfigs{"mock": {"identity": "mockuser", "aspect": "mock", "display_name": "MockUser"}}

func Register(name string, id IdentHandler) {
	name = strings.ToLower(name)
	IdentSet[name] = id
}

func RegisterConfig(name string, config map[string]string) {
	log.Infof("IDENT: Registered config for %s.", name)

	name = strings.ToLower(name)
	Config[name] = config
}

func GetIdent(authList string, r *http.Request) Ident {
	for _, name := range strings.Fields(authList) {
		var i IdentHandler
		var ok bool

		if i, ok = IdentSet[name]; !ok {
			log.Fatalf("GetIdentity Plugin [%s] does not exist!", name)
			panic("GetIdentity Plugin does not exist!")
		}

		u := i(r)

		if u.IsActive() {
			return u
		}
	}

	return NewNullUser("anon", "none", "Guest User", false)
}

type NullUser struct {
	Ident  string `json:"ident"`
	Aspect string `json:"aspect"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}
func NewNullUser(ident, aspect, name string, active bool) Ident {
	return NullUser{ident, aspect, name, active}
}
func (m NullUser) GetIdentity() string {
	return m.Ident
}
func (m NullUser) GetAspect() string {
	return m.Aspect
}
func (m NullUser) HasRole(r ...string) bool {
	return m.Active
}
func (m NullUser) HasGroup(g ...string) bool {
	return m.Active
}
func (m NullUser) IsActive() bool {
	return m.Active
}
func (m NullUser) GetDisplay() string {
	return m.Name
}
