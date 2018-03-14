package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"net/http"
	"os"
	"strings"
	"sour.is/x/toolbox/log"
)

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
