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
	"os"
	"os/signal"
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

	log.Notice("Listen and Serve on ", listen)

	wait := time.Second * 15

	srv := &http.Server{
		Addr: listen,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()
	close(SignalStartup)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	close(SignalShutdown)

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Notice("shutting down")
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

type ErrorMessage struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func WriteMsg(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorMessage{code, msg})
}
func WriteObject(w http.ResponseWriter, code int, o interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(o)
}
func WriteText(w http.ResponseWriter, code int, o string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(o))
}
