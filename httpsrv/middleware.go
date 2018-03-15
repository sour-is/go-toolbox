package httpsrv

import (
	"net/http"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
)

type Event int

const (
	EventPreAuth Event = iota
	EventPreHandle
	EventPostHandle
	EventUnknown = 1<<63 - 1
)

func (e Event) String() string {
	switch e {
	case EventPreAuth:
		return "EventPreAuth"
	case EventPreHandle:
		return "EventPreHandle"
	case EventPostHandle:
		return "EventPostHandle"
	case EventUnknown:
		fallthrough
	default:
		return "EventUnknown"
	}
}

type MiddlewareFunc func(name string, w ResponseWriter, r *http.Request, id ident.Ident)

type Middleware struct {
	Name        string
	Whitelist   map[string]bool
	Blacklist   map[string]bool
	ProcessHTTP MiddlewareFunc
}
type MiddlewareList []Middleware

var MiddlewareSet = make(map[Event][]Middleware)

func runMiddleware(e Event, name string, w ResponseWriter, r *http.Request, id ident.Ident) {
	for _, m := range MiddlewareSet[e] {
		if _, ok := m.Blacklist[name]; len(m.Whitelist) > 0 && ok {
			m.ProcessHTTP(name, w, r, id)
		}

		if _, ok := m.Whitelist[name]; len(m.Whitelist) == 0 || ok {
			m.ProcessHTTP(name, w, r, id)
		}
	}
}

func NewMiddleware(name string, hdlr MiddlewareFunc) (m Middleware) {
	m.Name = name
	m.ProcessHTTP = hdlr
	return m
}

func (m Middleware) SetWhitelist(whitelist []string) Middleware {
	lis := make(map[string]bool)
	for s := range lis {
		lis[s] = true
	}
	m.Whitelist = lis
	return m
}

func (m Middleware) SetBlacklist(whitelist []string) Middleware {
	lis := make(map[string]bool)
	for s := range lis {
		lis[s] = true
	}
	m.Blacklist = lis
	return m
}

func (m Middleware) Register(event Event) {
	MiddlewareSet[event] = append(MiddlewareSet[event], m)
}
