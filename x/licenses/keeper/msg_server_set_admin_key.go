package keeper

import (
	"context"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

var validPermissions = map[string]struct{}{
	"issue":  {},
	"revoke": {},
	"update": {},
}

func (k msgServer) SetAdminKey(goCtx context.Context, req *types.MsgSetAdminKey) (*types.MsgSetAdminKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if params.Owner != req.Owner {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "signer %s is not the module owner %s", req.Owner, params.Owner)
	}

	if _, err := sdk.AccAddressFromBech32(req.Address); err != nil {
		return nil, fmt.Errorf("invalid address %q: %w", req.Address, err)
	}

	for _, grant := range req.Grants {
		if _, ok := validPermissions[grant.Permission]; !ok {
			return nil, fmt.Errorf("invalid permission %q: must be one of issue, revoke, update", grant.Permission)
		}
		if len(grant.LicenseTypes) == 0 {
			return nil, fmt.Errorf("grant for permission %q must include at least one license type", grant.Permission)
		}
	}

	ak := types.AdminKey{
		Address: req.Address,
		Grants:  req.Grants,
	}

	if err := k.Keeper.SetAdminKey(ctx, ak); err != nil {
		return nil, err
	}

	// Build comma-joined summaries for the event
	var perms []string
	var grantTypes []string
	for _, grant := range req.Grants {
		perms = append(perms, grant.Permission)
		grantTypes = append(grantTypes, strings.Join(grant.LicenseTypes, ","))
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSetAdminKey,
		sdk.NewAttribute(types.AttributeKeyAddress, req.Address),
		sdk.NewAttribute(types.AttributeKeyPermissions, strings.Join(perms, ",")),
		sdk.NewAttribute(types.AttributeKeyGrantTypes, strings.Join(grantTypes, ";")),
	))

	return &types.MsgSetAdminKeyResponse{}, nil
}
