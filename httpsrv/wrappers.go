package httpsrv

import (
	"github.com/spf13/viper"
	"net/http"
	"sour.is/x/toolbox/ident"
)

func identWrapper(name string, hdlr HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var nw = wrapResponseWriter(w)

		if ok := runMiddleware(EventPreAuth, name, nw, r, ident.NullUser{}); !ok {
			return
		}

		id := ident.GetIdent(viper.GetString("http.idm"), r)

		runMiddleware(EventPreHandle, name, nw, r, id)

		hdlr.ServeHTTP(nw, r, id)

		runMiddleware(EventPostHandle, name, nw, r, id)

		nw.StopTime()

		runMiddleware(EventComplete, name, nw, r, id)
	})
}

func stdWrapper(hdlr http.HandlerFunc) HandlerFunc {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request, i ident.Ident) {
		hdlr.ServeHTTP(w, r)
	})
}

func assetWrapper(name, prefix string, hdlr http.FileSystem) http.Handler {
	fn := http.StripPrefix(prefix, http.FileServer(hdlr))
	return identWrapper(name, stdWrapper(fn.ServeHTTP))
}
