package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/OptioNetwork/optio/testutil/keeper"
	"github.com/OptioNetwork/optio/x/distro/keeper"
	"github.com/OptioNetwork/optio/x/distro/types"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, context.Context) {
	k, ctx := keepertest.DistroKeeper(t)
	return k, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)
}
