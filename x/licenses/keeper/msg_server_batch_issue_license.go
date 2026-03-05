package keeper

import (
	"context"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

func (k msgServer) BatchIssueLicense(goCtx context.Context, req *types.MsgBatchIssueLicense) (*types.MsgBatchIssueLicenseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(req.Entries) == 0 {
		return nil, fmt.Errorf("entries must not be empty")
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

	count := int64(len(req.Entries))
	if !lt.MaxSupply.IsZero() && lt.IssuedCount.AddRaw(count).GT(lt.MaxSupply) {
		return nil, errorsmod.Wrapf(types.ErrMaxSupplyReached, "license type %s: issuing %d would exceed max supply of %s (current: %s)", req.LicenseTypeId, count, lt.MaxSupply.String(), lt.IssuedCount.String())
	}

	// Validate all entries before issuing any
	for i, entry := range req.Entries {
		if _, err := sdk.AccAddressFromBech32(entry.Holder); err != nil {
			return nil, fmt.Errorf("entry %d: invalid holder address %q: %w", i, entry.Holder, err)
		}
		if entry.StartDate == "" {
			return nil, fmt.Errorf("entry %d: start_date is required", i)
		}
		startDate, err := time.Parse("2006-01-02", entry.StartDate)
		if err != nil {
			return nil, fmt.Errorf("entry %d: invalid start_date %q: must be YYYY-MM-DD format", i, entry.StartDate)
		}
		if entry.EndDate != "" {
			endDate, err := time.Parse("2006-01-02", entry.EndDate)
			if err != nil {
				return nil, fmt.Errorf("entry %d: invalid end_date %q: must be YYYY-MM-DD format", i, entry.EndDate)
			}
			if endDate.Before(startDate) {
				return nil, fmt.Errorf("entry %d: end_date %s must not be before start_date %s", i, entry.EndDate, entry.StartDate)
			}
		}
	}

	ids := make([]uint64, 0, count)
	for _, entry := range req.Entries {
		id, err := k.GetNextLicenseID(ctx, req.LicenseTypeId)
		if err != nil {
			return nil, err
		}

		license := types.License{
			Id:        id,
			Type:      req.LicenseTypeId,
			Holder:    entry.Holder,
			StartDate: entry.StartDate,
			EndDate:   entry.EndDate,
			Status:    "active",
		}

		if err := k.SetLicense(ctx, license); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	lt.IssuedCount = lt.IssuedCount.AddRaw(count)
	if err := k.SetLicenseType(ctx, lt); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBatchIssueLicense,
		sdk.NewAttribute(types.AttributeKeyLicenseTypeID, req.LicenseTypeId),
		sdk.NewAttribute("count", fmt.Sprintf("%d", count)),
	))

	return &types.MsgBatchIssueLicenseResponse{Ids: ids}, nil
}
