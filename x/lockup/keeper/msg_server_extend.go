package keeper

import (
	"context"
	"time"

	"github.com/OptioNetwork/optio/x/lockup/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) Extend(goCtx context.Context, msg *types.MsgExtend) (*types.MsgExtendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.ExtendingAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid lockup address: %s", err)
	}

	acc := k.accountKeeper.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("no account found for address: %s", msg.ExtendingAddress)
	}

	lockupAcc, ok := acc.(*types.Account)
	if !ok {
		return nil, types.ErrInvalidAccount.Wrapf("account is not a long-term stake account: %s", msg.ExtendingAddress)
	}

	newLockups := map[string]*types.Lock{}
	for _, extension := range msg.Extensions {
		from, err := time.Parse(time.DateOnly, extension.From)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid from format: %s", extension.From)
		}

		to, err := time.Parse(time.DateOnly, extension.Lock.UnlockDate)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid unlock date format: %s", extension.Lock.UnlockDate)
		}

		if !to.After(from) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("unlock date must be after from date")
		}

		existingLock, exists := newLockups[extension.Lock.UnlockDate]
		if exists {
			existingLock.Coin.Amount = existingLock.Coin.Amount.Add(extension.Lock.Coin.Amount)
			newLockups[extension.Lock.UnlockDate] = existingLock
			continue
		}
		newLockups[extension.Lock.UnlockDate] = extension.Lock
	}

	for _, lockup := range newLockups {
		existingLockup, exists := lockupAcc.Lockups[lockup.UnlockDate]
		if exists {
			existingLockup.Coin.Amount = existingLockup.Coin.Amount.Add(lockup.Coin.Amount)
			lockupAcc.Lockups[lockup.UnlockDate] = existingLockup
		} else {
			lockupAcc.Lockups[lockup.UnlockDate] = &types.Lockup{
				Coin: lockup.Coin,
			}
		}
	}
	k.accountKeeper.SetAccount(ctx, lockupAcc)

	return &types.MsgExtendResponse{}, nil
}
