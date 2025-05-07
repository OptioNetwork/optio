package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgMint{}

func NewMsgMint(fromAddress string, amount uint64) *MsgMint {
	return &MsgMint{
		Amount: amount,
	}
}

func (msg *MsgMint) ValidateBasic() error {

	if msg.Amount == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "amount cannot be zero")
	}
	return nil
}
