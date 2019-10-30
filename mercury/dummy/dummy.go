package dummy

import (
	"sour.is/x/toolbox/dbm/rsql"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/mercury"
)

// IndexDummy implements Handler.Index
type IndexDummy struct{}

// GetIndex returns nil
func (IndexDummy) GetIndex(mercury.NamespaceSearch, *rsql.Program) mercury.Config { return nil }

// ObjectsDummy implements Handler.Objects
type ObjectsDummy struct{}

// GetObjects returns nil
func (ObjectsDummy) GetObjects(mercury.NamespaceSearch, *rsql.Program, []string) mercury.Config {
	return nil
}

// WriteDummy implements Handler.Write
type WriteDummy struct{}

// WriteObjects returns nil
func (WriteDummy) WriteObjects(mercury.Config) error { return nil }

// RulesDummy implements Handler.Rules
type RulesDummy struct{}

// GetRules returns nil
func (RulesDummy) GetRules(ident.Ident) mercury.Rules { return nil }

// GroupsDummy implements Handler.Groups
type GroupsDummy struct{}

// GetGroups returns nil
func (GroupsDummy) GetGroups(ident.Ident) mercury.Rules { return nil }

// NotifyDummy implements Handler.Groups
type NotifyDummy struct{}

// GetNotify returns nil
func (NotifyDummy) GetNotify(string) mercury.ListNotify { return nil }
