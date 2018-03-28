package httpsrv

import (
	"github.com/spf13/viper"
	"net/http"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
)

func identWrapper(name string, hdlr HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var nw = WrapResponseWriter(w)

		var ok = true
		var id ident.Ident = ident.NullUser{}

		log.NilDebug("EventPreAuth")
		if ok = runMiddleware(EventPreAuth, name, nw, r, id); ok {
			log.NilDebug("EventAuth")
			id = ident.GetIdent(viper.GetString("http.idm"), r)
			log.NilDebug("EventPreHandle", id)
			if ok = runMiddleware(EventPreHandle, name, nw, r, id); ok {
				hdlr.ServeHTTP(nw, r, id)
				log.NilDebug("EventPreHandle")
				runMiddleware(EventPostHandle, name, nw, r, id)
			}
		}

		nw.StopTime()
		log.NilDebug("EventComplete")
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
