package mercury

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"sour.is/x/toolbox/dbm/rsql"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
)

// Handler interface for backends
type Handler interface {
	GetIndex(NamespaceSearch, *rsql.Program) ArraySpace
	GetObjects(NamespaceSearch, *rsql.Program, []string) ArraySpace
	WriteObjects(ArraySpace) error
	GetRules(ident.Ident) Rules
	GetNotify(string) ListNotify
}

// HandlerItem a single handler matching
type HandlerItem struct {
	Handler
	Match    string
	Priority int
}

// HandlerList a list of handlers
type HandlerList []HandlerItem

func (h HandlerItem) String() string {
	return fmt.Sprintf("%d: %s", h.Priority, h.Match)
}

// Registry handler
var Registry HandlerList

func (hl HandlerList) String() string {
	var buf strings.Builder
	for _, h := range hl {
		buf.WriteString(h.String())
		buf.WriteRune('\n')
	}

	return buf.String()
}

// Register add a handler to registry
func Register(match string, priority int, hdlr Handler) {
	fmt.Println("mercury regster ", match)
	Registry = append(Registry, HandlerItem{Match: match, Priority: priority, Handler: hdlr})
	sort.Sort(Registry)
}

// Len implements Len for sort.interface
func (hl HandlerList) Len() int {
	return len(hl)
}

// Less implements Less for sort.interface
func (hl HandlerList) Less(i, j int) bool {
	return hl[i].Priority > hl[j].Priority
}

// Swap implements Swap for sort.interface
func (hl HandlerList) Swap(i, j int) { hl[i], hl[j] = hl[j], hl[i] }

// GetIndex query each handler that match namespace.
func (hl HandlerList) GetIndex(match, search string) (lis ArraySpace) {
	spec := ParseNamespace(match)
	pgm := rsql.DefaultParse(search)
	matches := make([]NamespaceSearch, len(hl))

	for _, c := range spec {
		for i, hldr := range hl {

			ok, err := filepath.Match(hldr.Match, c.Value())
			if !ok || err != nil {
				continue
			}
			matches[i] = append(matches[i], c)
		}
	}

	for i, hldr := range hl {
		log.Debug("INDEX ", hldr.Match)
		arr := hldr.GetIndex(matches[i], pgm)
		if arr != nil {
			lis = append(lis, arr...)
		}
	}

	return
}

// Search query each handler with a key=value search

// GetObjects query each handler that match for fully qualified namespaces.
func (hl HandlerList) GetObjects(match, search, fields string) (out ArraySpace) {
	spec := ParseNamespace(match)
	pgm := rsql.DefaultParse(search)
	flds := strings.Split(fields, ",")

	matches := make([]NamespaceSearch, len(hl))

	for _, c := range spec {
		for i, hldr := range hl {
			ok, err := filepath.Match(hldr.Match, c.Value())
			if !ok || err != nil {
				continue
			}
			matches[i] = append(matches[i], c)
		}
	}

	for i, hldr := range hl {
		log.Debug("QUERY ", hldr.Match)
		arr := hldr.GetObjects(matches[i], pgm, flds)
		if arr != nil {
			out = append(out, arr...)
		}
	}

	return
}

// WriteObjects write objects to backends
func (hl HandlerList) WriteObjects(spaces ArraySpace) error {
	matches := make([]ArraySpace, len(hl))

	for _, s := range spaces {
		for i, hldr := range hl {
			ok, err := filepath.Match(hldr.Match, s.Space)
			if !ok || err != nil {
				continue
			}
			log.Debug("MATCH ", i, " ", s.Space)
			matches[i] = append(matches[i], s)
			break
		}
	}

	for i, hldr := range hl {
		log.Debug("WRITE MATCH ", hldr.Match)
		err := hldr.WriteObjects(matches[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// GetRules query each of the handlers for rules.
func (hl HandlerList) GetRules(user ident.Ident) (lis Rules) {
	for _, hldr := range hl {
		arr := hldr.GetRules(user)
		if arr != nil {
			log.Debug("RULES ", hldr.Match)
			lis = append(lis, arr...)
		}
	}

	return
}

// GetNotify query each of the handlers for rules.
func (hl HandlerList) GetNotify(event string) (lis ListNotify, err error) {
	for _, hldr := range hl {
		log.Debug("NOTIFY ", hldr.Match)

		arr := hldr.GetNotify(event)
		if err != nil {
			continue
		}
		if arr != nil {
			log.Debug("NOTIFY ", hldr.Match)
			lis = append(lis, arr...)
		}
	}

	return
}
