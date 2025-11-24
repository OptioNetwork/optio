package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/OptioNetwork/optio/testutil/keeper"
	"github.com/OptioNetwork/optio/x/lockup/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.LockupKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
