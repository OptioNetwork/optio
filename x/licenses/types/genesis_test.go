package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/OptioNetwork/optio/testutil/sample"
	"github.com/OptioNetwork/optio/x/licenses/types"
)

func TestGenesisState_Validate(t *testing.T) {
	validAddr1 := sample.AccAddress()
	validAddr2 := sample.AccAddress()

	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid full state",
			genState: &types.GenesisState{
				Params: types.Params{Owner: validAddr1},
				LicenseTypes: []types.LicenseType{
					{Id: "node", MaxSupply: math.NewInt(100), IssuedCount: math.NewInt(1)},
				},
				Licenses: []types.License{
					{Id: 1, Type: "node", Holder: validAddr1, StartDate: "2026-01-01", Status: "active"},
				},
				AdminKeys: []types.AdminKey{
					{Address: validAddr2, Grants: []types.AdminKeyGrant{{Permission: "issue", LicenseTypes: []string{"node"}}}},
				},
			},
			valid: true,
		},
		{
			desc: "duplicate license type ID",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				LicenseTypes: []types.LicenseType{
					{Id: "node", MaxSupply: math.ZeroInt(), IssuedCount: math.ZeroInt()},
					{Id: "node", MaxSupply: math.ZeroInt(), IssuedCount: math.ZeroInt()},
				},
			},
			valid: false,
		},
		{
			desc: "duplicate license (type, id)",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				LicenseTypes: []types.LicenseType{
					{Id: "node", MaxSupply: math.ZeroInt(), IssuedCount: math.ZeroInt()},
				},
				Licenses: []types.License{
					{Id: 1, Type: "node", Holder: validAddr1, StartDate: "2026-01-01", Status: "active"},
					{Id: 1, Type: "node", Holder: validAddr2, StartDate: "2026-01-01", Status: "active"},
				},
			},
			valid: false,
		},
		{
			desc: "license references unknown type",
			genState: &types.GenesisState{
				Params:       types.DefaultParams(),
				LicenseTypes: []types.LicenseType{},
				Licenses: []types.License{
					{Id: 1, Type: "missing", Holder: validAddr1, StartDate: "2026-01-01", Status: "active"},
				},
			},
			valid: false,
		},
		{
			desc: "invalid license status",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				LicenseTypes: []types.LicenseType{
					{Id: "node", MaxSupply: math.ZeroInt(), IssuedCount: math.ZeroInt()},
				},
				Licenses: []types.License{
					{Id: 1, Type: "node", Holder: validAddr1, StartDate: "2026-01-01", Status: "suspended"},
				},
			},
			valid: false,
		},
		{
			desc: "invalid holder address",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				LicenseTypes: []types.LicenseType{
					{Id: "node", MaxSupply: math.ZeroInt(), IssuedCount: math.ZeroInt()},
				},
				Licenses: []types.License{
					{Id: 1, Type: "node", Holder: "bad", StartDate: "2026-01-01", Status: "active"},
				},
			},
			valid: false,
		},
		{
			desc: "duplicate admin key address",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				AdminKeys: []types.AdminKey{
					{Address: validAddr1, Grants: []types.AdminKeyGrant{{Permission: "issue", LicenseTypes: []string{"t1"}}}},
					{Address: validAddr1, Grants: []types.AdminKeyGrant{{Permission: "revoke", LicenseTypes: []string{"t1"}}}},
				},
			},
			valid: false,
		},
		{
			desc: "invalid admin key address",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				AdminKeys: []types.AdminKey{
					{Address: "bad", Grants: []types.AdminKeyGrant{{Permission: "issue", LicenseTypes: []string{"t1"}}}},
				},
			},
			valid: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
