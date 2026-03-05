package licenses_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	keepertest "github.com/OptioNetwork/optio/testutil/keeper"
	"github.com/OptioNetwork/optio/testutil/sample"
	licenses "github.com/OptioNetwork/optio/x/licenses/module"
	"github.com/OptioNetwork/optio/x/licenses/types"
)

func TestGenesis(t *testing.T) {
	holder1 := sample.AccAddress()
	holder2 := sample.AccAddress()
	admin1 := sample.AccAddress()
	admin2 := sample.AccAddress()
	owner := sample.AccAddress()

	genesisState := types.GenesisState{
		Params: types.Params{Owner: owner},
		LicenseTypes: []types.LicenseType{
			{Id: "node", Transferrable: true, MaxSupply: math.NewInt(100), IssuedCount: math.NewInt(2)},
			{Id: "validator", Transferrable: false, MaxSupply: math.ZeroInt(), IssuedCount: math.NewInt(1)},
		},
		Licenses: []types.License{
			{Id: 1, Type: "node", Holder: holder1, StartDate: "2026-01-01", EndDate: "2027-01-01", Status: "active"},
			{Id: 2, Type: "node", Holder: holder2, StartDate: "2026-02-01", Status: "active"},
			{Id: 1, Type: "validator", Holder: holder1, StartDate: "2026-01-01", Status: "revoked"},
		},
		AdminKeys: []types.AdminKey{
			{Address: admin1, Grants: []types.AdminKeyGrant{{Permission: "issue", LicenseTypes: []string{"node", "validator"}}}},
			{Address: admin2, Grants: []types.AdminKeyGrant{{Permission: "revoke", LicenseTypes: []string{"node"}}}},
		},
	}

	k, ctx := keepertest.LicensesKeeper(t)
	licenses.InitGenesis(ctx, k, genesisState)
	got := licenses.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	// Params
	require.Equal(t, owner, got.Params.Owner)

	// License types
	require.Len(t, got.LicenseTypes, 2)

	// Licenses
	require.Len(t, got.Licenses, 3)

	// Admin keys
	require.Len(t, got.AdminKeys, 2)

	// Verify license counter was rebuilt: issuing a new "node" license should get ID 3
	nextID, err := k.GetNextLicenseID(ctx, "node")
	require.NoError(t, err)
	require.Equal(t, uint64(3), nextID)

	// Verify "validator" counter: next ID should be 2
	nextID, err = k.GetNextLicenseID(ctx, "validator")
	require.NoError(t, err)
	require.Equal(t, uint64(2), nextID)
}
