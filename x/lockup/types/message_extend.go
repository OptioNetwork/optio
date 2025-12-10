package types

import (
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgExtend{}

func NewMsgExtend(extendingAddress string, extensions []*Extension) *MsgExtend {
	return &MsgExtend{
		Address:    extendingAddress,
		Extensions: extensions,
	}
}

func (msg *MsgExtend) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid extendingAddress address (%s)", err)
	}

	if len(msg.Extensions) == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "at least one extension is required")
	}

	for i, extension := range msg.Extensions {
		if extension == nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "extension cannot be nil")
		}

		if extension.FromDate == "" {
			return errorsmod.Wrapf(ErrInvalidDate, "extension from date cannot be empty")
		}

		fromTime, err := time.Parse(time.DateOnly, extension.FromDate)
		if err != nil {
			return errorsmod.Wrapf(ErrInvalidDate, "extension at index %d has invalid 'from' date format: %s", i, err)
		}

		if extension.Lock == nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "extension lock cannot be nil")
		}

		if !extension.Lock.Amount.IsValid() || extension.Lock.Amount.IsZero() {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "invalid lock amount: %s", extension.Lock.Amount.String())
		}

		if extension.Lock.Amount.Denom == "" {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "lock coin denom cannot be empty")
		}

		if extension.Lock.Amount.Denom != "uOPT" {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "lock coin denom must be 'uOPT'")
		}

		if extension.Lock.UnlockDate == "" {
			return errorsmod.Wrapf(ErrInvalidDate, "extension lock unlock date cannot be empty")
		}

		toTime, err := time.Parse(time.DateOnly, extension.Lock.UnlockDate)
		if err != nil {
			return errorsmod.Wrapf(ErrInvalidDate, "invalid unlock date format: %s", extension.Lock.UnlockDate)
		}

		if !toTime.After(fromTime) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "unlock date must be after from date")
		}
	}

	return nil
}
