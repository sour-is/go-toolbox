
{{range .Types}}//go:generate dataloaden <...model...>.{{.Name}}
{{end}}

package loader

import (
	"context"
	"fmt"
	"time"
	"sync"

	ctrl "<...ctrl...>"
	model "<...model...>"

	"sour.is/x/toolbox/ident"
)

type contextKey struct {
   name string
}
func (c *contextKey) String() string {
    return "loader context key " + c.name
}

// ManagerKey defines the context key for finding the dataloader
var ManagerKey = contextKey{"loader.Manger"}

// Manager holds all of the loaders for lazy loading from database
type Manager struct {
    mu sync.RWMutex
    Ident     ident.Ident
{{range .Types}}
{{if .HasID}}{{.Name}}    {{.Name}}Loader{{end}}
{{end}}
}
func (m *Manager) GetIdent() *ident.Ident {
    m.mu.RLock()
    defer m.mu.RUnlock()

    return &m.Ident
}
func (m *Manager) SetIdent(u ident.Ident) {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.Ident = u
}

func Get(id ident.Ident) *Manager {
	return GetContext(context.Background(), id)
}
func GetContext(ctx context.Context, id ident.Ident) *Manager {
	return &Manager{
		Ident: id,

        {{range .Types}}
		{{if .HasID}}
        {{.Name}}: {{.Name}}Loader{
			maxBatch: 1000,
			wait:     75 * time.Millisecond,
			fetch: fetch{{.Name}}(ctx),
		},
		{{end}}
        {{end}}
	}
}

{{range .Types}}
{{if .HasID}}
func fetch{{.Name}}(ctx context.Context) func(ids []int) (lisP []*model.{{.Name}}, errs []error) {
	return func(ids []int) (lisP []*model.{{.Name}}, errs []error) {
		fn := ctrl.List{{.Name}}ByIDContext
		null := model.{{.Name}}{}
		lisP = make([]*model.{{.Name}}, len(ids))

		// ----
		m := make(map[uint64]int, len(ids))

		args := make([]uint64, len(ids))
		for i := 0; i < len(ids); i++ {
			args[i] = uint64(ids[i])
			m[uint64(ids[i])] = i
		}
		lis, _, err := fn(ctx, args)

		if err != nil {
			errs = append(errs, err)
			return
		}

		for i, o := range lis {
			lisP[m[o.ID]] = &lis[i]
		}

		for i, o := range lisP {
			if o == nil {
				lisP[i] = &null
				errs = append(errs, fmt.Errorf("not found: %d", i))
			}
		}
		return	
	}
}
{{end}}
{{end}}
