package anchor

import "context"

// Checkpoint is chain-agnostic proof data for optional anchoring.
type Checkpoint struct {
	ChainID   string `json:"chain_id"`
	Height    uint64 `json:"height"`
	RootHash  string `json:"root_hash"`
	Reference string `json:"reference"`
}

// Adapter allows anchoring into any blockchain without core protocol lock-in.
type Adapter interface {
	Name() string
	Anchor(ctx context.Context, rootHash string) (Checkpoint, error)
	Verify(ctx context.Context, checkpoint Checkpoint) error
}
