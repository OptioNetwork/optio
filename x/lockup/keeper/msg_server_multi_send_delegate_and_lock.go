package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/OptioNetwork/optio/x/lockup/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) MultiSendDelegateAndLock(goCtx context.Context, msg *types.MsgMultiSendDelegateAndLock) (*types.MsgMultiSendDelegateAndLockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	totalOutputs := math.ZeroInt()
	for _, output := range msg.Outputs {
		if !output.Amount.IsPositive() || output.Amount.IsZero() {
			return nil, sdkerrors.ErrInvalidCoins.Wrapf("invalid coin in output to %s: %s", output.ToAddress, err)
		}
		totalOutputs = totalOutputs.Add(output.Amount)
	}
	if !msg.Input.Equal(totalOutputs) {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("total amount %s does not match sum of outputs %s", msg.Input.String(), totalOutputs.String())
	}

	for _, output := range msg.Outputs {
		msg := &types.MsgSendDelegateAndLock{
			FromAddress:      msg.FromAddress,
			ToAddress:        output.ToAddress,
			ValidatorAddress: output.ValidatorAddress,
			Lock: &types.Lock{
				UnlockDate: output.UnlockDate,
				Amount:     output.Amount,
			},
		}
		_, err := k.SendDelegateAndLock(ctx, msg)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgMultiSendDelegateAndLockResponse{}, nil
}
