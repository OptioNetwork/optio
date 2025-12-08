package keeper

import (
	"context"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	"github.com/OptioNetwork/optio/x/lockup/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k msgServer) Lock(goCtx context.Context, msg *types.MsgLock) (*types.MsgLockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lockupAddr, err := sdk.AccAddressFromBech32(msg.LockupAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid lockup address: %s", err)
	}
	acc := k.accountKeeper.GetAccount(ctx, lockupAddr)

	if acc == nil {
		baseAcc := authtypes.NewBaseAccountWithAddress(lockupAddr)
		lAcc := types.NewLockupAccount(baseAcc)
		k.accountKeeper.SetAccount(ctx, lAcc)
		acc = lAcc
	} else {
		acc = k.accountKeeper.GetAccount(ctx, lockupAddr)
		_, ok := acc.(*types.Account)
		if !ok {
			lAcc := types.NewLockupAccount(acc.(*authtypes.BaseAccount))
			k.accountKeeper.SetAccount(ctx, lAcc)
			acc = lAcc
		}
	}

	lockupAcc := acc.(*types.Account)
	amountToLock := math.ZeroInt()

	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}

	for _, lock := range msg.Locks {
		if !lock.Coin.IsValid() || lock.Coin.IsZero() {
			return nil, sdkerrors.ErrInvalidCoins.Wrapf("invalid lock amount: %s", lock.Coin.String())
		}

		if lock.Coin.Denom != bondDenom {
			return nil, sdkerrors.ErrInvalidCoins.Wrapf("invalid denom: %s, expected: %s", lock.Coin.Denom, bondDenom)
		}

		unlockTime, err := time.Parse(time.DateOnly, lock.UnlockDate)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid unlock date format: %s", lock.UnlockDate)
		}

		if ctx.BlockTime().AddDate(0, 6, 0).After(unlockTime) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("unlock time must be at least 6 months from now")
		}
		if ctx.BlockTime().AddDate(2, 0, 0).Before(unlockTime) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("unlock time cannot be more than 2 years from now")
		}
		amountToLock = amountToLock.Add(lock.Coin.Amount)
	}

	currentLockedAmount := lockupAcc.GetLockedAmount(ctx.BlockTime())
	totalDelegatedAmount, err := k.GetTotalDelegatedAmount(ctx, lockupAddr)
	if err != nil {
		return nil, err
	}

	if totalDelegatedAmount.LT(currentLockedAmount.Add(amountToLock)) {
		return nil, errorsmod.Wrapf(
			types.ErrInsufficientDelegations,
			"insufficient delegated tokens to create new locks by the requested amount: %s < %s",
			totalDelegatedAmount.String(),
			currentLockedAmount.Add(amountToLock).String(),
		)
	}

	for _, lock := range msg.Locks {
		lockupAcc.Locks = lockupAcc.UpsertLock(lock.UnlockDate, lock.Coin)

		unlockTime, err := time.Parse(time.DateOnly, lock.UnlockDate)
		if err != nil {
			return nil, err
		}

		if err := k.AddToExpirationQueue(ctx, unlockTime, lockupAddr, lock.Coin.Amount); err != nil {
			return nil, err
		}
	}

	currentTotal, err := k.GetTotalLocked(ctx)
	if err != nil {
		return nil, err
	}
	newTotal := currentTotal.Add(amountToLock)
	if err := k.SetTotalLocked(ctx, newTotal); err != nil {
		return nil, err
	}

	k.accountKeeper.SetAccount(ctx, lockupAcc)

	return &types.MsgLockResponse{}, nil
}
