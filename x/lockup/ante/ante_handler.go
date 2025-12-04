package ante

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/OptioNetwork/optio/x/lockup/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// NewLockedBalanceDecorator checks that the sender has sufficient unlocked balance for msgs that spend funds.
type LockedBalanceDecorator struct {
	accountKeeper ante.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	stakingKeeper stakingkeeper.Keeper
}

func NewLockedBalanceDecorator(accountKeeper ante.AccountKeeper, bankKeeper bankkeeper.Keeper) LockedBalanceDecorator {
	return LockedBalanceDecorator{
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

func (lbd LockedBalanceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	msgs := tx.GetMsgs()

	err = handleMsgs(ctx, msgs, lbd)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func handleMsgs(ctx sdk.Context, msgs []sdk.Msg, lbd LockedBalanceDecorator) error {
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
			return nil
		}
	}

	return nil
}

// handleMsgUndelegate checks if the undelegation would cause the delegator's total delegated amount to be less than their locked amount.
func handleMsgUndelegate(ctx sdk.Context, msg sdk.Msg, lbd LockedBalanceDecorator) error {
	msgUndelegate := msg.(*stakingtypes.MsgUndelegate)
	fromAddr, err := sdk.AccAddressFromBech32(msgUndelegate.DelegatorAddress)
	if err != nil {
		return err
	}

	acc := lbd.accountKeeper.GetAccount(ctx, fromAddr)
	if ltsAcc, ok := acc.(*types.Account); ok {
		totalLocked := math.NewInt(0)
		blockTime := ctx.BlockTime()
		for unlockDate, lockup := range ltsAcc.Lockups {
			if types.IsLocked(blockTime, unlockDate) {
				totalLocked = totalLocked.Add(lockup.Coin.Amount)
			}
		}

		delegations, err := lbd.stakingKeeper.GetDelegatorDelegations(ctx, fromAddr, 1000)
		if err != nil {
			return err
		}

		totalDelegated := math.NewInt(0)
		for _, delegation := range delegations {
			totalDelegated = totalDelegated.Add(delegation.Shares.TruncateInt())
		}

		delegatedAfterUndelegate := totalDelegated.Sub(msgUndelegate.Amount.Amount)

		if delegatedAfterUndelegate.LT(totalLocked) {
			return errorsmod.Wrapf(
				types.ErrInsufficientDelegations,
				"undelegation would cause new delegated amount to be less than the locked amount: %s < %s",
				delegatedAfterUndelegate.String(),
				totalLocked.String(),
			)
		}
	}

	return nil
}

// handleMsgExec checks each inner message of MsgExec for locked balance constraints.
func handleMsgExec(ctx sdk.Context, msg sdk.Msg, lbd LockedBalanceDecorator) error {
	msgGrant := msg.(*authztypes.MsgExec)
	var sdkMsgs []sdk.Msg
	for _, innerMsg := range msgGrant.Msgs {
		var sdkMsg sdk.Msg
		registry := codectypes.NewInterfaceRegistry()
		err := registry.UnpackAny(innerMsg, &sdkMsg)
		if err != nil {
			return err
		}
		sdkMsgs = append(sdkMsgs, sdkMsg)
	}

	err := handleMsgs(ctx, sdkMsgs, lbd)
	if err != nil {
		return err
	}

	return nil
}
