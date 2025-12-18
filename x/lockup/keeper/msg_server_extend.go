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

		bondDenom, err := k.stakingKeeper.BondDenom(ctx)
		if err != nil {
			return nil, err
		}

		if extension.Amount.Denom != bondDenom {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid denom: %s, expected: %s", extension.Amount.Denom, bondDenom)
		}

		fromDate, err := time.Parse(time.DateOnly, extension.FromDate)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid from date format: %s", extension.FromDate)
		}

		toDate, err := time.Parse(time.DateOnly, extension.ToDate)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid to date format (%s)", extension.ToDate)
		}

		if !toDate.After(fromDate) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("to date must be after from date")
		}

		blockTime := ctx.BlockTime()
		blockDay := time.Date(blockTime.Year(), blockTime.Month(), blockTime.Day(), 0, 0, 0, 0, time.UTC)

		if blockDay.After(toDate) || blockDay.Equal(toDate) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("to date must be in the future")
		}

		if blockDay.AddDate(2, 0, 0).Before(toDate) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("to date cannot be more than 2 years from now")
		}

		existingLock, idx, found := lockupAcc.FindLock(extension.FromDate)
		if !found {
			return nil, types.ErrLockupNotFound.Wrapf("no lockup found for from date (%s)", extension.FromDate)
		}

		amountToMove := extension.Amount.Amount
		if existingLock.Amount.Amount.LT(amountToMove) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("extension amount exceeds existing lock amount date (%s)", extension.FromDate)
		} else if existingLock.Amount.Amount.Equal(amountToMove) {
			lockupAcc.Locks = lockupAcc.RemoveLock(idx)
		} else {
			updatedLock := &types.Lock{
				UnlockDate: existingLock.UnlockDate,
				Amount:     sdk.Coin{Denom: bondDenom, Amount: existingLock.Amount.Amount.Sub(amountToMove)},
			}
			lockupAcc.Locks = lockupAcc.UpdateLock(idx, updatedLock)
		}

		lockupAcc.Locks = lockupAcc.UpsertLock(extension.ToDate, sdk.Coin{Denom: bondDenom, Amount: amountToMove})

		if err := k.RemoveFromExpirationQueue(ctx, fromDate, addr, amountToMove); err != nil {
			return nil, err
		}

		if err := k.AddToExpirationQueue(ctx, toDate, addr, amountToMove); err != nil {
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
