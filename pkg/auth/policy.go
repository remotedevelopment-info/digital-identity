package auth

import "errors"

// RiskLevel controls whether 2FA or 3FA is required.
type RiskLevel string

const (
	RiskNormal RiskLevel = "normal"
	RiskHigh   RiskLevel = "high"
)

// AuthContext carries factors required for root actions.
type AuthContext struct {
	LongPhrase bool
	EmailOTP   bool
	TOTP       bool
	Hardware   bool
	Risk       RiskLevel
}

// ValidateRootAction enforces long phrase + 2FA/3FA policy.
func ValidateRootAction(ctx AuthContext) error {
	if !ctx.LongPhrase {
		return errors.New("long phrase is required")
	}

	factorCount := 0
	if ctx.EmailOTP {
		factorCount++
	}
	if ctx.TOTP {
		factorCount++
	}
	if ctx.Hardware {
		factorCount++
	}

	required := 2
	if ctx.Risk == RiskHigh {
		required = 3
	}
	if factorCount < required {
		return errors.New("insufficient secondary factors")
	}

	return nil
}
