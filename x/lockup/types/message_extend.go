package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgExtend{}

func NewMsgExtend(extendingAddress string, extensions []*Extension) *MsgExtend {
	return &MsgExtend{
		ExtendingAddress: extendingAddress,
		Extensions:       extensions,
	}
}

func (msg *MsgExtend) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.ExtendingAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid extendingAddress address (%s)", err)
	}
	return nil
}
