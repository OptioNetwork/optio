package antehandler

import (
	txsigning "cosmossdk.io/x/tx/signing"

	lockupante "github.com/OptioNetwork/optio/x/lockup/ante"
	"github.com/OptioNetwork/optio/x/lockup/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type HandlerOptions struct {
	AccountKeeper          ante.AccountKeeper
	BankKeeper             bankkeeper.Keeper
	ExtensionOptionChecker ante.ExtensionOptionChecker
	FeegrantKeeper         ante.FeegrantKeeper
	SignModeHandler        *txsigning.HandlerMap
	SigGasConsumer         ante.SignatureVerificationGasConsumer
	TxFeeChecker           ante.TxFeeChecker
	LockupKeeper           keeper.Keeper
	StakingKeeper          stakingkeeper.Keeper
}

func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		lockupante.NewLockedDelegationsDecorator(options.AccountKeeper, options.BankKeeper, options.LockupKeeper, options.StakingKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
	), nil
}
