package keeper

import (
	"context"
	"time"

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
		// create a new lockup account if it doesn't exist
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
	totalLockedAmount := math.ZeroInt()
	for unlockDate, lockup := range lockupAcc.Lockups {
		if types.IsLocked(ctx.BlockTime(), unlockDate) {
			totalLockedAmount = totalLockedAmount.Add(lockup.Coin.Amount)
		}
	}

	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}
	for _, lock := range msg.Lockups {
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
		totalLockedAmount = totalLockedAmount.Add(lock.Coin.Amount)
	}

	// Validate that total locked amount doesn't exceed total delegations
	totalDelegations := math.ZeroInt()
	delegations, err := k.stakingKeeper.GetDelegatorDelegations(ctx, lockupAddr, 1000)
	if err != nil {
		return nil, err
	}

	for _, delegation := range delegations {
		valAddr, err := sdk.ValAddressFromBech32(delegation.GetValidatorAddr())
		if err != nil {
			return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
		}
		validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
		if err != nil {
			return nil, err
		}
		tokens := validator.TokensFromShares(delegation.GetShares())
		totalDelegations = totalDelegations.Add(tokens.TruncateInt())
	}

	if totalDelegations.LT(totalLockedAmount) {
		return nil, sdkerrors.ErrInsufficientFunds.Wrapf("trying to lock more than total delegations: total delegations %s < locked amount %s", totalDelegations.String(), totalLockedAmount.String())
	}

	for _, msgLockup := range msg.Lockups {
		existingLockup, exists := lockupAcc.Lockups[msgLockup.UnlockDate]
		if exists {
			existingLockup.Coin.Amount = existingLockup.Coin.Amount.Add(msgLockup.Coin.Amount)
			lockupAcc.Lockups[msgLockup.UnlockDate] = existingLockup
		} else {
			lockupAcc.Lockups[msgLockup.UnlockDate] = &types.Lockup{
				Coin: msgLockup.Coin,
			}
		}
	}

	k.accountKeeper.SetAccount(ctx, lockupAcc)

	return &types.MsgLockResponse{}, nil
}
