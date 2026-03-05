package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

func (k msgServer) TransferLicense(goCtx context.Context, req *types.MsgTransferLicense) (*types.MsgTransferLicenseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := sdk.AccAddressFromBech32(req.Recipient); err != nil {
		return nil, fmt.Errorf("invalid recipient address %q: %w", req.Recipient, err)
	}

	if req.Holder == req.Recipient {
		return nil, fmt.Errorf("cannot transfer license to the current holder")
	}

	license, found, err := k.GetLicense(ctx, req.LicenseTypeId, req.Id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrLicenseNotFound, "license (type=%s, id=%d) not found", req.LicenseTypeId, req.Id)
	}

	if license.Holder != req.Holder {
		return nil, errorsmod.Wrapf(types.ErrNotLicenseHolder, "signer %s is not the holder of license (type=%s, id=%d)", req.Holder, req.LicenseTypeId, req.Id)
	}

	lt, found, err := k.GetLicenseType(ctx, license.Type)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrLicenseTypeNotFound, "license type %s not found", license.Type)
	}
	if !lt.Transferrable {
		return nil, errorsmod.Wrapf(types.ErrLicenseNotTransferable, "license type %s is not transferrable", license.Type)
	}

	if err := k.DeleteLicense(ctx, license); err != nil {
		return nil, err
	}

	license.Holder = req.Recipient
	if err := k.SetLicense(ctx, license); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeTransferLicense,
		sdk.NewAttribute(types.AttributeKeyLicenseTypeID, req.LicenseTypeId),
		sdk.NewAttribute(types.AttributeKeyLicenseID, fmt.Sprintf("%d", req.Id)),
		sdk.NewAttribute(types.AttributeKeyHolder, req.Holder),
		sdk.NewAttribute(types.AttributeKeyRecipient, req.Recipient),
	))

	return &types.MsgTransferLicenseResponse{}, nil
}
