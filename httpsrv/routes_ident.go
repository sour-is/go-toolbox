package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"net/http"
	"strings"

	"sour.is/x/toolbox/ident"
)

// IdentRoute is a single route to handle
type IdentRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc HandlerFunc
}

// IdentRoutes is a list of routes to handle
type IdentRoutes []IdentRoute

// IdentRouteSet is a set of lists of routes to handle
var IdentRouteSet = make(map[string]IdentRoutes)

// IdentRegister adds an ident route to list
func IdentRegister(name string, routes IdentRoutes) {
	name = strings.ToLower(name)
	IdentRouteSet[name] = routes
}

// HandlerFunc is used by registered routes
type HandlerFunc func(ResponseWriter, *http.Request, ident.Ident)

// ServeHTTP handles a request
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *http.Request, i ident.Ident) {
	h := w.Header()
	h.Add("x-user-ident", i.GetIdentity())
	h.Add("x-user-aspect", i.GetAspect())
	h.Add("x-user-name", i.GetDisplay())

	f(w, r, i)
}
