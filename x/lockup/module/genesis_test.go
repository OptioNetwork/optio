package lockup_test

import (
	"testing"

	keepertest "github.com/OptioNetwork/optio/testutil/keeper"
	"github.com/OptioNetwork/optio/testutil/nullify"
	lockup "github.com/OptioNetwork/optio/x/lockup/module"
	"github.com/OptioNetwork/optio/x/lockup/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.LockupKeeper(t)
	lockup.InitGenesis(ctx, k, genesisState)
	got := lockup.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
