package rubicon // import "sour.is/x/toolbox/ident/rubicon"

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"sour.is/x/toolbox/ident"
)

var store *cache.Cache

func init() {
	store = cache.New(5*time.Minute, 30*time.Second)

	ident.Register("rubicon", NewRubicon)
}

// User implements the ident.Ident interface
type User struct {
	ident    string
	name     string
	loggedIn bool
}

// IdmUser is the responce from rubicon api
type IdmUser struct {
	UserId    int    `json:"userId"`
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// NewRubicon checks requests and returns an ident.Ident
func NewRubicon(r *http.Request) ident.Ident {
	c := ident.Config["rubicon"]

	var id string

	if id = r.URL.Query().Get("access_token"); id == "" {
		return User{
			"anon",
			"Guest User",
			false,
		}
	}

	if user, ok := store.Get(id); ok == true {
		u := user.(*IdmUser)
		return User{
			u.UserName,
			u.FirstName + " " + u.LastName,
			true,
		}
	}

	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	netClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	req, _ := http.NewRequest("GET", c["idm"], nil)
	q := req.URL.Query()
	q.Add("user_token", id)
	req.URL.RawQuery = q.Encode()

	response, _ := netClient.Get(req.URL.String())

	buf, _ := ioutil.ReadAll(response.Body)

	var u = new(IdmUser)
	json.Unmarshal(buf, &u)

	if u.UserName == "" {
		return User{
			"anon",
			"Guest User",
			false,
		}
	}

	store.SetDefault(id, u)

	return User{
		u.UserName,
		u.FirstName + " " + u.LastName,
		true,
	}
}

// GetIdentity returns username
func (u User) GetIdentity() string {
	return u.ident
}

// GetAspect returns aspect
func (u User) GetAspect() string {
	return "rubicon"
}

// HasRole returns bool for tested roles
func (u User) HasRole(r ...string) bool {
	return true
}

// HasGroup returns bool for tested groups
func (u User) HasGroup(g ...string) bool {
	return true
}

// IsActive returns bool for logged in state
func (u User) IsActive() bool {
	return u.loggedIn
}

// GetDisplay returns human friendly name
func (u User) GetDisplay() string {
	return u.name
}
