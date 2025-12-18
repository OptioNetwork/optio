package keeper

import (
	"context"
	"time"

	errorsmod "cosmossdk.io/errors"

	"github.com/OptioNetwork/optio/x/lockup/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k msgServer) Lock(goCtx context.Context, msg *types.MsgLock) (*types.MsgLockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid lockup address: %s", err)
	}
	acc := k.accountKeeper.GetAccount(ctx, address)

	if acc == nil {
		baseAcc := authtypes.NewBaseAccountWithAddress(address)
		lAcc := types.NewLockupAccount(baseAcc)
		k.accountKeeper.SetAccount(ctx, lAcc)
		acc = lAcc
	} else {
		_, ok := acc.(*types.Account)
		if !ok {
			lAcc := types.NewLockupAccount(acc.(*authtypes.BaseAccount))
			k.accountKeeper.SetAccount(ctx, lAcc)
			acc = lAcc
		}
	}

	lockupAcc := acc.(*types.Account)

	if !msg.Amount.IsPositive() || msg.Amount.IsZero() {
		return nil, sdkerrors.ErrInvalidCoins.Wrapf("invalid lock amount: %s", msg.Amount.String())
	}

	unlockDate, err := time.Parse(time.DateOnly, msg.UnlockDate)
	if err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid unlock date format: %s", msg.UnlockDate)
	}

	blockTime := ctx.BlockTime()

	if blockTime.After(unlockDate) {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("unlock date must be in the future")
	}

	if blockTime.AddDate(2, 0, 0).Before(unlockDate) {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("unlock date cannot be more than 2 years from now")
	}

	currentLockedAmount := lockupAcc.GetLockedAmount(ctx.BlockTime())
	totalDelegatedAmount, err := k.GetTotalDelegatedAmount(ctx, address)
	if err != nil {
		return nil, err
	}

	if totalDelegatedAmount.LT(currentLockedAmount.Add(msg.Amount)) {
		return nil, errorsmod.Wrapf(
			types.ErrInsufficientDelegations,
			"insufficient delegated tokens to create new locks by the requested amount: %s < %s",
			totalDelegatedAmount.String(),
			currentLockedAmount.Add(msg.Amount).String(),
		)
	}

	lockupAcc.Locks = lockupAcc.UpsertLock(msg.UnlockDate, msg.Amount)

	if err := k.AddToExpirationQueue(ctx, unlockDate, address, msg.Amount); err != nil {
		return nil, err
	}

	k.accountKeeper.SetAccount(ctx, lockupAcc)

	ctx.EventManager().EmitEvents([]sdk.Event{
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyLockAddress, msg.Address),
			sdk.NewAttribute(types.AttributeKeyUnlockDate, msg.UnlockDate),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	})

	return &types.MsgLockResponse{}, nil
}
