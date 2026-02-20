package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/nicholasalexander/digital-identity/pkg/chain"
)

var ErrNotFound = errors.New("chain not found")

// ChainStore persists identity chains.
type ChainStore interface {
	Create(ctx context.Context, c *chain.IdentityChain) error
	Get(ctx context.Context, ownerID string) (*chain.IdentityChain, error)
	Update(ctx context.Context, c *chain.IdentityChain) error
	List(ctx context.Context) ([]*chain.IdentityChain, error)
}

type fileModel struct {
	Chains map[string]*chain.IdentityChain `json:"chains"`
}

// FileStore stores chains in a local JSON file.
type FileStore struct {
	mu     sync.RWMutex
	path   string
	chains map[string]*chain.IdentityChain
}

func NewFileStore(path string) (*FileStore, error) {
	fs := &FileStore{path: path, chains: map[string]*chain.IdentityChain{}}
	if err := fs.load(); err != nil {
		return nil, err
	}
	return fs, nil
}

func (s *FileStore) Create(_ context.Context, c *chain.IdentityChain) error {
	if c == nil {
		return errors.New("chain is nil")
	}
	if c.OwnerID == "" {
		return errors.New("owner_id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.chains[c.OwnerID]; ok {
		return fmt.Errorf("chain already exists for owner %s", c.OwnerID)
	}
	s.chains[c.OwnerID] = clone(c)
	return s.saveLocked()
}

func (s *FileStore) Get(_ context.Context, ownerID string) (*chain.IdentityChain, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.chains[ownerID]
	if !ok {
		return nil, ErrNotFound
	}
	return clone(c), nil
}

func (s *FileStore) Update(_ context.Context, c *chain.IdentityChain) error {
	if c == nil || c.OwnerID == "" {
		return errors.New("valid chain is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.chains[c.OwnerID]; !ok {
		return ErrNotFound
	}
	s.chains[c.OwnerID] = clone(c)
	return s.saveLocked()
}

func (s *FileStore) List(_ context.Context) ([]*chain.IdentityChain, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*chain.IdentityChain, 0, len(s.chains))
	for _, c := range s.chains {
		out = append(out, clone(c))
	}
	return out, nil
}

func (s *FileStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(s.path); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	b, err := os.ReadFile(s.path)
	if err != nil {
		return fmt.Errorf("read store: %w", err)
	}
	if len(b) == 0 {
		return nil
	}

	var m fileModel
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("unmarshal store: %w", err)
	}
	if m.Chains != nil {
		s.chains = m.Chains
	}
	return nil
}

func (s *FileStore) saveLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("mkdir store dir: %w", err)
	}

	b, err := json.MarshalIndent(fileModel{Chains: s.chains}, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal store: %w", err)
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return fmt.Errorf("write temp store: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("commit store: %w", err)
	}
	return nil
}

func clone(c *chain.IdentityChain) *chain.IdentityChain {
	if c == nil {
		return nil
	}
	dup := *c
	dup.Blocks = append([]chain.BlockLink(nil), c.Blocks...)
	return &dup
}
