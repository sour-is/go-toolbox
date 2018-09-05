package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"net/http"
	"strings"
)

// HttpRoute is a route to handle
type HttpRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// HttpRoutes is a list of routes
type HttpRoutes []HttpRoute

// RouteSet is a set of routelists
var RouteSet = make(map[string]HttpRoutes)

// HttpRegister registers a routeset
func HttpRegister(name string, routes HttpRoutes) {
	name = strings.ToLower(name)
	RouteSet[name] = routes
}
