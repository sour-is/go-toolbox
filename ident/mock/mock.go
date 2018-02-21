package mock // sour.is/go/ident/mock

import "fmt"
import "net/http"
import "hash/crc32"
import "sour.is/go/ident"

func init() {
	ident.Register("mock", CheckMock)
}

type MockUser struct {
	ident  string
	aspect string
	name   string
	active bool
}

func NewMock(ident, aspect, name string, active bool) ident.Ident {
	return MockUser{ident, aspect, name, active}
}

func CheckMock(r *http.Request) ident.Ident {
	c := ident.Config["mock"]

	crc := crc32.ChecksumIEEE([]byte(r.RemoteAddr))

	return NewMock(
		c["identity"],
		c["aspect"],
		fmt.Sprintf("%s-%s-%x",
			c["display_name"],
			r.RemoteAddr,
			crc,
		),
		true)
}

func (m MockUser) GetIdentity() string {
	return m.ident
}
func (m MockUser) GetAspect() string {
	return m.aspect
}
func (m MockUser) HasRole(r ...string) bool {
	return m.active
}
func (m MockUser) HasGroup(g ...string) bool {
	return m.active
}
func (m MockUser) IsActive() bool {
	return m.active
}
func (m MockUser) GetDisplay() string {
	return m.name
}
