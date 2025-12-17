package ante

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/OptioNetwork/optio/x/lockup/keeper"
	"github.com/OptioNetwork/optio/x/lockup/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"cosmossdk.io/x/feegrant"

	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

type LockedDelegationsDecorator struct {
	accountKeeper ante.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	lockupKeeper  keeper.Keeper
	stakingKeeper stakingkeeper.Keeper
}

func NewLockedDelegationsDecorator(accountKeeper ante.AccountKeeper, bankKeeper bankkeeper.Keeper, lockupKeeper keeper.Keeper, stakingKeeper stakingkeeper.Keeper) LockedDelegationsDecorator {
	return LockedDelegationsDecorator{
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		lockupKeeper:  lockupKeeper,
		stakingKeeper: stakingKeeper,
	}
}

func (d LockedDelegationsDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	msgs := tx.GetMsgs()

	bondDenom, err := d.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return ctx, err
	}

	err = d.handleMsgs(ctx, msgs, bondDenom)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func (d LockedDelegationsDecorator) handleMsgs(ctx sdk.Context, msgs []sdk.Msg, bondDenom string) error {
	for _, msg := range msgs {
		switch msg.(type) {
		case *stakingtypes.MsgUndelegate:

			err := d.handleMsgUndelegate(ctx, msg, bondDenom)
			if err != nil {
				return err
			}

		case *authztypes.MsgExec:

			err := d.handleMsgExec(ctx, msg, bondDenom)
			if err != nil {
				return err
			}

		case *banktypes.MsgSend:

			err := d.handleMsgSend(ctx, msg, bondDenom)
			if err != nil {
				return err
			}

		case *banktypes.MsgMultiSend:

			err := d.handleMsgMultiSend(ctx, msg, bondDenom)
			if err != nil {
				return err
			}

		case *authztypes.MsgGrant: // Is this needed? This doesn't actually move funds, just grants permission.

			err := d.handleMsgGrant(ctx, msg, bondDenom)
			if err != nil {
				return err
			}

		case *govtypes.MsgDeposit:

			err := d.handleMsgDeposit(ctx, msg, bondDenom)
			if err != nil {
				return err
			}

		case *distributiontypes.MsgFundCommunityPool:

			err := d.handleMsgFundCommunityPool(ctx, msg, bondDenom)
			if err != nil {
				return err
			}

		case *feegrant.MsgGrantAllowance:

			err := d.handleMsgGrantAllowance(ctx, msg, bondDenom)
			if err != nil {
				return err
			}

		case *ibctransfertypes.MsgTransfer:

			err := d.handleMsgTransfer(ctx, msg, bondDenom)
			if err != nil {
				return err
			}

		default:
			continue
		}
	}

	return nil
}

func (d LockedDelegationsDecorator) handleMsgSend(ctx sdk.Context, msg sdk.Msg, bondDenom string) error {
	sendMsg := msg.(*banktypes.MsgSend)
	fromAddr, err := sdk.AccAddressFromBech32(sendMsg.FromAddress)
	if err != nil {
		return err
	}

	acc := d.accountKeeper.GetAccount(ctx, fromAddr)
	if ltsAcc, ok := acc.(*types.Account); ok {
		ok, lockedAboveDelegated, err := checkDelegationsAgainstLocked(ctx, ltsAcc, d.lockupKeeper)
		if err != nil {
			return err
		}

		if ok {
			return nil
		}

		for _, coin := range sendMsg.Amount {
			if coin.Denom != bondDenom {
				continue
			}

			totalBalance := d.bankKeeper.GetBalance(ctx, fromAddr, coin.Denom).Amount
			available := totalBalance.Sub(*lockedAboveDelegated)
			if available.LT(coin.Amount) {
				return errorsmod.Wrapf(errortypes.ErrInsufficientFunds, "insufficient unlocked balance: available %s, required %s", available, coin.Amount)
			}
		}
	}

	return nil
}

func (d LockedDelegationsDecorator) handleMsgMultiSend(ctx sdk.Context, msg sdk.Msg, bondDenom string) error {
	multiSendMsg := msg.(*banktypes.MsgMultiSend)
	// Check each input
	for _, input := range multiSendMsg.Inputs {
		fromAddr, err := sdk.AccAddressFromBech32(input.Address)
		if err != nil {
			return err
		}
		acc := d.accountKeeper.GetAccount(ctx, fromAddr)
		if ltsAcc, ok := acc.(*types.Account); ok {
			ok, lockedAboveDelegated, err := checkDelegationsAgainstLocked(ctx, ltsAcc, d.lockupKeeper)
			if err != nil {
				return err
			}

			if ok {
				continue
			}

			for _, coin := range input.Coins {
				if coin.Denom != bondDenom {
					continue
				}

				totalBalance := d.bankKeeper.GetBalance(ctx, fromAddr, coin.Denom).Amount
				available := totalBalance.Sub(*lockedAboveDelegated)
				if available.LT(coin.Amount) {
					return errorsmod.Wrapf(errortypes.ErrInsufficientFunds, "insufficient unlocked balance: available %s, required %s", available, coin.Amount)
				}
			}
		}
	}

	return nil
}

func (d LockedDelegationsDecorator) handleMsgGrant(ctx sdk.Context, msg sdk.Msg, bondDenom string) error {
	msgGrant := msg.(*authztypes.MsgGrant)

	fromAddr, err := sdk.AccAddressFromBech32(msgGrant.Granter)
	if err != nil {
		return err
	}
	acc := d.accountKeeper.GetAccount(ctx, fromAddr)
	if ltsAcc, ok := acc.(*types.Account); ok {

		ok, lockedAboveDelegated, err := checkDelegationsAgainstLocked(ctx, ltsAcc, d.lockupKeeper)
		if err != nil {
			return err
		}

		if ok {
			return nil
		}

		if msgGrant.Grant.Authorization.GetTypeUrl() == sdk.MsgTypeURL(&banktypes.SendAuthorization{}) {
			authorization, err := msgGrant.Grant.GetAuthorization()
			if err != nil {
				return err
			}
			sendAuth, ok := authorization.(*banktypes.SendAuthorization)
			if ok {
				for _, coin := range sendAuth.SpendLimit {
					if coin.Denom != bondDenom {
						continue
					}

					totalBalance := d.bankKeeper.GetBalance(ctx, fromAddr, coin.Denom).Amount
					available := totalBalance.Sub(*lockedAboveDelegated)
					if available.LT(coin.Amount) {
						return errorsmod.Wrapf(errortypes.ErrInsufficientFunds, "insufficient unlocked balance: available %s, required %s", available, coin.Amount)
					}
				}
			}
		}
	}

	return nil
}

// handleMsgUndelegate checks if the undelegation would cause the delegator's total delegated amount to be less than their locked amount.
func (d LockedDelegationsDecorator) handleMsgUndelegate(ctx sdk.Context, msg sdk.Msg, bondDenom string) error {
	msgUndelegate := msg.(*stakingtypes.MsgUndelegate)
	if msgUndelegate.Amount.Denom != bondDenom {
		return nil
	}

	fromAddr, err := sdk.AccAddressFromBech32(msgUndelegate.DelegatorAddress)
	if err != nil {
		return err
	}

	acc := d.accountKeeper.GetAccount(ctx, fromAddr)
	if ltsAcc, ok := acc.(*types.Account); ok {
		totalLocked := math.ZeroInt()
		blockTime := ctx.BlockTime()
		for _, lock := range ltsAcc.Locks {
			if types.IsLocked(blockTime, lock.UnlockDate) {
				totalLocked = totalLocked.Add(lock.Amount.Amount)
			}
		}

		totalDelegated, err := d.lockupKeeper.GetTotalDelegatedAmount(ctx, fromAddr)
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
func (d LockedDelegationsDecorator) handleMsgExec(ctx sdk.Context, msg sdk.Msg, bondDenom string) error {
	msgExec := msg.(*authztypes.MsgExec)
	msgs, err := msgExec.GetMessages()
	if err != nil {
		return err
	}

	err = d.handleMsgs(ctx, msgs, bondDenom)
	if err != nil {
		return err
	}

	return nil
}

// handleMsgDeposit checks if the deposit is sending locked funds.
func (d LockedDelegationsDecorator) handleMsgDeposit(ctx sdk.Context, msg sdk.Msg, bondDenom string) error {
	msgDeposit := msg.(*govtypes.MsgDeposit)
	fromAddr, err := sdk.AccAddressFromBech32(msgDeposit.Depositor)
	if err != nil {
		return err
	}
	acc := d.accountKeeper.GetAccount(ctx, fromAddr)

	if ltsAcc, ok := acc.(*types.Account); ok {

		ok, lockedAboveDelegated, err := checkDelegationsAgainstLocked(ctx, ltsAcc, d.lockupKeeper)
		if err != nil {
			return err
		}

		if ok {
			return nil
		}

		for _, coin := range msgDeposit.Amount {
			if coin.Denom != bondDenom {
				continue
			}

			totalBalance := d.bankKeeper.GetBalance(ctx, fromAddr, coin.Denom).Amount
			available := totalBalance.Sub(*lockedAboveDelegated)
			if available.LT(coin.Amount) {
				return errorsmod.Wrapf(errortypes.ErrInsufficientFunds, "insufficient unlocked balance: available %s, required %s", available, coin.Amount)
			}
		}
	}

	return nil
}

// handleMsgFundCommunityPool checks if the funder has sufficient unlocked balance to fund the community pool.
func (d LockedDelegationsDecorator) handleMsgFundCommunityPool(ctx sdk.Context, msg sdk.Msg, bondDenom string) error {
	msgFCP := msg.(*distributiontypes.MsgFundCommunityPool)
	fromAddr, err := sdk.AccAddressFromBech32(msgFCP.Depositor)
	if err != nil {
		return err
	}
	acc := d.accountKeeper.GetAccount(ctx, fromAddr)

	if ltsAcc, ok := acc.(*types.Account); ok {

		ok, lockedAboveDelegated, err := checkDelegationsAgainstLocked(ctx, ltsAcc, d.lockupKeeper)
		if err != nil {
			return err
		}

		if ok {
			return nil
		}

		for _, coin := range msgFCP.Amount {
			if coin.Denom != bondDenom {
				continue
			}

			totalBalance := d.bankKeeper.GetBalance(ctx, fromAddr, coin.Denom).Amount
			available := totalBalance.Sub(*lockedAboveDelegated)
			if available.LT(coin.Amount) {
				return errorsmod.Wrapf(errortypes.ErrInsufficientFunds, "insufficient unlocked balance: available %s, required %s", available, coin.Amount)
			}
		}
	}

	return nil
}

// handleMsgGrantAllowance checks if the granter has sufficient unlocked balance to grant the allowance.
func (d LockedDelegationsDecorator) handleMsgGrantAllowance(ctx sdk.Context, msg sdk.Msg, bondDenom string) error {
	msgFeeGrant := msg.(*feegrant.MsgGrantAllowance)
	fromAddr, err := sdk.AccAddressFromBech32(msgFeeGrant.Granter)
	if err != nil {
		return err
	}
	acc := d.accountKeeper.GetAccount(ctx, fromAddr)

	if ltsAcc, ok := acc.(*types.Account); ok {

		ok, lockedAboveDelegated, err := checkDelegationsAgainstLocked(ctx, ltsAcc, d.lockupKeeper)
		if err != nil {
			return err
		}

		if ok {
			return nil
		}

		var coins []sdk.Coin
		allowance, ok := msgFeeGrant.Allowance.GetCachedValue().(feegrant.FeeAllowanceI)
		if !ok {
			return nil
		}

		switch a := allowance.(type) {
		case *feegrant.BasicAllowance:
			if a.SpendLimit == nil {
				return nil
			}
			coins = a.SpendLimit
		case *feegrant.PeriodicAllowance:
			if a.Basic.SpendLimit == nil {
				return nil
			}
			coins = a.Basic.SpendLimit
		case *feegrant.AllowedMsgAllowance:
			allowanceInner, ok := a.Allowance.GetCachedValue().(feegrant.FeeAllowanceI)
			if !ok {
				return nil
			}

			switch ai := allowanceInner.(type) {
			case *feegrant.BasicAllowance:
				if ai.SpendLimit == nil {
					return nil
				}
				coins = ai.SpendLimit
			case *feegrant.PeriodicAllowance:
				if ai.Basic.SpendLimit == nil {
					return nil
				}
				coins = ai.Basic.SpendLimit
			default:
				return nil
			}
		default:
			return nil
		}

		for _, coin := range coins {
			if coin.Denom != bondDenom {
				continue
			}

			totalBalance := d.bankKeeper.GetBalance(ctx, fromAddr, coin.Denom).Amount
			available := totalBalance.Sub(*lockedAboveDelegated)
			if available.LT(coin.Amount) {
				return errorsmod.Wrapf(errortypes.ErrInsufficientFunds, "insufficient unlocked balance: available %s, required %s", available, coin.Amount)
			}
		}
	}

	return nil
}

// handleMsgTransfer checks if the transfer is sending locked funds.
func (d LockedDelegationsDecorator) handleMsgTransfer(ctx sdk.Context, msg sdk.Msg, bondDenom string) error {
	msgTransfer := msg.(*ibctransfertypes.MsgTransfer)
	fromAddr, err := sdk.AccAddressFromBech32(msgTransfer.Sender)
	if err != nil {
		return err
	}
	acc := d.accountKeeper.GetAccount(ctx, fromAddr)

	if ltsAcc, ok := acc.(*types.Account); ok {

		ok, lockedAboveDelegated, err := checkDelegationsAgainstLocked(ctx, ltsAcc, d.lockupKeeper)
		if err != nil {
			return err
		}

		if ok {
			return nil
		}

		if msgTransfer.Token.Denom != bondDenom {
			return nil
		}

		totalBalance := d.bankKeeper.GetBalance(ctx, fromAddr, msgTransfer.Token.Denom).Amount
		available := totalBalance.Sub(*lockedAboveDelegated)
		if available.LT(msgTransfer.Token.Amount) {
			return errorsmod.Wrapf(errortypes.ErrInsufficientFunds, "insufficient unlocked balance: available %s, required %s", available, msgTransfer.Token.Amount)
		}
	}

	return nil
}

// checkDelegationsAgainstLocked checks if the total delegated amount is greater than the total locked amount.
// Returns true if total delegated amount is greater than total locked amount.
func checkDelegationsAgainstLocked(ctx sdk.Context, ltsAcc *types.Account, lockupKeeper keeper.Keeper) (bool, *math.Int, error) {

	delegationsTotal, err := lockupKeeper.GetTotalDelegatedAmount(ctx, ltsAcc.GetAddress())
	if err != nil {
		return false, nil, err
	}

	totalLocked := math.NewInt(0)
	for _, lock := range ltsAcc.Locks {
		if types.IsLocked(ctx.BlockTime(), lock.UnlockDate) {
			totalLocked = totalLocked.Add(lock.Amount.Amount)
		}
	}

	lockedAboveDelegated := totalLocked.Sub(*delegationsTotal)
	if lockedAboveDelegated.IsNegative() {
		lockedAboveDelegated = math.ZeroInt()
	}

	return delegationsTotal.GTE(totalLocked), &lockedAboveDelegated, nil

}
