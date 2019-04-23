package app

import (
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/mercury"
	"sour.is/x/toolbox/mercury/dummy"
)

func init() {
	mercury.Register("*", 0, appDefault{})
}

type appDefault struct {
	dummy.IndexDummy
	dummy.ObjectsDummy
	dummy.WriteDummy
	dummy.NotifyDummy
}

// GetRules returns default rules for user role.
func (appDefault) GetRules(u ident.Ident) (lis mercury.Rules) {

	if u.HasRole("admin") {
		lis = append(lis,
			mercury.Rule{
				Role:  "admin",
				Type:  "NS",
				Match: "*",
			},
			mercury.Rule{
				Role:  "write",
				Type:  "NS",
				Match: "*",
			},
			mercury.Rule{
				Role:  "admin",
				Type:  "GR",
				Match: "*",
			},
		)
	} else if u.HasRole("write") {
		lis = append(lis,
			mercury.Rule{
				Role:  "write",
				Type:  "NS",
				Match: "*",
			},
		)
	} else if u.HasRole("read") {
		lis = append(lis,
			mercury.Rule{
				Role:  "read",
				Type:  "NS",
				Match: "*",
			},
		)
	}

	return lis
}
