package chain

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
)

func TestAppendAndVerify(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("key gen: %v", err)
	}

	c := NewIdentityChain("user:test", priv.Public().(ed25519.PublicKey))
	e := IdentityEvent{
		Type:    EventLogin,
		ActorID: "user:test",
		Risk:    "normal",
		Payload: map[string]string{"ip": "127.0.0.1"},
	}

	if err := AppendEvent(c, e, priv); err != nil {
		t.Fatalf("append: %v", err)
	}
	if err := Verify(c); err != nil {
		t.Fatalf("verify: %v", err)
	}
}
