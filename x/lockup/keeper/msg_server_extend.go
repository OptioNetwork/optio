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

	events := sdk.Events{}
	for _, extension := range msg.Extensions {

		from, err := time.Parse(time.DateOnly, extension.FromDate)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid from format: %s", extension.FromDate)
		}

		unlock, err := time.Parse(time.DateOnly, extension.ToDate)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid unlock date format: %s", extension.ToDate)
		}

		if !unlock.After(from) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("unlock date must be after from date")
		}

		existingLock, idx, found := lockupAcc.FindLock(extension.FromDate)
		if !found {
			return nil, types.ErrLockupNotFound.Wrapf("no lockup found for unlock date: %s", extension.FromDate)
		}

		if !existingLock.Amount.Equal(extension.Amount) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("extension amount mismatch for unlock from date: %s. you must extend the entire amount: %s", extension.FromDate, existingLock.Amount.String())
		}

		amountToMove := extension.Amount
		if existingLock.Amount.LT(amountToMove) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("extension amount exceeds existing lock amount for unlock from date: %s", extension.FromDate)
		}

		existingLock.Amount = existingLock.Amount.Sub(amountToMove)

		lockupAcc.Locks = lockupAcc.UpdateLock(idx, existingLock)
		lockupAcc.Locks = lockupAcc.UpsertLock(extension.ToDate, amountToMove)

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
			sdk.NewAttribute(types.AttributeKeyUnlockDate, extension.ToDate),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amountToMove.String()),
		))
	}

	k.accountKeeper.SetAccount(ctx, lockupAcc)

	ctx.EventManager().EmitEvents(events)

	return &types.MsgExtendResponse{}, nil
}
