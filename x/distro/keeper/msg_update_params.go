package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OptioNetwork/optio/x/distro/types"
)

func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.GetAuthority() != req.Authority {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	updateParams := req.Params
	currentParams := k.GetParams(ctx)

	// DistributionStartDate, MonthsInHalvingPeriod, MaxSupply, Denom are not allowed to be changed
	updateParams.DistributionStartDate = currentParams.DistributionStartDate
	updateParams.MonthsInHalvingPeriod = currentParams.MonthsInHalvingPeriod
	updateParams.MaxSupply = currentParams.MaxSupply
	updateParams.Denom = currentParams.Denom

	if err := k.SetParams(ctx, updateParams); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
