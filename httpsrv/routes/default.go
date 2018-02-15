package routes

import (
	"fmt"
	"net/http"
	"sour.is/x/httpsrv"
	"sour.is/x/ident"
)

func init() {
	httpsrv.IdentRegister("default", httpsrv.IdentRoutes{
		{"Index", "GET", "/", Index},
	})
}

// swagger:route GET / defaultIndex
//
// Welcomes user.
//
// This welcome user based on identity used.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Security:
//       api_key:
//
//     Responses:
//       default: genericError
//       200: someResponse
//       422: validationError
func Index(w http.ResponseWriter, r *http.Request, i ident.Ident) {
	fmt.Fprintf(w, "Welcome, %s (%s)!\n", i.GetDisplay(), i.GetIdentity())
}
