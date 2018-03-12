package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"net/http"
	"strings"

	"sour.is/x/toolbox/ident"
)

type IdentRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc HandlerFunc
}
type IdentRoutes []IdentRoute

var IdentRouteSet = make(map[string]IdentRoutes)

func IdentRegister(name string, routes IdentRoutes) {
	name = strings.ToLower(name)
	IdentRouteSet[name] = routes
}

type HandlerFunc func(http.ResponseWriter, *http.Request, ident.Ident)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, i ident.Ident) {
	h := w.Header()
	h.Add("x-user-ident", i.GetIdentity())
	h.Add("x-user-aspect", i.GetAspect())
	h.Add("x-user-name", i.GetDisplay())

	f(w, r, i)
}
