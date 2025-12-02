package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSendDelegateAndLock{}

func NewMsgSendDelegateAndLock(fromAddress string, toAddress string, valAddress string, amount string, unlockDate string) *MsgSendDelegateAndLock {
	return &MsgSendDelegateAndLock{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		ValAddress:  valAddress,
		Amount:      amount,
		UnlockDate:  unlockDate,
	}
}

func (msg *MsgSendDelegateAndLock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid fromAddress address (%s)", err)
	}
	return nil
}
