package mercury

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"path/filepath"

	"sour.is/x/toolbox/log"
	"sour.is/x/toolbox/mqtt"
)

func (n Notify) sendNotify() (err error) {
	if n.Method == "MQTT" {
		var m mqtt.Message
		m, err = mqtt.NewMessage(n.URL, n)
		if err != nil {
			return
		}
		log.Debug(n)
		err = mqtt.Publish(m)
		return
	}

	cl := &http.Client{}
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		caCertPool = x509.NewCertPool()
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	cl.Transport = transport

	var req *http.Request
	req, err = http.NewRequest(n.Method, n.URL, bytes.NewBufferString(""))
	if err != nil {
		return
	}
	req.Header.Set("content-type", "application/json")

	log.Notice("URL: ", n.URL)
	res, err := cl.Do(req)
	if err != nil {
		return
	}
	res.Body.Close()
	log.Debug(res.Status)
	if res.StatusCode != 200 {
		err = fmt.Errorf("unable to read config")
		return
	}

	return
}

// Check if name matches notify
func (n Notify) Check(name string) bool {
	ok, err := filepath.Match(n.Match, name)
	if err != nil {
		return false
	}
	return ok
}

// Notify stores the attributes for a registry space
type Notify struct {
	Name   string
	Match  string
	Event  string
	Method string
	URL    string
}

// ListNotify array of notify
type ListNotify []Notify

// Find returns list of notify that match name.
func (ln ListNotify) Find(name string) (lis ListNotify) {
	lis = make(ListNotify, 0, len(ln))
	for _, o := range ln {
		if o.Check(name) {
			lis = append(lis, o)
		}
	}
	return
}
