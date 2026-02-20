package assurance

// Level expresses confidence tiers without forcing KYC for all users.
type Level string

const (
	LevelSelfAsserted Level = "self_asserted"
	LevelVerified     Level = "verified"
	LevelGovernment   Level = "government_verified"
)

// Evidence represents a proof source for trust elevation.
type Evidence struct {
	Type      string `json:"type"`
	Reference string `json:"reference"`
	Issuer    string `json:"issuer"`
	Hash      string `json:"hash"`
}

// HasGovernmentEvidence indicates optional passport/KYC style trust elevation.
func HasGovernmentEvidence(items []Evidence) bool {
	for _, e := range items {
		if e.Type == "passport" || e.Type == "government_id" || e.Type == "kyc" {
			return true
		}
	}
	return false
}
