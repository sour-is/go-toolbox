package mock // import "sour.is/x/toolbox/ident/mock"

import "fmt"
import "net/http"
import "hash/crc32"
import "sour.is/x/toolbox/ident"

func init() {
	ident.Register("mock", CheckMock)
}

// MockUser implements ident.Ident
type MockUser struct {
	ident  string
	aspect string
	name   string
	active bool
}

// NewMock returns a new MockUser
func NewMock(ident, aspect, name string, active bool) ident.Ident {
	return MockUser{ident, aspect, name, active}
}

// CheckMock checks the http request for mock username
func CheckMock(r *http.Request) ident.Ident {
	c := ident.Config["mock"]

	crc := crc32.ChecksumIEEE([]byte(r.RemoteAddr))

	return NewMock(
		c["ident"],
		c["aspect"],
		fmt.Sprintf("%s-%s-%x",
			c["name"],
			r.RemoteAddr,
			crc,
		),
		true)
}

// GetIdentity returns username
func (m MockUser) GetIdentity() string {
	return m.ident
}

// GetAspect returns aspect
func (m MockUser) GetAspect() string {
	return m.aspect
}

// HasRole returns bool for tested roles
func (m MockUser) HasRole(r ...string) bool {
	return m.active
}

// HasGroup returns bool for tested groups
func (m MockUser) HasGroup(g ...string) bool {
	return m.active
}

// IsActive returns bool for logged in state
func (m MockUser) IsActive() bool {
	return m.active
}

// GetDisplay returns human friendly name
func (m MockUser) GetDisplay() string {
	return m.name
}
