package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"net/http"
	"strings"
	"os"
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


// This function will return the app for subdirectories
// of the application so it can remove the /#/ urls.
func FsHtml5(fn http.FileSystem) http.FileSystem {
	return fsWrap{F: fn}
}

type fsWrap struct {
	F http.FileSystem
}

func (fs fsWrap) Open(name string) (b http.File, err error) {
	b, err = fs.F.Open(name)

	if os.IsNotExist(err) {
		if strings.HasSuffix(name, "/app.js") {
			return fs.F.Open("/app.js")
		}
		if strings.HasSuffix(name, "/index.html") {
			return fs.F.Open("/index.html")
		}
		for _, n := range []string{".ico", ".css", ".js", ".png", ".jpg", ".svg", ".html", ".txt", ".json"} {
			if strings.HasSuffix(name, n) {
				return
			}
		}

		return fs.F.Open("/")
	}
	return
}
