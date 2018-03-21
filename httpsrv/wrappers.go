package httpsrv

import (
	"github.com/spf13/viper"
	"net/http"
	"sour.is/x/toolbox/ident"
)

func identWrapper(name string, hdlr HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var nw = WrapResponseWriter(w)

		var ok = true
		var id ident.Ident

		//log.Debug("EventPreAuth")
		if ok = runMiddleware(EventPreAuth, name, nw, r, ident.NullUser{}); ok {
			//log.Debug("EventAuth")
			id = ident.GetIdent(viper.GetString("http.idm"), r)
			//log.Debug("EventPreHandle", id)
			if ok = runMiddleware(EventPreHandle, name, nw, r, id); ok {
				hdlr.ServeHTTP(nw, r, id)
				//log.Debug("EventPreHandle")
				runMiddleware(EventPostHandle, name, nw, r, id)
			}
		}

		nw.StopTime()
		//log.Debug("EventComplete")
		runMiddleware(EventComplete, name, nw, r, id)
	})
}

func stdWrapper(hdlr http.HandlerFunc) HandlerFunc {
	return HandlerFunc(func(w ResponseWriter, r *http.Request, i ident.Ident) {
		hdlr.ServeHTTP(w, r)
	})
}

func assetWrapper(name, prefix string, hdlr http.FileSystem) http.Handler {
	fn := http.StripPrefix(prefix, http.FileServer(hdlr))
	return identWrapper(name, stdWrapper(fn.ServeHTTP))
}
