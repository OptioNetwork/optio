package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgLock{}

func NewMsgLock(lockupAddress string, lockups []*Lock) *MsgLock {
	return &MsgLock{
		LockupAddress: lockupAddress,
		Lockups:       lockups,
	}
}

func (msg *MsgLock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.LockupAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid lockupAddress address (%s)", err)
	}
	return nil
}
