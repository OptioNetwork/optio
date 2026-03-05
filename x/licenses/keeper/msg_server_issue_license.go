package keeper

import (
	"context"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

func (k msgServer) IssueLicense(goCtx context.Context, req *types.MsgIssueLicense) (*types.MsgIssueLicenseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := sdk.AccAddressFromBech32(req.Holder); err != nil {
		return nil, fmt.Errorf("invalid holder address %q: %w", req.Holder, err)
	}

	if req.StartDate == "" {
		return nil, fmt.Errorf("start_date is required")
	}
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date %q: must be YYYY-MM-DD format", req.StartDate)
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date %q: must be YYYY-MM-DD format", req.EndDate)
		}
		if endDate.Before(startDate) {
			return nil, fmt.Errorf("end_date %s must not be before start_date %s", req.EndDate, req.StartDate)
		}
	}

	if !k.HasPermission(ctx, req.Issuer, "issue", req.LicenseTypeId) {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "%s does not have issue permission for license type %s", req.Issuer, req.LicenseTypeId)
	}

	lt, found, err := k.GetLicenseType(ctx, req.LicenseTypeId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrLicenseTypeNotFound, "license type %s not found", req.LicenseTypeId)
	}

	count := req.Count
	if count == 0 {
		count = 1
	}

	if !lt.MaxSupply.IsZero() && lt.IssuedCount.AddRaw(int64(count)).GT(lt.MaxSupply) {
		return nil, errorsmod.Wrapf(types.ErrMaxSupplyReached, "license type %s: issuing %d would exceed max supply of %s (current: %s)", req.LicenseTypeId, count, lt.MaxSupply.String(), lt.IssuedCount.String())
	}

	ids := make([]uint64, 0, count)
	for i := uint64(0); i < count; i++ {
		id, err := k.GetNextLicenseID(ctx, req.LicenseTypeId)
		if err != nil {
			return nil, err
		}

		license := types.License{
			Id:        id,
			Type:      req.LicenseTypeId,
			Holder:    req.Holder,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
			Status:    "active",
		}

		if err := k.SetLicense(ctx, license); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	lt.IssuedCount = lt.IssuedCount.AddRaw(int64(count))
	if err := k.SetLicenseType(ctx, lt); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeIssueLicense,
		sdk.NewAttribute(types.AttributeKeyLicenseTypeID, req.LicenseTypeId),
		sdk.NewAttribute(types.AttributeKeyHolder, req.Holder),
		sdk.NewAttribute("count", fmt.Sprintf("%d", count)),
	))

	return &types.MsgIssueLicenseResponse{Ids: ids}, nil
}
