package mock // import "sour.is/x/toolbox/ident/mock"

import "fmt"
import "net/http"
import "hash/crc32"
import "sour.is/x/toolbox/ident"

func init() {
	ident.Register("mock", CheckMock)
}

// User implements ident.Ident
type User struct {
	ident  string
	aspect string
	name   string
	groups map[string]struct{}
	roles  map[string]struct{}
	meta   map[string]string
	active bool
}

// NewMock returns a new MockUser
func NewMock(ident, aspect, name string, groups []string, roles []string, meta map[string]string, active bool) ident.Ident {
	r := make(map[string]struct{}, len(roles))
	for _, name := range roles {
		r[name] = struct{}{}
	}
	g := make(map[string]struct{}, len(groups))
	for _, name := range groups {
		g[name] = struct{}{}
	}

	return User{ident, aspect, name, g, r, meta, active}
}

// CheckMock checks the http request for mock username
func CheckMock(r *http.Request) ident.Ident {
	c := ident.GetConfigs()["mock"]

	crc := crc32.ChecksumIEEE([]byte(r.RemoteAddr))

	return NewMock(
		c["ident"],
		c["aspect"],
		fmt.Sprintf("%s-%s-%x",
			c["name"],
			r.RemoteAddr,
			crc,
		),
		nil,
		nil,
		nil,
		true)
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
func (m User) HasRole(lis ...string) bool {
	if !m.active {
		return false
	}

	for _, name := range lis {
		if _, ok := m.roles[name]; ok {
			return true
		}
	}

	return false
}

// HasGroup returns bool for tested groups
func (m User) HasGroup(lis ...string) bool {
	if !m.active {
		return false
	}

	for _, name := range lis {
		if _, ok := m.groups[name]; ok {
			return true
		}
	}

	return false
}

// IsActive returns bool for logged in state
func (m User) IsActive() bool {
	return m.active
}

// GetDisplay returns human friendly name
func (m User) GetDisplay() string {
	return m.name
}

// GetGroups returns empty list
func (m User) GetGroups() []string {
	lis := make([]string, 0, len(m.roles))
	for i := range m.roles {
		lis = append(lis, i)
	}
	return lis
}

// GetRoles returns empty list
func (m User) GetRoles() []string {
	lis := make([]string, 0, len(m.groups))
	for i := range m.roles {
		lis = append(lis, i)
	}
	return lis
}

// GetMeta returns empty list
func (m User) GetMeta() map[string]string {
	return m.meta
}
