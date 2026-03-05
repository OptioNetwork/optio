package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultParams returns a default set of parameters.
// The owner is intentionally empty and must be set before the module is used.
func DefaultParams() Params {
	return Params{
		Owner: "",
	}
}

// Validate validates the Params.
func (p Params) Validate() error {
	if p.Owner != "" {
		if _, err := sdk.AccAddressFromBech32(p.Owner); err != nil {
			return fmt.Errorf("invalid owner address: %w", err)
		}
	}
	return nil
}
