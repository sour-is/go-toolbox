package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"net/http"
	"strings"
)

type HttpRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type HttpRoutes []HttpRoute

var RouteSet = make(map[string]HttpRoutes)

func HttpRegister(name string, routes HttpRoutes) {
	name = strings.ToLower(name)
	RouteSet[name] = routes
}
