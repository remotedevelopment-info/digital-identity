package chain

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const genesisPrevHash = "GENESIS"

// NewIdentityChain initializes an empty chain for one owner.
func NewIdentityChain(ownerID string, rootPublicKey ed25519.PublicKey) *IdentityChain {
	return &IdentityChain{
		OwnerID:       ownerID,
		RootPublicKey: base64.StdEncoding.EncodeToString(rootPublicKey),
		CreatedAt:     time.Now().UTC(),
		Blocks:        make([]BlockLink, 0),
	}
}

// AppendEvent appends an event as the next immutable block.
func AppendEvent(c *IdentityChain, e IdentityEvent, signerPrivateKey ed25519.PrivateKey) error {
	if c == nil {
		return fmt.Errorf("chain is nil")
	}
	if len(signerPrivateKey) == 0 {
		return fmt.Errorf("signer private key is required")
	}
	signerPub := base64.StdEncoding.EncodeToString(signerPrivateKey.Public().(ed25519.PublicKey))
	if signerPub != c.RootPublicKey {
		return fmt.Errorf("signer is not the root owner key")
	}
	if e.Type == "" {
		return fmt.Errorf("event type is required")
	}
	if e.ActorID == "" {
		return fmt.Errorf("actor_id is required")
	}
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	if e.ID == "" {
		e.ID = newEventID(e)
	}

	eh, err := hashEvent(e)
	if err != nil {
		return err
	}

	prevHash := genesisPrevHash
	if n := len(c.Blocks); n > 0 {
		prevHash = c.Blocks[n-1].Hash
	}

	next := BlockLink{
		Index:     uint64(len(c.Blocks)),
		PrevHash:  prevHash,
		EventHash: eh,
	}
	next.Hash = hashLink(next)

	hashBytes, err := hex.DecodeString(next.Hash)
	if err != nil {
		return fmt.Errorf("decode link hash: %w", err)
	}
	sig := ed25519.Sign(signerPrivateKey, hashBytes)
	next.Signature = base64.StdEncoding.EncodeToString(sig)
	next.SignerPublicKey = signerPub

	c.Blocks = append(c.Blocks, next)
	return nil
}

// Verify checks cryptographic integrity and signatures across the full chain.
func Verify(c *IdentityChain) error {
	if c == nil {
		return fmt.Errorf("chain is nil")
	}

	for i, b := range c.Blocks {
		expectedPrev := genesisPrevHash
		if i > 0 {
			expectedPrev = c.Blocks[i-1].Hash
		}
		if b.PrevHash != expectedPrev {
			return fmt.Errorf("block %d prev hash mismatch", i)
		}

		rebuilt := hashLink(BlockLink{Index: b.Index, PrevHash: b.PrevHash, EventHash: b.EventHash})
		if b.Hash != rebuilt {
			return fmt.Errorf("block %d hash mismatch", i)
		}

		pub, err := base64.StdEncoding.DecodeString(b.SignerPublicKey)
		if err != nil {
			return fmt.Errorf("block %d signer key decode error: %w", i, err)
		}
		sig, err := base64.StdEncoding.DecodeString(b.Signature)
		if err != nil {
			return fmt.Errorf("block %d signature decode error: %w", i, err)
		}
		hashBytes, err := hex.DecodeString(b.Hash)
		if err != nil {
			return fmt.Errorf("block %d hash decode error: %w", i, err)
		}
		if !ed25519.Verify(ed25519.PublicKey(pub), hashBytes, sig) {
			return fmt.Errorf("block %d signature verification failed", i)
		}
	}
	return nil
}

func newEventID(e IdentityEvent) string {
	raw := fmt.Sprintf("%s|%s|%s|%d", e.Type, e.ActorID, e.Risk, time.Now().UnixNano())
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func hashEvent(e IdentityEvent) (string, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("marshal event: %w", err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

func hashLink(b BlockLink) string {
	raw := fmt.Sprintf("%d|%s|%s", b.Index, b.PrevHash, b.EventHash)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
