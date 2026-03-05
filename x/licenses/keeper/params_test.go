package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/OptioNetwork/optio/testutil/keeper"
	"github.com/OptioNetwork/optio/x/licenses/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.LicensesKeeper(t)
	params := k.GetParams(ctx)

	// Params were set by LicensesKeeper helper
	require.NotEmpty(t, params.Owner)

	// Round-trip: set then get
	newParams := types.Params{Owner: params.Owner}
	require.NoError(t, k.SetParams(ctx, newParams))
	require.EqualValues(t, newParams, k.GetParams(ctx))
}
