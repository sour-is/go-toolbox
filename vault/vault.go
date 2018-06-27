package vault

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"sour.is/x/toolbox/httpsrv"
	"sour.is/x/toolbox/log"
)

type vaultAuth struct {
	Auth   vaultAuthMap `json:"auth"`
	Data   vaultAuthMap `json:"data"`
	Errors []string     `json:"errors"`
}
type vaultData struct {
	Data   vaultDataMap `json:"data"`
	Errors []string     `json:"errors"`
}

type vaultDataMap map[string]interface{}
type vaultAuthMap struct {
	Accessor      string            `json:"accessor"`
	ClientToken   string            `json:"client_token"`
	EntityID      string            `json:"entity_id"`
	LeaseDuration int               `json:"lease_duration"`
	Metadata      map[string]string `json:"metadata"`
	Policies      []string          `json:"policies"`
	Renewable     bool              `json:"renewable"`
	TTL           int               `json:"ttl"`
}

type pki struct {
	Addr     string
	CAfile   string
	CA       string
	CertFile string
	Cert     string
	KeyFile  string
	Key      string
}

var vault struct {
	Addr   string
	Secret string
	Token  string
	CA     string
	Cert   string
	Key    string
	PKI    pki
}

const (
	methodGET       = "GET"
	methodPOST      = "POST"
	vaultPKIAuth    = "v1/auth/cert/login"
	vaultLookupSelf = "v1/auth/token/lookup-self"
	vaultRenewSelf  = "v1/auth/token/renew-self"
)

func config() error {
	// Read VAULT specific config environment variables.
	viper.BindEnv("vault.addr",   "VAULT_ADDR")
	viper.BindEnv("vault.secret", "VAULT_SECRET")
	viper.BindEnv("vault.token",  "VAULT_TOKEN")
	viper.BindEnv("vault.ca",     "VAULT_CAFILE")
	viper.BindEnv("vault.cert",   "VAULT_CLIENT_CERT")
	viper.BindEnv("vault.key",    "VAULT_CLIENT_KEY")

	vault.Addr = viper.GetString("vault.addr")
	vault.Secret = viper.GetString("vault.secret")

	if viper.IsSet("vault.ca") {
		b, err := ioutil.ReadFile(viper.GetString("vault.ca"))
		if err != nil {
			log.Fatal(err)
		}
		vault.CA = string(b)
	}

	if viper.IsSet("vault.cert") {
		b, err := ioutil.ReadFile(viper.GetString("vault.cert"))
		if err != nil {
			log.Fatal(err)
		}
		vault.Cert = string(b)
	}

	if viper.IsSet("vault.key") {
		b, err := ioutil.ReadFile(viper.GetString("vault.key"))
		if err != nil {
			log.Fatal(err)
		}
		vault.Key = string(b)
	}

	if viper.IsSet("vault.token") {
		vault.Token = viper.GetString("vault.token")
		log.Noticef("Using Token %v****", vault.Token[:9])
	} else {
		log.Debug("Attempting PKI for auth ...")

		err := viper.UnmarshalKey("vault.pki", &vault.PKI)
		if err != nil {
			return err
		}

		if vault.PKI.Addr == "" {
			vault.PKI.Addr = vault.Addr
		}

		if vault.PKI.CAfile != "" {
			b, err := ioutil.ReadFile(vault.PKI.CAfile)
			if err != nil {
				log.Fatal(err)
			}
			vault.PKI.CA = string(b)
		}
		if vault.PKI.CA == "" {
			vault.PKI.CA = vault.CA
		}

		if vault.PKI.CertFile != "" {
			b, err := ioutil.ReadFile(vault.PKI.CertFile)
			if err != nil {
				log.Fatal(err)
			}
			vault.PKI.Cert = string(b)
		}
		if vault.PKI.Cert == "" {
			vault.PKI.Cert = vault.Cert
		}

		if vault.PKI.KeyFile != "" {
			b, err := ioutil.ReadFile(vault.PKI.KeyFile)
			if err != nil {
				log.Fatal(err)
			}
			vault.PKI.Key = string(b)
		}
		if vault.PKI.Key == "" {
			vault.PKI.Key = vault.Key
		}

		return certAuth(vault.PKI)
	}

	return checkToken()
}

func LoadVault() error {
	err := config()
	if err != nil {
		return err
	}

	if vault.Secret == "" {
		log.Debug("No Secret defined for vault.")
		return nil
	}

	log.Noticef("Read config from: %s secret: %s", vault.Addr, vault.Secret)

	cl := newClient(vault.CA, "", "")
	cl.Token = vault.Token
	data, err := cl.Req(methodGET, fmt.Sprintf("%s/v1/%s", vault.Addr, vault.Secret))
	if err != nil {
		return err
	}

	log.NilDebugf("%#v", data.Data)
	for key, value := range data.Data {
		viper.Set(key, value)
	}

	return nil
}

func certAuth(pki pki) error {
	if pki.Cert == "" || pki.Key == "" {
		log.Fatal("Certificate not defined for pki authentication")
	}

	cl := newClient(pki.CA, pki.Cert, pki.Key)
	auth, err := cl.Auth(methodPOST, fmt.Sprintf("%s/%s", pki.Addr, vaultPKIAuth))
	if err != nil {
		log.Fatalf("unable to authenticate vault: %s", err.Error())
	}

	log.Noticef("Authenticated with vault: %v token: %v****", pki.Addr, auth.Auth.ClientToken[:9])
	log.NilDebug(auth)
	vault.Token = auth.Auth.ClientToken

	return nil
}

func tokenRenewer(timeout time.Duration) {
	ticker60m := time.NewTicker(timeout)
	defer ticker60m.Stop()

	httpsrv.WaitShutdown.Add(1)
	defer httpsrv.WaitShutdown.Done()

	for {
		select {
		case <-ticker60m.C:
			log.Debug("Renewing Token ")
			cl := newClient(vault.CA, vault.PKI.Cert, vault.PKI.Key)
			cl.Token = vault.Token
			_, err := cl.Req(methodPOST, fmt.Sprintf("%s/%s", vault.Addr, vaultRenewSelf))
			if err != nil {
				log.Error(err)
				return
			}

		case <-httpsrv.SignalShutdown:
			log.Debug("Shutting Down Token Renew")
			return

		}
	}
}

func checkToken() error {
	cl := newClient(vault.CA, "", "")
	cl.Token = vault.Token
	auth, err := cl.Auth(methodGET, fmt.Sprintf("%s/%s", vault.Addr, vaultLookupSelf))
	if err != nil {
		return err
	}
	log.Debug("Token Policies", auth.Data.Policies)

	//if auth.Data.Renewable && auth.Data.TTL > 0 {
	//	log.Debug("Starting Token Renewer ", time.Second * time.Duration(auth.Data.TTL / 2))
	//	go tokenRenewer(time.Second * time.Duration(auth.Data.TTL / 2))
	//}

	return nil
}

type client struct {
	Client *http.Client
	Token  string
}

func newClient(ca, cert, key string) (c client) {
	cl := &http.Client{}

	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		caCertPool = x509.NewCertPool()
	}

	if ca != "" {
		// Load CA cert
		caCertPool.AppendCertsFromPEM([]byte(ca))
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	if cert != "" && key != "" {
		// Load client cert
		clientCert, err := tls.X509KeyPair([]byte(cert), []byte(key))
		if err != nil {
			log.Fatal(err)
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, clientCert)
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	cl.Transport = transport
	c.Client = cl

	return
}
func (c client) Req(method, url string) (data vaultData, err error) {

	var req *http.Request
	req, err = http.NewRequest(method, url, bytes.NewBufferString(""))
	if err != nil {
		return
	}
	req.Header.Set("content-type", "application/json")
	if c.Token != "" {
		req.Header.Set("x-vault-token", c.Token)
	}
	log.NilNotice("URL: ", url)
	res, err := c.Client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	log.NilInfo(method, url)
	log.NilInfo(res.Status)
	//res.Body = log.Tee(res.Body)
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("unable to read config")
		return
	}

	return
}
func (c client) Auth(method, url string) (auth vaultAuth, err error) {

	var req *http.Request
	req, err = http.NewRequest(method, url, bytes.NewBufferString(""))
	if err != nil {
		return
	}
	req.Header.Set("content-type", "application/json")
	if c.Token != "" {
		req.Header.Set("x-vault-token", c.Token)
	}
	log.NilNotice("URL: ", url)
	res, err := c.Client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	log.NilInfo(method, url)
	log.NilInfo(res.Status)
	//res.Body = log.Tee(res.Body)
	err = json.NewDecoder(res.Body).Decode(&auth)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("unable to read config: %v", auth.Errors)
		return
	}

	return
}
