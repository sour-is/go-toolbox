package ident

import (
	"context"
	"sync"
)

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return "ident context key " + c.name
}

// ManagerKey defines the context key for finding the dataloader
var ManagerKey = contextKey{"ident.Manger"}

// Manager holds a mutex and ident
type Manager struct {
	mu    sync.RWMutex
	Ident Ident
}

// GetManager returns a new context
func GetManager(id Ident) *Manager {
	return &Manager{
		Ident: id,
	}
}

// WithContext wrap the context with a ident manager
func WithContext(ctx context.Context, id Ident) context.Context {
	return context.WithValue(ctx, ManagerKey, GetManager(id))
}

// GetIdent from the manager
func (m *Manager) GetIdent() *Ident {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return &m.Ident
}

// SetIdent into manager
func (m *Manager) SetIdent(u Ident) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Ident = u
}

// SetContextIdent writes ident into context
func SetContextIdent(ctx context.Context, s Ident) {
	ctx.Value(ManagerKey).(*Manager).SetIdent(s)
}

// GetContextIdent returns a user object from session
func GetContextIdent(ctx context.Context) (s Ident) {
	return ctx.Value(ManagerKey).(*Manager).Ident
}
