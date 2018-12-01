
package resolver

import (
	"context"

	ctrl "<...ctrl...>"
	loader "<...loader...>"
	model "<...model...>"

    "sour.is/x/toolbox/dbm/qry"
	opentracing "github.com/opentracing/opentracing-go"
)

{{range .Types}}
// nolint: deadcode
func get{{.Name}}(ctx context.Context, id uint64) (model.{{.Name}}, error) {
    sp, _ := opentracing.StartSpanFromContext(ctx, "get{{.Name}}List")
    defer sp.Finish()

    ptr, err := ctx.Value(loader.ManagerKey).(*loader.Manager).{{.Name}}.Load(int(id))
    return *ptr, err
}
// nolint: deadcode
func get{{.Name}}List(ctx context.Context, q qry.Input) (lis []model.{{.Name}}, err error) {
    sp, octx := opentracing.StartSpanFromContext(ctx, "get{{.Name}}List")
    defer sp.Finish()

	// ----
	
    return ctrl.List{{.Name}}Qry(octx, q)
}
{{end}}
