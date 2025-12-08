package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSendDelegateAndLock{}

func NewMsgSendDelegateAndLock(fromAddress string, toAddress string, valAddress string, amount string, unlockDate string) *MsgSendDelegateAndLock {
	a, ok := math.NewIntFromString(amount)
	if !ok {
		panic("invalid amount string")
	}

	return &MsgSendDelegateAndLock{
		FromAddress:      fromAddress,
		ToAddress:        toAddress,
		ValidatorAddress: valAddress,
		Lock: &Lock{
			UnlockDate: unlockDate,
			Coin:       sdk.Coin{Denom: "uOPT", Amount: a},
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
