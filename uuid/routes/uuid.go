package routes

import (
	"fmt"
	"net/http"
	"sour.is/x/httpsrv"
	"sour.is/x/uuid"
)

func init() {
	httpsrv.HttpRegister("uuid", httpsrv.HttpRoutes{
		{
			"v4",
			"GET",
			"/v1/uuid/v4",
			v4,
		},
		{
			"v5",
			"GET",
			"/v1/uuid/v5",
			v5,
		},
		{
			"v6",
			"GET",
			"/v1/uuid/v6",
			v6,
		},
	})
}

func v4(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, uuid.V4())
}

func v5(w http.ResponseWriter, r *http.Request) {
	var ns string
	var name string

	if ns = r.URL.Query().Get("ns"); ns == "" {
		ns = "00000000-0000-0000-0000-000000000000"
	}

	if name = r.URL.Query().Get("name"); name == "" {
		name = "text"
	}

	fmt.Fprint(w, uuid.V5(name, ns))
}

func v6(w http.ResponseWriter, r *http.Request) {
	var ns string
	var origin bool

	if ns = r.URL.Query().Get("ns"); ns == "" {
		ns = "00000000-0000-0000-0000-000000000000"
	}

	origin = false
	if ok := r.URL.Query().Get("origin"); ok == "true" {
		origin = true
	}

	fmt.Fprint(w, uuid.V6(ns, origin))
}
