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

func (a AssetRoutes) Len() int      { return len(a) }
func (a AssetRoutes) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a AssetRoutes) Less(i, j int) bool {
	if len(a[i].Path) == len(a[j].Path) {
		return a[i].Path >= a[j].Path
	}

	return len(a[i].Path) > len(a[j].Path)
}
