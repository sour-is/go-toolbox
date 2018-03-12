package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
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

			handler := Wrapper(route.HandlerFunc, name)

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

			handler := IdentWrapper(route.HandlerFunc, name)

			log.Infof("Registered IDENT: %s for %s %s",
				name, route.Method, route.Pattern)

			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(name).
				Handler(handler)
		}
	}

	for setName, assets := range AssetSet {
		for _, route := range assets {
			name := setName + "::" + route.Name

			fn := AssetWrapper(name, route.Path, route.HandlerFunc)
			log.Infof("Registered ASSET: %s for %s", name, route.Path)

			router.PathPrefix(route.Path).Name(name).Handler(fn)
		}
	}

	if viper.IsSet("idm") {
		lis := viper.GetStringMapString("idm")
		for idm, _ := range lis {
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

	/* TODO This requires Go 1.8+ */
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	go server.Shutdown(ctx)

	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // prints "context deadline exceeded"
		server.Close()
	}

	WaitShutdown.Wait()
}

func init() {
	HttpRegister("info", HttpRoutes{
		{"GetAppInfo", "GET", "/app-info", GetAppInfo},
	})
}

func GetAppInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	s := fmt.Sprintf("%s (%s %s)",
		viper.GetString("app.name"),
		viper.GetString("app.version"),
		viper.GetString("app.build"))

	w.Write([]byte(s))
}
