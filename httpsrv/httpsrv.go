package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"sort"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
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

type ModuleHandler func(map[string]string)

var modules = make(map[string]ModuleHandler)
var server *http.Server

var SignalStartup = make(chan struct{})
var SignalShutdown = make(chan struct{})
var WaitShutdown sync.WaitGroup

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

func RegisterModule(name string, fn ModuleHandler) {
	name = strings.ToLower(name)

	modules[name] = fn
}

func Run() {
	router := NewRouter()
	listen := viper.GetString("http.listen")

	log.Notice("Listen and Serve on", listen)

	server = new(http.Server)
	server.Addr = listen
	server.Handler = router

	close(SignalStartup)
	log.Notice(server.ListenAndServe())
}

func Shutdown() {
	close(SignalShutdown)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	go server.Shutdown(ctx)

	select {
	case <-ctx.Done():
		log.Error(ctx.Err()) // prints "context deadline exceeded"
		server.Close()
	}

	WaitShutdown.Wait()
}

func init() {
	HttpRegister("info", HttpRoutes{
		{"GetAppInfo", "GET", "/app-info", getAppInfo},
		{"GetAppInfo", "GET", "/v1/app-info", v1GetAppInfo},
	})
}

func getAppInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	s := fmt.Sprintf("%s (%s %s)",
		viper.GetString("app.name"),
		viper.GetString("app.version"),
		viper.GetString("app.build"))

	w.Write([]byte(s))
}

func v1GetAppInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	app := viper.GetStringMapString("app")
	json.NewEncoder(w).Encode(app)
}
