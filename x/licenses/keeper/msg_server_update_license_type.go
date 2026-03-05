package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

func (k msgServer) UpdateLicenseType(goCtx context.Context, req *types.MsgUpdateLicenseType) (*types.MsgUpdateLicenseTypeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if params.Owner != req.Owner {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "signer %s is not the module owner %s", req.Owner, params.Owner)
	}

	lt, found, err := k.GetLicenseType(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrLicenseTypeNotFound, "license type %s not found", req.Id)
	}

	if !req.MaxSupply.IsZero() && lt.IssuedCount.GT(req.MaxSupply) {
		return nil, errorsmod.Wrapf(types.ErrMaxSupplyReached, "cannot set max_supply to %s: %s licenses already issued", req.MaxSupply.String(), lt.IssuedCount.String())
	}

	lt.Transferrable = req.Transferrable
	lt.MaxSupply = req.MaxSupply

	if err := k.SetLicenseType(ctx, lt); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUpdateLicenseType,
		sdk.NewAttribute(types.AttributeKeyLicenseTypeID, req.Id),
	))

	return &types.MsgUpdateLicenseTypeResponse{}, nil
}
