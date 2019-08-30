package httpsrv

import (
	"net/http"

	"sour.is/x/toolbox/ident"
)

// Event lifecycle state
type Event int

// Event lifecycle states
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

// MiddlewareFunc is a function that handles a request/response through request lifecycle
type MiddlewareFunc func(name string, w ResponseWriter, r *http.Request, id ident.Ident) bool

// Middleware defines when middleware should be executed in event lifecycle
type Middleware struct {
	Name        string
	Whitelist   map[string]bool
	Blacklist   map[string]bool
	ProcessHTTP MiddlewareFunc
}

// MiddlewareList is a list of registered middlewares
type MiddlewareList []Middleware

// MiddlewareSet is a set of middlewares grouped by Event lifecycle
var MiddlewareSet = make(map[Event][]Middleware)

func runMiddleware(e Event, name string, w ResponseWriter, r *http.Request, id ident.Ident) (ok bool) {
	ok = true

	for _, m := range MiddlewareSet[e] {

		// WL0  WL?  BL0  BL? :: RUN
		//  Y    -    Y    -  ::  Y
		//  Y    -    N    N  ::  Y
		//  Y    -    N    Y  ::  N

		if len(m.Whitelist) == 0 {
			if _, ck := m.Blacklist[name]; ck {
				//log.Debugf("Event %s is in blacklist: %v", name, m.Blacklist)
				return
			}
		}

		//  Y    -    Y    -  ::  Y
		//  Y    -    N    N  ::  Y
		//  N    N    Y    -  ::  N
		//  N    N    N    N  ::  N
		//  N    N    N    Y  ::  N

		if len(m.Whitelist) > 0 {
			if _, ck := m.Whitelist[name]; !ck {
				//log.Debugf("Event %s is NOT in whitelist: %v", name, m.Whitelist)
				return
			}
		}

		//  Y    -    Y    -  ::  Y
		//  Y    -    N    N  ::  Y
		//  N    Y    Y    -  ::  Y
		//  N    Y    N    N  ::  Y
		//  N    Y    N    Y  ::  N

		if _, ck := m.Blacklist[name]; ck {
			// log.Debugf("Event %s is in blacklist: %v", name, m.Blacklist)
			return
		}

		//log.Debugf("Event %s passes white/black lists", name)
		ok = m.ProcessHTTP(name, w, r, id)
	}
	return
}

// NewMiddleware defines a new middleware
func NewMiddleware(name string, hdlr MiddlewareFunc) (m Middleware) {
	m.Name = name
	m.ProcessHTTP = hdlr
	m.Whitelist = make(map[string]bool)
	m.Blacklist = make(map[string]bool)
	return m
}

// SetWhitelist sets the events names to target
func (m Middleware) SetWhitelist(whitelist []string) Middleware {
	for _, s := range whitelist {
		m.Whitelist[s] = true
	}

	return m
}

// SetBlacklist sets the events to avoid
func (m Middleware) SetBlacklist(whitelist []string) Middleware {
	for _, s := range whitelist {
		m.Blacklist[s] = true
	}

	return m
}

// Register inserts the middleware into the lifecycle map
func (m Middleware) Register(event Event) {
	MiddlewareSet[event] = append(MiddlewareSet[event], m)
}
