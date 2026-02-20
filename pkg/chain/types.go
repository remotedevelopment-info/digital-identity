package chain

import "time"

// EventType identifies the category of identity event.
type EventType string

const (
	EventIdentityAssertion EventType = "identity_assertion"
	EventLogin             EventType = "login"
	EventVerification      EventType = "verification"
	EventAuthorization     EventType = "authorization"
	EventRecovery          EventType = "recovery"
)

// IdentityEvent is the user-owned auditable event payload.
type IdentityEvent struct {
	ID        string            `json:"id"`
	Type      EventType         `json:"type"`
	Timestamp time.Time         `json:"timestamp"`
	ActorID   string            `json:"actor_id"`
	Risk      string            `json:"risk"`
	Payload   map[string]string `json:"payload"`
}

// BlockLink is a cryptographically linked append-only record.
type BlockLink struct {
	Index           uint64 `json:"index"`
	PrevHash        string `json:"prev_hash"`
	EventHash       string `json:"event_hash"`
	Hash            string `json:"hash"`
	Signature       string `json:"signature"`
	SignerPublicKey string `json:"signer_public_key"`
}

// IdentityChain is a sovereign per-user event chain.
type IdentityChain struct {
	OwnerID       string      `json:"owner_id"`
	RootPublicKey string      `json:"root_public_key"`
	CreatedAt     time.Time   `json:"created_at"`
	Blocks        []BlockLink `json:"blocks"`
}
