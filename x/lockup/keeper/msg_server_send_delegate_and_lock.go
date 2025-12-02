package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/OptioNetwork/optio/x/lockup/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k msgServer) SendDelegateAndLock(goCtx context.Context, msg *types.MsgSendDelegateAndLock) (*types.MsgSendDelegateAndLockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}

	amount, ok := math.NewIntFromString(msg.Amount)
	if !ok {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid amount: %s", msg.Amount)
	}

	lock := &types.Lock{
		Coin:       sdk.NewCoin(bondDenom, amount),
		UnlockDate: msg.UnlockDate,
	}

	if !lock.Coin.IsValid() || !lock.Coin.IsPositive() {
		return nil, sdkerrors.ErrInvalidCoins.Wrapf("invalid lock amount: %s", lock.Coin.String())
	}

	if lock.Coin.Denom != bondDenom {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid coin denomination: got %s, expected %s", lock.Coin.Denom, bondDenom)
	}

	fromAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	toAddr, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.ValAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoins(ctx, fromAddr, toAddr, sdk.NewCoins(lock.Coin))
	if err != nil {
		return nil, err
	}

	newShares, err := k.stakingKeeper.Delegate(ctx, toAddr, lock.Coin.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValAddress),
			sdk.NewAttribute(stakingtypes.AttributeKeyDelegator, msg.ToAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, lock.Coin.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
		),
	})

	// now lock the funds (skip balance validation since we just delegated)
	lockupMsg := &types.MsgLock{
		LockupAddress: msg.ToAddress,
		Lockups:       []*types.Lock{lock},
	}

	_, err = k.Lock(goCtx, lockupMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendDelegateAndLockResponse{}, nil
}
