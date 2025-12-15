package posthandler

import (
	lockupkeeper "github.com/OptioNetwork/optio/x/lockup/keeper"
	lockuppost "github.com/OptioNetwork/optio/x/lockup/post"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
)

// HandlerOptions are the options required for constructing a default SDK PostHandler.
type HandlerOptions struct {
	AccountKeeper authkeeper.AccountKeeper
	LockupKeeper  lockupkeeper.Keeper
}

// NewPostHandler returns an empty PostHandler chain.
func NewPostHandler(options HandlerOptions) (sdk.PostHandler, error) {
	return sdk.ChainPostDecorators(
		lockuppost.NewRemoveExpiredLocksDecorator(options.AccountKeeper, options.LockupKeeper),
	), nil

}
