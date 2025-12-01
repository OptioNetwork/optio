package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/lockup module sentinel errors
var (
	ErrInvalidSigner  = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidAccount = sdkerrors.Register(ModuleName, 1101, "invalid account type for lockup")
)
