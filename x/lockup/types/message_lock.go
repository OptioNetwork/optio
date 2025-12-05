package types

import (
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgLock{}

func NewMsgLock(lockupAddress string, locks []*Lock) *MsgLock {
	return &MsgLock{
		LockupAddress: lockupAddress,
		Locks:         locks,
	}
}

func (msg *MsgLock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.LockupAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid lockupAddress address (%s)", err)
	}

	if len(msg.Locks) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "locks cannot be empty")
	}

	for i, lock := range msg.Locks {
		if lock == nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "lock at index %d is nil", i)
		}

		if lock.UnlockDate == "" {
			return errorsmod.Wrapf(ErrInvalidDate, "lock at index %d has empty unlock date", i)
		}
		_, err := time.Parse(time.DateOnly, lock.UnlockDate)
		if err != nil {
			return errorsmod.Wrapf(ErrInvalidDate, "lock at index %d has invalid unlock date format: %s", i, err)
		}

		if !lock.Coin.IsValid() {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "lock at index %d has invalid coin: %s", i, lock.Coin.String())
		}
		if !lock.Coin.IsPositive() {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "lock at index %d has non-positive coin amount: %s", i, lock.Coin.String())
		}
	}

	return nil
}
