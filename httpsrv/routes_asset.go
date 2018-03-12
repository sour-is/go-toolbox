package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"net/http"
	"strings"
)

type AssetRoute struct {
	Name        string
	Path        string
	HandlerFunc http.FileSystem
}
type AssetRoutes []AssetRoute

var AssetSet = make(map[string]AssetRoutes)

func AssetRegister(name string, routes AssetRoutes) {
	name = strings.ToLower(name)
	AssetSet[name] = routes
}
