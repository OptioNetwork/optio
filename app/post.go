package app

import (
	lockuppost "github.com/OptioNetwork/optio/x/lockup/post"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewPostHandler returns an empty PostHandler chain.
func NewPostHandler(options HandlerOptions) (sdk.PostHandler, error) {
	return sdk.ChainPostDecorators(
		lockuppost.NewRemoveExpiredLocksDecorator(options.AccountKeeper, *options.LockupKeeper),
	), nil

}
