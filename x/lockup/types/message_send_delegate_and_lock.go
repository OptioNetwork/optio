package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSendDelegateAndLock{}

func NewMsgSendDelegateAndLock(fromAddress string, toAddress string, ValidatorAddress string, amount math.Int, unlockDate string) *MsgSendDelegateAndLock {
	return &MsgSendDelegateAndLock{
		FromAddress:      fromAddress,
		ToAddress:        toAddress,
		ValidatorAddress: ValidatorAddress,
		Lock: &Lock{
			UnlockDate: unlockDate,
			Amount:     amount,
		},
	}
}

func (msg *MsgSendDelegateAndLock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid fromAddress address (%s)", err)
	}
	return nil
}
