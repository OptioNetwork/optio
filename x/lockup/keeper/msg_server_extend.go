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

	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}

	events := sdk.Events{}
	for _, extension := range msg.Extensions {
		if extension.Lock.Coin.Denom != bondDenom {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid coin denomination: got %s, expected %s", extension.Lock.Coin.Denom, bondDenom)
		}

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

		fromLockup, idx, found := lockupAcc.FindLock(extension.From)
		if !found {
			return nil, types.ErrLockupNotFound.Wrapf("no lockup found for unlock date: %s", extension.From)
		}

		if !fromLockup.Coin.Amount.Equal(extension.Lock.Coin.Amount) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("extension amount mismatch for unlock from date: %s. you must extend the entire amount: %s", extension.From, fromLockup.Coin.Amount.String())
		}

		amountToMove := fromLockup.Coin.Amount
		lockupAcc.Locks = lockupAcc.RemoveLock(idx)
		lockupAcc.Locks = lockupAcc.UpsertLock(extension.Lock.UnlockDate, sdk.NewCoin(bondDenom, amountToMove))

		events = events.AppendEvent(sdk.NewEvent(
			types.EventTypeLockExtended,
			sdk.NewAttribute(types.AttributeKeyLockAddress, msg.ExtendingAddress),
			sdk.NewAttribute(types.AttributeKeyOldUnlockDate, extension.From),
			sdk.NewAttribute(types.AttributeKeyUnlockDate, extension.Lock.UnlockDate),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amountToMove.String()),
		))
	}

	k.accountKeeper.SetAccount(ctx, lockupAcc)

	ctx.EventManager().EmitEvents(events)

	return &types.MsgExtendResponse{}, nil
}
