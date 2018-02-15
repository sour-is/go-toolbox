package session

import (
	"github.com/patrickmn/go-cache"
	"net/http"
	"sour.is/x/ident"
	"time"
	"sour.is/x/uuid"
	"strings"
)

var store *cache.Cache

func init() {
	store = cache.New(4*time.Hour, 30*time.Second)

	ident.Register("session", CheckSession)
}

type User struct {
	Ident  string `json:"ident"`
	Aspect string `json:"aspect"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
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
		store.Set(id, u, cache.DefaultExpiration)

		return u
	}

	return User{}
}

func NewSession(ident, aspect, display string) (id string) {
	id = uuid.V4()

	u := User{
		ident,
		aspect,
		display,
		true,
	}

	store.Set(id, u, cache.DefaultExpiration)

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
	return false
}
func (u User) HasGroup(g ...string) bool {
	return false
}
func (u User) IsActive() bool {
	return u.Active
}
func (u User) GetDisplay() string {
	return u.Name
}
