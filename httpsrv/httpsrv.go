package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"sort"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
	"sync"
)

// Example Usage

// package main
//
// import (
// 	"log"
// 	"sour.is/x/httpsrv"
//   _ "sour.is/x/httpsrv/routes"
// )
//
// func main() {
//     log.Println("Listen and Serve on", ":8080")
//     httpsrv.Run(config)
// }

// ModuleHandler holds registered handlers for httpsrv
type ModuleHandler func(map[string]string)

var modules = make(map[string]ModuleHandler)

// SignalStartup channel is closed when httpsrv starts up
var SignalStartup = make(chan struct{})
// SignalShutdown channel is closed when httpsrv shuts down
var SignalShutdown = make(chan struct{})
// WaitShutdown registers services to wait for graceful shutdown
var WaitShutdown sync.WaitGroup
// NewRouter
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for setName, routes := range RouteSet {
		for _, route := range routes {
			name := setName + "::" + route.Name

			handler := identWrapper(name, stdWrapper(route.HandlerFunc))

			log.Infof("Registered HTTP: %s for %s %s",
				name, route.Method, route.Pattern)

			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(name).
				Handler(handler)
		}
	}

	for setName, routes := range IdentRouteSet {
		for _, route := range routes {
			name := setName + "::" + route.Name

			handler := identWrapper(name, route.HandlerFunc)

			log.Infof("Registered IDENT: %s for %s %s",
				name, route.Method, route.Pattern)

			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(name).
				Handler(handler)
		}
	}

	var assets AssetRoutes
	for setName, assetRoutes := range AssetSet {
		for _, route := range assetRoutes {
			route.Name = setName + "::" + route.Name
			assets = append(assets, route)
		}
	}
	sort.Sort(assets)

	for _, route := range assets {
		fn := assetWrapper(route.Name, route.Path, route.HandlerFunc)
		log.Infof("Registered ASSET: %s for %s", route.Name, route.Path)

		router.PathPrefix(route.Path).Name(route.Name).Handler(fn)
	}

	for e, lis := range MiddlewareSet {
		for _, m := range lis {
			log.Infof("Registered Middleware: %s for %s", m.Name, e)
		}
	}

	if viper.IsSet("idm") {
		lis := viper.GetStringMapString("idm")
		for idm := range lis {
			ident.RegisterConfig(idm, viper.GetStringMapString("idm."+idm))
		}
	}

	for name, fn := range modules {
		module := viper.GetStringMapString("module." + name)
		fn(module)
	}

	return router
}
// Config reads settings from viper
func Config() {
	if viper.IsSet("http.fileserver") {
		fileserver := viper.GetString("http.fileserver")

		log.Infof("Configured: FileServer for %s", fileserver)
		s := strings.SplitN(fileserver, ":", 2)

		path := "/"
		dir := "./"

		if len(s) == 2 {
			dir = s[1]
			path = s[0]
		} else if len(s) == 1 {
			path = s[0]
		}

		AssetRegister("asset", AssetRoutes{
			{"Files", path, http.Dir(dir)},
		})
	}
}
// RegisterModule stores a module
func RegisterModule(name string, fn ModuleHandler) {
	name = strings.ToLower(name)

	modules[name] = fn
}
// Run startup a new server
func Run() {
	router := NewRouter()
	listen := viper.GetString("http.listen")
	log.Notice("Listen and Serve on ", listen)
	srv := &http.Server{
		Addr: listen,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}
	
	if viper.GetBool("http.tls") {
		cfg := &tls.Config{
        		MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
			    tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			    tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			    tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			    tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
    		}
		
		srv.TLSConfig = cfg
        	srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
		
		crt := viper.GetString("http.tls_crt")
		key := viper.GetString("http.tls_key")
		
		go func() {
			err := srv.ListenAndServeTLS(crt, key)
			if err != nil {
				log.Error(err)
			}
		}()
		
	} else {

		// Run our server in a goroutine so that it doesn't block.
		go func() {
			err := srv.ListenAndServe()
			if err != nil {
				log.Error(err)
			}

		}()

	}	
	close(SignalStartup)
}
// Shutdown graceful shutdown of server
func Shutdown() {
	close(SignalShutdown)
	log.Notice("shutting down...")
    done := make(chan struct{})
	go func() {
		WaitShutdown.Wait()
		close(done)
    }()

    select {
    case <-done:
    	log.Notice("all done.")
    case <-time.After(15 * time.Second):
		log.Notice("times up. forcing shutdown.")
	}
}

func init() {
	HttpRegister("info", HttpRoutes{
		{"get-app-info", "GET", "/app-info", getAppInfo},
		{"get-app-info", "GET", "/v1/app-info", v1GetAppInfo},
	})
}

// swagger:operation GET /app-info appInfo get-app-info
//
// Get App Info
//
// ---
// produces:
//   - "application/json"
// responses:
//   "200":
//     description: Success
//     schema:
//       type: object
//       properties:
//          items:
func getAppInfo(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	s := fmt.Sprintf("%s (%s %s)",
		viper.GetString("app.name"),
		viper.GetString("app.version"),
		viper.GetString("app.build"))

	w.Write([]byte(s))
}
// swagger:operation GET /v1/app-info appInfo v1-get-app-info
//
// Get App Info
//
// ---
// produces:
//   - "application/json"
// responses:
//   "200":
//     description: Success
//     schema:
//       type: object
//       properties:
//          items:
func v1GetAppInfo(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	app := viper.GetStringMapString("app")
	json.NewEncoder(w).Encode(app)
}

// ResultError is a message error
// swagger:model ResultError
type ResultError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
// WriteError write an error message
func WriteError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(ResultError{code, msg}); err != nil {
		log.Error(err)
	}
}
// WriteObject write object as json
func WriteObject(w http.ResponseWriter, code int, o interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(o); err != nil {
		log.Error(err)
	}
}
// WriteText writes plain text
func WriteText(w http.ResponseWriter, code int, o string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(o))
}
// ResultWindow represents a windowed struct of items.
// swagger:model ResultWindow
type ResultWindow struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Results uint64      `json:"results"`
	Limit   uint64      `json:"limit"`
	Offset  uint64      `json:"offset"`
	Items   interface{} `json:"items"`
}
// WriteWindow writes a window object of items
func WriteWindow(w http.ResponseWriter, code int, results, limit, offset uint64, o interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(
		ResultWindow{
			code,
			"OK",
			results,
			limit,
			offset,
			o}); err != nil {
		log.Error(err)
	}
}
