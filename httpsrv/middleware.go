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

		// WL0  WL?  BL0  BL? :: RUN
		//  Y    -    Y    -  ::  Y
		//  Y    -    N    N  ::  Y
		//  Y    -    N    Y  ::  N

		if len(m.Whitelist) == 0 {
			if _, ck := m.Blacklist[name]; ck {
				log.Debugf("Event %s is in blacklist: %v", name, m.Blacklist)
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
				log.Debugf("Event %s is NOT in whitelist: %v", name, m.Whitelist)
				return
			}
		}

		//  Y    -    Y    -  ::  Y
		//  Y    -    N    N  ::  Y
		//  N    Y    Y    -  ::  Y
		//  N    Y    N    N  ::  Y
		//  N    Y    N    Y  ::  N

		if _, ck := m.Blacklist[name]; ck {
			log.Debugf("Event %s is in blacklist: %v", name, m.Blacklist)

			return
		}

		log.Debugf("Event %s passes white/black lists", name)
		ok = m.ProcessHTTP(name, w, r, id)
	}
	return
}

func NewMiddleware(name string, hdlr MiddlewareFunc) (m Middleware) {
	m.Name = name
	m.ProcessHTTP = hdlr
	m.Whitelist = make(map[string]bool)
	m.Blacklist = make(map[string]bool)
	return m
}

func (m Middleware) SetWhitelist(whitelist []string) Middleware {
	for _, s := range whitelist {
		m.Whitelist[s] = true
	}

	return m
}

func (m Middleware) SetBlacklist(whitelist []string) Middleware {
	for _, s := range whitelist {
		m.Blacklist[s] = true
	}

	return m
}

func (m Middleware) Register(event Event) {
	MiddlewareSet[event] = append(MiddlewareSet[event], m)
}
