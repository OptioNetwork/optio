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

	for _, extension := range msg.Extensions {
		from, err := time.Parse(time.DateOnly, extension.From)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid from format: %s", extension.From)
		}

		unlock, err := time.Parse(time.DateOnly, extension.Lock.UnlockDate)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid unlock date format: %s", extension.Lock.UnlockDate)
		}

		if !unlock.After(from) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("unlock date must be after from date")
		}

		fromLockup, exists := lockupAcc.Lockups[extension.From]
		if !exists {
			return nil, types.ErrLockupNotFound.Wrapf("no lockup found for unlock date: %s", extension.From)
		}

		if !fromLockup.Coin.Amount.Equal(extension.Lock.Coin.Amount) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("lockup amount mismatch for unlock date: %s. you must extend the entire amount", extension.From)
		}

		delete(lockupAcc.Lockups, extension.From)

		newLockup, exists := lockupAcc.Lockups[extension.Lock.UnlockDate]
		if exists {
			newLockup.Coin.Amount = newLockup.Coin.Amount.Add(fromLockup.Coin.Amount)
		} else {
			newLockup = &types.Lockup{
				Coin: fromLockup.Coin,
			}
		}
		lockupAcc.Lockups[extension.Lock.UnlockDate] = newLockup
	}

	k.accountKeeper.SetAccount(ctx, lockupAcc)

	return &types.MsgExtendResponse{}, nil
}
