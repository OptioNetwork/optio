package post

import (
	"time"

	"github.com/OptioNetwork/optio/x/lockup/keeper"

	lockuptypes "github.com/OptioNetwork/optio/x/lockup/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

type RemoveExpiredLocksDecorator struct {
	accountKeeper ante.AccountKeeper
	lockupKeeper  keeper.Keeper
}

func NewRemoveExpiredLocksDecorator(accountKeeper ante.AccountKeeper, lockupKeeper keeper.Keeper) RemoveExpiredLocksDecorator {
	return RemoveExpiredLocksDecorator{
		accountKeeper: accountKeeper,
		lockupKeeper:  lockupKeeper,
	}
}

func (d RemoveExpiredLocksDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate, success bool, next sdk.PostHandler) (newCtx sdk.Context, err error) {
	if simulate {
		return next(ctx, tx, simulate, success)
	}

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return next(ctx, tx, simulate, success)
	}

	feePayer := feeTx.FeePayer()
	if feePayer == nil {
		return next(ctx, tx, simulate, success)
	}

	addr := sdk.AccAddress(feePayer)
	locks, err := d.lockupKeeper.GetLocksByAddress(ctx, addr)
	if err != nil {
		return ctx, err
	}

	blockTime := ctx.BlockTime()
	blockDate := time.Date(blockTime.Year(), blockTime.Month(), blockTime.Day(), 0, 0, 0, 0, time.UTC)

	updated := false
	// Iterate backwards to safely delete while iterating
	for i := len(locks) - 1; i >= 0; i-- {
		lock := locks[i]
		unlockTime, err := time.Parse(time.DateOnly, lock.UnlockDate)
		if err != nil {
			continue
		}

		if !lockuptypes.IsLocked(blockDate, lock.UnlockDate) {
			if err := d.lockupKeeper.RemoveFromExpirationQueue(ctx, unlockTime, addr, lock.Amount); err != nil {
				return ctx, err
			}

			locks = append(locks[:i], locks[i+1:]...)
			updated = true
		}
	}

	if updated {
		if err := d.lockupKeeper.SetLocksByAddress(ctx, addr, locks); err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate, success)
}
