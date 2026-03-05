package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

func (k msgServer) CreateLicenseType(goCtx context.Context, req *types.MsgCreateLicenseType) (*types.MsgCreateLicenseTypeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if params.Owner != req.Owner {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "signer %s is not the module owner %s", req.Owner, params.Owner)
	}

	if req.Id == "" {
		return nil, errorsmod.Wrap(types.ErrLicenseTypeNotFound, "license type id cannot be empty")
	}

	_, exists, err := k.GetLicenseType(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errorsmod.Wrapf(types.ErrLicenseTypeExists, "license type %s already exists", req.Id)
	}

	lt := types.LicenseType{
		Id:           req.Id,
		Transferrable: req.Transferrable,
		MaxSupply:    req.MaxSupply,
		IssuedCount:  math.ZeroInt(),
	}

	if err := k.SetLicenseType(ctx, lt); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCreateLicenseType,
		sdk.NewAttribute(types.AttributeKeyLicenseTypeID, req.Id),
	))

	return &types.MsgCreateLicenseTypeResponse{}, nil
}
