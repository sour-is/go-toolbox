package gql

import (
	"context"

	"github.com/spf13/viper"
)

// GraphGql Implements the AppInfo method for query
type GraphGql struct{}

// AppInfo returns app information
func (GraphGql) AppInfo(ctx context.Context) (s *AppInfo, err error) {
	return getAppInfo()
}

// AppInfo about running server
type AppInfo struct {
	// Application Name
	Name string `json:"name"`
	// Version number
	Version string `json:"version"`
	// Build information
	Build string `json:"build"`
}

func getAppInfo() (o *AppInfo, err error) {
	app := viper.GetStringMapString("app")
	o = new(AppInfo)
	o.Name = app["name"]
	o.Version = app["version"]
	o.Build = app["build"]
	return
}
