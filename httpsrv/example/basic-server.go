package main

import (
	"bytes"

	"github.com/spf13/viper"
	"sour.is/x/toolbox/httpsrv"
	_ "sour.is/x/toolbox/httpsrv/routes"
	_ "sour.is/x/toolbox/ident/mock"
	"sour.is/x/toolbox/log"
)

var defaultConfig = `
[http]
listen   = ":8060"
idm = "mock"

[idm.mock]
ident = "user"
aspect = "default"
name = "User Name"
`

func init() {
	log.SetVerbose(log.Vdebug)

	viper.SetConfigType("toml")
	viper.ReadConfig(bytes.NewBuffer([]byte(defaultConfig)))

	httpsrv.Config()

}

func main() {
	httpsrv.Run()
}
