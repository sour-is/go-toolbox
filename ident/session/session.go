package session // import "sour.is/x/toolbox/ident/session"

import (
	"net/http"
		"time"

	"github.com/patrickmn/go-cache"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/uuid"
	"github.com/spf13/viper"
	"sour.is/x/toolbox/log"
	"strings"
)

var store *cache.Cache

var cookieName string
var groupRoles map[string][]string
var userRoles map[string][]string
var userGroups map[string][]string

var sessionExpire = 10 * time.Minute
var cookieExpire =  24 * time.Hour

func init() {
	store = cache.New(cookieExpire, 30*time.Second)

	ident.Register("session", CheckSession)
}

// User is an ident.Ident
type User struct {
	Ident  string   `json:"ident"`
	Aspect string   `json:"aspect"`
	Name   string   `json:"name"`
	Active bool     `json:"active"`
	Groups map[string]struct{} `json:"groups"`
	Roles  map[string]struct{} `json:"roles"`
	Session string `json:"session"`
	Cookie  string `json:"cookie"`
}

// Config sets up the session module
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
	if viper.IsSet("idm.session.cookie") {
		cookieName = viper.GetString("idm.session.cookie")
	}
	if viper.IsSet("idm.session.cookie-ttl") {
		cookieExpire = time.Duration(viper.GetInt64("idm.session.cookie-ttl")) * time.Minute
	}
	if viper.IsSet("idm.session.session-ttl") {
		sessionExpire = time.Duration(viper.GetInt64("idm.session.cookie-ttl")) * time.Minute
	}
}

// GetSessionId attempts to read a session id out of request
func getSessionID(r *http.Request) string {

	// Try reading from cookies
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Debug(err)
	}

	if cookie != nil {
		return cookie.Value
	}

	// Try reading from Authorization
	if auth := r.Header.Get("authorization"); auth != "" {
		f := strings.Fields(auth)
		if len(f) < 2 || f[0] != "session" {
			return ""
		}

		return f[1]
	}

	return ""
}

// CheckSession is called by ident to lookup the user
func CheckSession(r *http.Request) ident.Ident {
	id := getSessionID(r)
	return GetSessionID(id)
}

func GetSessionID(id string) ident.Ident {
	if id == "" {
		return User{}
	}

	if user, ok := store.Get(id); ok == true {
		u := user.(User)
		store.Set(u.Session, u, sessionExpire)
		store.Set(u.Session, u, cookieExpire)

		return u
	}

	return User{}
}

// NewSession creates a new session and returns an ident.Ident
func NewSession(ident, aspect, display string, groups []string, roles []string) (ident.Ident) {
	u := User{
		Ident:  ident,
		Aspect: aspect,
		Name:   display,
		Active: true,
		Groups: make(map[string]struct{}),
		Roles: make(map[string]struct{}),
		Session: "S" + uuid.V4(),
		Cookie:  "C" + uuid.V4(),
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

	store.Set(u.Session, u, sessionExpire)
	store.Set(u.Cookie, u, cookieExpire)

	return u
}
// DeleteSession removes the session
func DeleteSession(id string) {
	u, ok := store.Get(id)
	if !ok {
		return
	}

	store.Delete(u.(User).Session)
	store.Delete(u.(User).Cookie)
}

// GetIdentity returns the identity of user
func (u User) GetIdentity() string {
	return u.Ident
}
// GetAspect returns the current aspect of user
func (u User) GetAspect() string {
	return u.Aspect
}
// HasRole returns true if any roles match
func (u User) HasRole(r ...string) bool {
	for _, k := range r {
		if _, ok := u.Roles[k]; ok {
			return true
		}
	}
	return false
}
// HasGroup returns true if any groups match
func (u User) HasGroup(g ...string) bool {
	for _, k := range g {
		if _, ok := u.Groups[k]; ok {
			return true
		}
	}
	return false
}
// IsActive returns true if user is active
func (u User) IsActive() bool {
	return u.Active
}
// GetDisplay returns the display name of user
func (u User) GetDisplay() string {
	return u.Name
}

func (u User) GetCookie() *http.Cookie {
	return &http.Cookie {
		Name:     cookieName,
		Value:    u.Cookie,
		HttpOnly: true,
		Secure:   viper.GetBool("idm.session.secure"),
		Path:     "/",
		MaxAge:   84600,
	}
}