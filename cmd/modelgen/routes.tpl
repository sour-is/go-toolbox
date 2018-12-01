package routes

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"

    "github.com/gorilla/mux"

    ctrl "<...ctrl...>"
    model "<...model...>"
    
    "sour.is/x/toolbox/dbm/rsql/squirrel"
    "sour.is/x/toolbox/dbm"
    "sour.is/x/toolbox/httpsrv"
    "sour.is/x/toolbox/ident"

)

func init() {
    {{range .Types}}
    httpsrv.IdentRegister("{{.Name}}", httpsrv.IdentRoutes{
      {Name: "get-{{spineCase .Name}}", Method: "GET", Pattern: "/v1/{{spineCase .Name}}", HandlerFunc: get{{.Name}} },
      {Name: "get-{{spineCase .Name}}-by-id", Method: "GET", Pattern: "/v1/{{spineCase .Name}}({ids})", HandlerFunc: get{{.Name}}ByID},
      {{if .ROnly}}{{else}}
          {Name: "post-{{spineCase .Name}}", Method: "POST", Pattern: "/v1/{{spineCase .Name}}", HandlerFunc: post{{.Name}} },
          {Name: "put-{{spineCase .Name}}", Method: "PUT", Pattern: "/v1/{{spineCase .Name}}({ids})", HandlerFunc: put{{.Name}} },
          {Name: "delete-{{spineCase .Name}}", Method: "DELETE", Pattern: "/v1/{{spineCase .Name}}({ids})", HandlerFunc: delete{{.Name}} },
      {{end}}
    })
    {{end}}
}

{{range .Types}}
// swagger:operation GET /v1/{{spineCase .Name}} {{spineCase .Name}} get-{{spineCase .Name}}
//
// Get {{.Name}}
//
// ---
// parameters:
//   - name: search
//     in: query
//     description: Search
//     required: false
//     type: string
//     format: string
//   - name: limit
//     in: query
//     description: Limit
//     required: true
//     type: integer
//     format: int
//     default: 100
//   - name: offset
//     in: query
//     description: Offset
//     required: true
//     type: integer
//     format: int
//     default: 0
// consumes:
//   - "application/json"
// produces:
//   - "application/json"
// responses:
//   "200":
//     description: Success
//     schema:
//       type: object
//       allOf:
//       - "$ref": "#/definitions/ResultWindow"
//       - type: object
//         properties:
//           items:
//             type: array
//             items:
//               "$ref": "#/definitions/{{.Name}}"
//   "5xx":
//     description: unexpected error
//     schema:
//       "$ref": "#/definitions/ResultError"
func get{{.Name}}(w httpsrv.ResponseWriter, r *http.Request, user ident.Ident) {
    if !user.HasRole("admin", "read", "{{.Name}}", "get-{{spineCase .Name}}") {
        httpsrv.WriteError(w, http.StatusForbidden, "Access Denied")
        return
    }
    fn := ctrl.List{{.Name}}Count

    // ----
    var limit, offset uint64
    var err error

    if limit, err = strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64); err != nil {
    w.WriteError(http.StatusBadRequest, "invalid limit")
        return
    }
    if offset, err = strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64); err != nil {
        w.WriteError(http.StatusBadRequest, "invalid offset")
        return
    }

    search, err := squirrel.Query(r.URL.Query().Get("search"), dbm.GetDbInfo(model.{{.Name}}{}))
    if err != nil {
        w.WriteError(http.StatusBadRequest, err.Error())
        return
    }
    lis, count, err := fn(search, limit, offset)
    if err != nil {
        w.WriteError(http.StatusInternalServerError, err.Error())
        return
    }

    w.WriteWindow(http.StatusOK, count, limit, offset, lis)
}

// swagger:operation GET /v1/{{spineCase .Name}}({ids}) {{spineCase .Name}} get-{{.Name}}-by-id
//
// Get {{.Name}} by ID
//
// ---
// parameters:
//   - name: ids
//     in: path
//     description: {{.Name}} IDs
//     required: true
//     schema:
//        type: array
//        items:
//           type: int
// consumes:
//   - "application/json"
// produces:
//   - "application/json"
// responses:
//   "200":
//     description: Success
//     schema:
//       "$ref": "#/definitions/{{.Name}}"
//   "5xx":
//     description: unexpected error
//     schema:
//       "$ref": "#/definitions/ResultError"
func get{{.Name}}ByID(w httpsrv.ResponseWriter, r *http.Request, user ident.Ident) {
    if !user.HasRole("admin", "read", "{{.Name}}", "get-{{spineCase .Name}}-by-id") {
        httpsrv.WriteError(w, http.StatusForbidden, "Access Denied")
        return
    }
    fn := ctrl.List{{.Name}}ByID

    // ----
    vars := mux.Vars(r)
    var ids []uint64
    for _, id := range strings.Split(vars["ids"], ",") {
        i, err := strconv.ParseUint(id, 10, 64)
        if err != nil {
            w.WriteError(http.StatusBadRequest, fmt.Sprintf("invalid id: %v", id))
            return
        }
        ids = append(ids, i)
    }
    lis, count, err := fn(ids)
    if err != nil {
        w.WriteError(http.StatusInternalServerError, err.Error())
        return
    }
    w.WriteWindow(http.StatusOK, count, count, 0, lis)
}

{{if .ROnly}}{{else}}
// swagger:operation POST /v1/{{spineCase .Name}} {{spineCase .Name}} post-{{spineCase .Name}}
//
// Post {{.Name}}
//
// ---
// parameters:
//   - name: {{.Name}}
//     in: body
//     description: {{.Name}}
//     required: true
//     schema:
//        "$ref": "#/definitions/{{.Name}}"
// consumes:
//   - "application/json"
// produces:
//   - "application/json"
// responses:
//   "200":
//     description: Success
//     schema:
//       "$ref": "#/definitions/{{.Name}}"
//   "5xx":
//     description: unexpected error
//     schema:
//       "$ref": "#/definitions/ResultError"
func post{{.Name}}(w httpsrv.ResponseWriter, r *http.Request, user ident.Ident) {
    if !user.HasRole("admin", "write", "{{.Name}}", "post-{{spineCase .Name}}") {
        httpsrv.WriteError(w, http.StatusForbidden, "Access Denied")
        return
    }

    o, err := ctrl.Put{{.Name}}(0, func(mode ctrl.Mode, d dbm.DbInfo, o ctrl.{{.Name}}) (err error) {
        // ----

        if err = json.NewDecoder(r.Body).Decode(o.{{.Name}}); err != nil {
            return ctrl.ParseError(err.Error())
        }

        // ----
        err = o.Save()
        return
    })
    if err != nil {
        switch err.(type) {
        case ctrl.ParseError:
            w.WriteError(http.StatusBadRequest, err.Error())
        default:
            w.WriteError(http.StatusInternalServerError, err.Error())
        }
        return
    }
    w.WriteObject(http.StatusOK, o)
}

// swagger:operation PUT /v1/{{spineCase .Name}}({id}) {{spineCase .Name}} put-{{spineCase .Name}}
//
// Post {{.Name}}
//
// ---
// parameters:
//   - name: id
//     in: path
//     description: {{.Name}} ID
//     required: true
//     schema:
//        type: int
//   - name: {{.Name}}
//     in: body
//     description: {{.Name}}
//     required: true
//     schema:
//        "$ref": "#/definitions/{{.Name}}"
// consumes:
//   - "application/json"
// produces:
//   - "application/json"
// responses:
//   "200":
//     description: Success
//     schema:
//       "$ref": "#/definitions/{{.Name}}"
//   "5xx":
//     description: unexpected error
//     schema:
//       "$ref": "#/definitions/ResultError"
func put{{.Name}}(w httpsrv.ResponseWriter, r *http.Request, user ident.Ident) {
    if !user.HasRole("admin", "write", "{{.Name}}", "put-{{spineCase .Name}}") {
        httpsrv.WriteError(w, http.StatusForbidden, "Access Denied")
        return
    }

    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 64)
    if err != nil {
        w.WriteError(http.StatusBadRequest, fmt.Sprintf("invalid id: %v", id))
        return
    }

    o, err := ctrl.Put{{.Name}}(id, func(mode ctrl.Mode, d dbm.DbInfo, o ctrl.{{.Name}}) (err error) {
        // ----

        if err = json.NewDecoder(r.Body).Decode(o.{{.Name}}); err != nil {
            return ctrl.ParseError(err.Error())
        }

        // ----
        err = o.Save()
        return
    })
    if err != nil {
        switch err.(type) {
        case ctrl.ParseError:
            w.WriteError(http.StatusBadRequest, err.Error())
        default:
            w.WriteError(http.StatusInternalServerError, err.Error())
        }
        return
    }
    w.WriteObject(http.StatusOK, o)
}

// swagger:operation DELETE /v1/{{spineCase .Name}}({id}) {{spineCase .Name}} delete-{{spineCase .Name}}
//
// Delete {{.Name}}
//
// ---
// parameters:
//   - name: id
//     in: path
//     description: {{.Name}} ID
//     required: true
//     schema:
//        type: int
// consumes:
//   - "application/json"
// produces:
//   - "application/json"
// responses:
//   "200":
//     description: Success
//     schema:
//       "$ref": "#/definitions/{{.Name}}"
//   "5xx":
//     description: unexpected error
//     schema:
//       "$ref": "#/definitions/ResultError"
func delete{{.Name}}(w httpsrv.ResponseWriter, r *http.Request, user ident.Ident) {
    if !user.HasRole("admin", "write", "{{.Name}}", "put-{{spineCase .Name}}") {
        httpsrv.WriteError(w, http.StatusForbidden, "Access Denied")
        return
    }

    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 64)
    if err != nil {
        w.WriteError(http.StatusBadRequest, fmt.Sprintf("invalid id: %v", id))
        return
    }

    err = ctrl.Delete{{.Name}}ByID(id)
    if err != nil {
        w.WriteError(http.StatusInternalServerError, err.Error())
        return
    }
    w.WriteText(http.StatusGone, "OK")
}
{{end}}
{{end}}
