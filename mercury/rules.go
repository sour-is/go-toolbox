package mercury

import (
	"path/filepath"
)

// Rule is a type of rule
type Rule struct {
	Role  string
	Type  string
	Match string
}

// Rules is a list of rules
type Rules []Rule

// GetNamespaceSearch returns a default search for users rules.
func (r Rules) GetNamespaceSearch() (lis NamespaceSearch) {
	for _, o := range r {
		if o.Type == "NS" && (o.Role == "read" || o.Role == "write") {
			lis = append(lis, NamespaceStar(o.Match))
		}
	}
	return
}

// Check if name matches rule
func (r Rule) Check(name string) bool {
	ok, err := filepath.Match(r.Match, name)
	if err != nil {
		return false
	}
	return ok
}

// CheckNamespace verifies user has access
func (r Rules) CheckNamespace(search NamespaceSearch) bool {
	for _, ns := range search {
		if !r.GetRoles("NS", ns.Value()).HasRole("read", "write") {
			return false
		}
	}

	return true
}

// ReduceSearch verifies user has access
func (r Rules) ReduceSearch(search NamespaceSearch) (out NamespaceSearch) {
	rules := r.GetNamespaceSearch()
	skip := make(map[string]struct{})
	out = make(NamespaceSearch, 0, len(rules))

	for _, rule := range rules {
		if _, ok := skip[rule.Raw()]; ok {
			continue
		}
		for _, ck := range search {
			if _, ok := skip[ck.Raw()]; ok {
				continue
			} else if rule.Match(ck.Raw()) {
				skip[ck.Raw()] = struct{}{}
				out = append(out, ck)
			} else if ck.Match(rule.Raw()) {
				out = append(out, rule)
			}
		}
	}

	return
}

// Roles is a list of roles for a resource
type Roles map[string]struct{}

// GetRoles returns a list of Roles
func (r Rules) GetRoles(typ, name string) (lis Roles) {
	for _, o := range r {
		if typ == o.Type && o.Check(name) {
			lis[o.Role] = struct{}{}
		}
	}
	return
}

// HasRole is a valid role
func (r Roles) HasRole(roles ...string) bool {
	for _, role := range roles {
		if _, ok := r[role]; ok {
			return true
		}
	}
	return false
}
