package keeper

import (
	"context"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlocker(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockTime := sdkCtx.BlockTime()

	globalTotal, err := k.GetTotalLocked(ctx)
	if err != nil {
		return err
	}

	newTotal := globalTotal
	save := false

	err = k.IterateAndDeleteExpiredLocks(ctx, blockTime, func(addr sdk.AccAddress, unlockTime time.Time, amount math.Int) error {
		save = true

		newTotal = newTotal.Sub(amount)

		return nil
	})
	if err != nil {
		return err
	}

	if save {
		if err := k.SetTotalLocked(ctx, newTotal); err != nil {
			return err
		}
	}

	return nil
}
