package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

func (k msgServer) RevokeLicense(goCtx context.Context, req *types.MsgRevokeLicense) (*types.MsgRevokeLicenseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	license, found, err := k.GetLicense(ctx, req.LicenseTypeId, req.Id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrLicenseNotFound, "license (type=%s, id=%d) not found", req.LicenseTypeId, req.Id)
	}

	if !k.HasPermission(ctx, req.Revoker, "revoke", license.Type) {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "%s does not have revoke permission for license type %s", req.Revoker, license.Type)
	}

	if err := k.DeleteLicense(ctx, license); err != nil {
		return nil, err
	}

	lt, found, err := k.GetLicenseType(ctx, license.Type)
	if err != nil {
		return nil, err
	}
	if found && !lt.IssuedCount.IsZero() {
		lt.IssuedCount = lt.IssuedCount.SubRaw(1)
		if err := k.SetLicenseType(ctx, lt); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRevokeLicense,
		sdk.NewAttribute(types.AttributeKeyLicenseTypeID, req.LicenseTypeId),
		sdk.NewAttribute(types.AttributeKeyLicenseID, fmt.Sprintf("%d", req.Id)),
	))

	return &types.MsgRevokeLicenseResponse{}, nil
}
