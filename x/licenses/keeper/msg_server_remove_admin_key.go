package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

func (k msgServer) RemoveAdminKey(goCtx context.Context, req *types.MsgRemoveAdminKey) (*types.MsgRemoveAdminKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if params.Owner != req.Owner {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "signer %s is not the module owner %s", req.Owner, params.Owner)
	}

	if err := k.DeleteAdminKey(ctx, req.Address); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRemoveAdminKey,
		sdk.NewAttribute(types.AttributeKeyAddress, req.Address),
	))

	return &types.MsgRemoveAdminKeyResponse{}, nil
}
