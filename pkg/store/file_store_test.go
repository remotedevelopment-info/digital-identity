package store

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"path/filepath"
	"testing"

	"github.com/nicholasalexander/digital-identity/pkg/chain"
)

func TestCreateGetUpdate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chains.json")

	s, err := NewFileStore(path)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("key gen: %v", err)
	}

	c := chain.NewIdentityChain("user:one", pub)
	if err := s.Create(context.Background(), c); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := s.Get(context.Background(), "user:one")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.OwnerID != "user:one" {
		t.Fatalf("unexpected owner: %s", got.OwnerID)
	}

	got.Blocks = append(got.Blocks, chain.BlockLink{Index: 0, PrevHash: "GENESIS"})
	if err := s.Update(context.Background(), got); err != nil {
		t.Fatalf("update: %v", err)
	}

	all, err := s.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 chain, got %d", len(all))
	}
}
