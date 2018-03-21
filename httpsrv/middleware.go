package httpsrv

import (
	"net/http"
	"sour.is/x/toolbox/ident"
)

type Event int

const (
	EventPreAuth Event = iota
	EventPreHandle
	EventPostHandle
	EventComplete
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
	case EventComplete:
		return "EventComplete"
	case EventUnknown:
		fallthrough
	default:
		return "EventUnknown"
	}
}

type MiddlewareFunc func(name string, w ResponseWriter, r *http.Request, id ident.Ident) bool

type Middleware struct {
	Name        string
	Whitelist   map[string]bool
	Blacklist   map[string]bool
	ProcessHTTP MiddlewareFunc
}
type MiddlewareList []Middleware

var MiddlewareSet = make(map[Event][]Middleware)

func runMiddleware(e Event, name string, w ResponseWriter, r *http.Request, id ident.Ident) (ok bool) {
	ok = true

	for _, m := range MiddlewareSet[e] {
		if _, ck := m.Blacklist[name]; len(m.Whitelist) > 0 && ck {
			ok = m.ProcessHTTP(name, w, r, id)
			if !ok {
				return
			}
		}

		if _, ck := m.Whitelist[name]; len(m.Whitelist) == 0 || ck {
			ok = m.ProcessHTTP(name, w, r, id)
			if !ok {
				return
			}
		}
	}
	return
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
