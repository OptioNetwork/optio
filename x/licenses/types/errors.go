package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrInvalidSigner          = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrLicenseTypeNotFound    = sdkerrors.Register(ModuleName, 1101, "license type not found")
	ErrLicenseTypeExists      = sdkerrors.Register(ModuleName, 1102, "license type already exists")
	ErrMaxSupplyReached       = sdkerrors.Register(ModuleName, 1103, "license type max supply reached")
	ErrLicenseNotFound        = sdkerrors.Register(ModuleName, 1104, "license not found")
	ErrNotLicenseHolder       = sdkerrors.Register(ModuleName, 1105, "signer is not the license holder")
	ErrLicenseNotTransferable = sdkerrors.Register(ModuleName, 1106, "license type is not transferrable")
	ErrUnauthorized           = sdkerrors.Register(ModuleName, 1107, "signer does not have the required admin key permission")
)
