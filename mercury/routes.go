package mercury

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/BurntSushi/toml"

	"sour.is/x/toolbox/httpsrv"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"

	"github.com/golang/gddo/httputil"
)

func init() {
	httpsrv.IdentRegister("config", httpsrv.IdentRoutes{
		{Name: "get-mercury-spaces", Method: "GET", Pattern: "/v1/mercury-spaces", HandlerFunc: getSpace},

		{Name: "get-mercury-config", Method: "GET", Pattern: "/v1/mercury-config", HandlerFunc: getConfig},
		{Name: "post-mercury-config", Method: "POST", Pattern: "/v1/mercury-config", HandlerFunc: postConfig},
	})
}

// swagger:operation GET /v1/mercury-config mercury get-mercury-config
//
// Get Mercury Config
//
// ---
// parameters:
//   - name: space
//     in: query
//     description: Space
//     required: false
//     type: string
//     format: string
// consumes:
//   - "application/json"
// produces:
//   - "text/plain"
//   - "application/environ"
//   - "application/ini"
//   - "application/json"
//   - "application/toml"
// responses:
//   "200":
//     description: Success
//     schema:
//       type: object
//       allOf:
//       - "$ref": "#/definitions/Space"
//       - type: object
//         properties:
//           items:
//             type: array
//             items:
//               "$ref": "#/definitions/Audience"
//   "5xx":
//     description: unexpected error
//     schema:
//       "$ref": "#/definitions/ResultError"
func getConfig(w httpsrv.ResponseWriter, r *http.Request, id ident.Ident) {
	if !id.IsActive() {
		w.WriteError(401, "NO_AUTH")
		return
	}

	rules := Registry.GetRules(id)
	space := r.URL.Query().Get("space")
	if space == "" {
		space = "*"
	}

	ns := ParseNamespace(space)
	ns = rules.ReduceSearch(ns)

	var err error
	lis := Registry.GetObjects(ns.String(), "", "")
	lis, err = lis.accessFilter(id)
	if err != nil {
		w.WriteError(500, "ERR: "+err.Error())
		return
	}

	sort.Sort(lis)
	var content string

	switch httputil.NegotiateContentType(r, []string{
		"text/plain",
		"application/environ",
		"application/ini",
		"application/json",
		"application/toml",
	}, "text/plain") {
	case "text/plain":
		content = lis.String()
	case "application/environ":
		content = lis.EnvString()
	case "application/ini":
		content = lis.IniString()
	case "application/json":
		w.WriteObject(200, lis)
	case "application/toml":
		w.WriteText(200, "")
		m := make(map[string]map[string][]string)
		for _, o := range lis {
			if _, ok := m[o.Space]; !ok {
				m[o.Space] = make(map[string][]string)
			}
			for _, v := range o.List {
				m[o.Space][v.Name] = append(m[o.Space][v.Name], v.Values...)
			}
		}
		err := toml.NewEncoder(w.W).Encode(m)
		if err != nil {
			log.Error(err)
		}
	}

	w.WriteText(200, content)
}

// swagger:operation POST /v1/mercury-config mercury post-mercury-config
//
// Get Mercury
//
// ---
// parameters:
//   - name: payload
//     in: post
//     description: Payload
//     required: true
//     type: string
//     format: string
// consumes:
//   - "application/json"
// produces:
//   - "application/json"
// responses:
//   "200":
//     description: Success
//     schema:
//       type: string
//   "5xx":
//     description: unexpected error
//     schema:
//       "$ref": "#/definitions/ResultError"
func postConfig(w httpsrv.ResponseWriter, r *http.Request, id ident.Ident) {
	if !id.IsActive() {
		w.WriteError(401, "NO_AUTH")
		return
	}
	config, err := parseText(r.Body)
	r.Body.Close()
	if err != nil {
		w.WriteError(400, "PARSE_ERR")
		return
	}

	c, _ := json.MarshalIndent(config, "", "  ")
	log.Debug(string(c))

	func() {
		rules := Registry.GetRules(id)

		notify, err := Registry.GetNotify("updated")
		if err != nil {
			log.Error(err)
		}
		_ = rules
		var notifyActive = make(map[string]struct{})
		var filteredConfigs Config
		for ns, c := range config {
			if !rules.GetRoles("NS", ns).HasRole("write") {
				log.Debug("SKIP", ns)
				continue
			}

			log.Debug("SAVE", ns)
			for _, n := range notify.Find(ns) {
				notifyActive[n.Name] = struct{}{}
			}
			filteredConfigs = append(filteredConfigs, c)
		}

		err = Registry.WriteObjects(filteredConfigs)
		if err != nil {
			log.Error(err)
			return
		}

		log.Debug("SEND NOTIFYS ", notifyActive)
		for _, n := range notify {
			if _, ok := notifyActive[n.Name]; ok {
				err = n.sendNotify()
				if err != nil {
					log.Debug(err)
				}
			}
		}
		log.Debug("DONE!")
	}()

	w.WriteText(202, "OK")
}

// swagger:operation GET /v1/mercury-spaces mercury get-mercury-spaces
//
// Get Mercury Space List
//
// ---
// parameters:
//   - name: space
//     in: query
//     description: Space
//     required: false
//     type: string
//     format: string
// consumes:
//   - "application/json"
// produces:
//   - "text/plain"
//   - "application/json"
// responses:
//   "200":
//     description: Success
//     schema:
//       type: object
//       allOf:
//       - "$ref": "#/definitions/Space"
//       - type: object
//         properties:
//           items:
//             type: array
//             items:
//               "$ref": "#/definitions/Audience"
//   "5xx":
//     description: unexpected error
//     schema:
//       "$ref": "#/definitions/ResultError"

func getSpace(w httpsrv.ResponseWriter, r *http.Request, id ident.Ident) {
	if !id.IsActive() {
		w.WriteError(401, "NO_AUTH")
		return
	}

	rules := Registry.GetRules(id)
	log.Debug(rules)

	space := r.URL.Query().Get("space")
	if space == "" {
		space = "*"
	}

	ns := ParseNamespace(space)
	ns = rules.ReduceSearch(ns)
	log.Debug(ns.String())

	lis := Registry.GetIndex(ns.String(), "")
	sort.Sort(lis)

	switch httputil.NegotiateContentType(r, []string{
		"text/plain",
		"application/json",
	}, "text/plain") {
	case "text/plain":
		w.WriteText(200, lis.StringList())
	case "application/json":
		w.WriteObject(200, lis)
	}
}
