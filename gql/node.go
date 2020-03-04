package gql

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
)

func FmtID(format string, args ...interface{}) string {
	s := fmt.Sprintf(format, args...)

	return base64.RawURLEncoding.EncodeToString([]byte(s))
}

// SplitID decodes base64 if needed and splits ID
func SplitID(id string) ([]string, error) {
	if !strings.Contains(id, ":") {
		s, err := base64.RawURLEncoding.DecodeString(id)
		if err != nil {
			return nil, err
		}
		id = string(s)
	}

	return strings.Split(id, ":"), nil
}

// NodeFn handles a Node request for ID
type NodeFn func(ctx context.Context, id []string) (Node, error)

// GraphNode implements the node interface for Graphql
type GraphNode struct {
	fns map[string]NodeFn
}

// Node for use with graphql schema @goModel()
type Node interface {
	IsNode()
}

// Register a new handler
func (nh *GraphNode) Register(name string, fn NodeFn) {
	if nh.fns == nil {
		nh.fns = make(map[string]NodeFn)
	}

	nh.fns[name] = fn
}

// Node handles request from graphql
func (nh GraphNode) Node(ctx context.Context, id string) (Node, error) {
	ids, err := SplitID(id)
	if err != nil {
		return nil, err
	}
	if h, ok := nh.fns[ids[0]]; ok {
		return h(ctx, ids)
	}

	return nil, fmt.Errorf("Unsuppored Node Type: %s", ids[0])
}
