package post

import (
	"time"

	"github.com/OptioNetwork/optio/x/lockup/keeper"
	"github.com/OptioNetwork/optio/x/lockup/types"

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

	blockTime := ctx.BlockTime()

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return next(ctx, tx, simulate, success)
	}

	feePayer := feeTx.FeePayer()
	if feePayer == nil {
		return next(ctx, tx, simulate, success)
	}

	acc := d.accountKeeper.GetAccount(ctx, feePayer)
	if acc == nil {
		return next(ctx, tx, simulate, success)
	}

	lockupAcc, ok := acc.(*types.Account)
	if !ok {
		return next(ctx, tx, simulate, success)
	}

	updated := false
	i := 0
	for i < len(lockupAcc.Locks) {
		lock := lockupAcc.Locks[i]

		unlockTime, err := time.Parse(time.DateOnly, lock.UnlockDate)
		if err != nil {
			i++
			continue
		}

		blockTimeDateOnly, err := time.Parse(time.DateOnly, blockTime.Format(time.DateOnly))
		if err != nil {
			i++
			continue
		}

		if unlockTime.Before(blockTimeDateOnly) {
			lockupAcc.Locks = lockupAcc.RemoveLock(i)
			updated = true
		} else {
			i++
		}
	}

	if updated {
		d.accountKeeper.SetAccount(ctx, lockupAcc)
	}

	return next(ctx, tx, simulate, success)
}
