package ante

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/OptioNetwork/optio/x/lockup/keeper"
	"github.com/OptioNetwork/optio/x/lockup/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type LockedDelegationsDecorator struct {
	accountKeeper ante.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	lockupKeeper  keeper.Keeper
}

func NewLockedDelegationsDecorator(accountKeeper ante.AccountKeeper, bankKeeper bankkeeper.Keeper, lockupKeeper keeper.Keeper) LockedDelegationsDecorator {
	return LockedDelegationsDecorator{
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		lockupKeeper:  lockupKeeper,
	}
}

func (d LockedDelegationsDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	msgs := tx.GetMsgs()

	err = handleMsgs(ctx, msgs, d)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func handleMsgs(ctx sdk.Context, msgs []sdk.Msg, lbd LockedDelegationsDecorator) error {
	for _, msg := range msgs {
		switch msg.(type) {
		case *stakingtypes.MsgUndelegate:

			err := handleMsgUndelegate(ctx, msg, lbd)
			if err != nil {
				return err
			}

		case *authztypes.MsgExec:

			err := handleMsgExec(ctx, msg, lbd)
			if err != nil {
				return err
			}

		default:
			continue
		}
	}

	return nil
}

// handleMsgUndelegate checks if the undelegation would cause the delegator's total delegated amount to be less than their locked amount.
func handleMsgUndelegate(ctx sdk.Context, msg sdk.Msg, lbd LockedDelegationsDecorator) error {
	msgUndelegate := msg.(*stakingtypes.MsgUndelegate)
	fromAddr, err := sdk.AccAddressFromBech32(msgUndelegate.DelegatorAddress)
	if err != nil {
		return err
	}

	acc := lbd.accountKeeper.GetAccount(ctx, fromAddr)
	if ltsAcc, ok := acc.(*types.Account); ok {
		totalLocked := math.ZeroInt()
		blockTime := ctx.BlockTime()
		for _, lock := range ltsAcc.Locks {
			if types.IsLocked(blockTime, lock.UnlockDate) {
				totalLocked = totalLocked.Add(lock.Amount.Amount)
			}
		}

		totalDelegated, err := lbd.lockupKeeper.GetTotalDelegatedAmount(ctx, fromAddr)
		if err != nil {
			return err
		}
		delegatedAfterUndelegate := totalDelegated.Sub(msgUndelegate.Amount.Amount)

		if delegatedAfterUndelegate.LT(totalLocked) {
			return errorsmod.Wrapf(
				types.ErrInsufficientDelegations,
				"unbond would cause new delegated amount to be less than the locked amount: %s < %s",
				delegatedAfterUndelegate.String(),
				totalLocked.String(),
			)
		}
	}

	return nil
}

// handleMsgExec checks each inner message of MsgExec for locked balance constraints.
func handleMsgExec(ctx sdk.Context, msg sdk.Msg, lbd LockedDelegationsDecorator) error {
	msgExec := msg.(*authztypes.MsgExec)
	msgs, err := msgExec.GetMessages()
	if err != nil {
		return err
	}

	err = handleMsgs(ctx, msgs, lbd)
	if err != nil {
		return err
	}

	return nil
}
