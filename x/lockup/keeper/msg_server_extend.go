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

	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid lockup address: %s", err)
	}

	acc := k.accountKeeper.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("no account found for address: %s", msg.Address)
	}

	lockupAcc, ok := acc.(*types.Account)
	if !ok {
		return nil, types.ErrInvalidAccount.Wrapf("account is not a long-term stake account: %s", msg.Address)
	}

	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}

	events := sdk.Events{}
	for _, extension := range msg.Extensions {
		if extension.Lock.Amount.Denom != bondDenom {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid coin denomination: got %s, expected %s", extension.Lock.Amount.Denom, bondDenom)
		}

		from, err := time.Parse(time.DateOnly, extension.FromDate)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid from format: %s", extension.FromDate)
		}

		unlock, err := time.Parse(time.DateOnly, extension.Lock.UnlockDate)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid unlock date format: %s", extension.Lock.UnlockDate)
		}

		if !unlock.After(from) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("unlock date must be after from date")
		}

		fromLockup, idx, found := lockupAcc.FindLock(extension.FromDate)
		if !found {
			return nil, types.ErrLockupNotFound.Wrapf("no lockup found for unlock date: %s", extension.FromDate)
		}

		if !fromLockup.Amount.Amount.Equal(extension.Lock.Amount.Amount) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("extension amount mismatch for unlock from date: %s. you must extend the entire amount: %s", extension.FromDate, fromLockup.Amount.Amount.String())
		}

		amountToMove := fromLockup.Amount.Amount
		lockupAcc.Locks = lockupAcc.RemoveLock(idx)
		lockupAcc.Locks = lockupAcc.UpsertLock(extension.Lock.UnlockDate, sdk.NewCoin(bondDenom, amountToMove))

		// Update expiration queue: remove from old date, add to new date
		if err := k.RemoveFromExpirationQueue(ctx, from, addr, amountToMove); err != nil {
			return nil, err
		}

		if err := k.AddToExpirationQueue(ctx, unlock, addr, amountToMove); err != nil {
			return nil, err
		}

		events = events.AppendEvent(sdk.NewEvent(
			types.EventTypeLockExtended,
			sdk.NewAttribute(types.AttributeKeyLockAddress, msg.Address),
			sdk.NewAttribute(types.AttributeKeyOldUnlockDate, extension.FromDate),
			sdk.NewAttribute(types.AttributeKeyUnlockDate, extension.Lock.UnlockDate),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amountToMove.String()),
		))
	}

	k.accountKeeper.SetAccount(ctx, lockupAcc)

	ctx.EventManager().EmitEvents(events)

	return &types.MsgExtendResponse{}, nil
}
