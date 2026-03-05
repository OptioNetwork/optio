package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

func (k msgServer) UpdateLicense(goCtx context.Context, req *types.MsgUpdateLicense) (*types.MsgUpdateLicenseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	license, found, err := k.GetLicense(ctx, req.LicenseTypeId, req.Id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrLicenseNotFound, "license (type=%s, id=%d) not found", req.LicenseTypeId, req.Id)
	}

	if !k.HasPermission(ctx, req.Updater, "update", license.Type) {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "%s does not have update permission for license type %s", req.Updater, license.Type)
	}

	if req.Status != "active" && req.Status != "revoked" {
		return nil, fmt.Errorf("invalid status %q: must be \"active\" or \"revoked\"", req.Status)
	}

	license.Status = req.Status

	if err := k.SetLicense(ctx, license); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUpdateLicense,
		sdk.NewAttribute(types.AttributeKeyLicenseTypeID, req.LicenseTypeId),
		sdk.NewAttribute(types.AttributeKeyLicenseID, fmt.Sprintf("%d", req.Id)),
		sdk.NewAttribute(types.AttributeKeyStatus, req.Status),
	))

	return &types.MsgUpdateLicenseResponse{}, nil
}
