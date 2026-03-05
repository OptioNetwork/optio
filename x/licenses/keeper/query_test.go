package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/OptioNetwork/optio/testutil/keeper"
	"github.com/OptioNetwork/optio/testutil/sample"
	"github.com/OptioNetwork/optio/x/licenses/types"
)

func TestParamsQuery(t *testing.T) {
	k, ctx := keepertest.LicensesKeeper(t)
	params := k.GetParams(ctx)

	resp, err := k.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, params, resp.Params)
}

func TestLicenseTypeQuery(t *testing.T) {
	k, ctx := keepertest.LicensesKeeper(t)

	lt := types.LicenseType{
		Id: "node", Transferrable: true, MaxSupply: math.NewInt(100), IssuedCount: math.ZeroInt(),
	}
	require.NoError(t, k.SetLicenseType(ctx, lt))

	// Found
	resp, err := k.LicenseType(ctx, &types.QueryLicenseTypeRequest{Id: "node"})
	require.NoError(t, err)
	require.Equal(t, "node", resp.LicenseType.Id)
	require.True(t, resp.LicenseType.Transferrable)

	// Not found
	_, err = k.LicenseType(ctx, &types.QueryLicenseTypeRequest{Id: "missing"})
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestLicenseTypesQuery(t *testing.T) {
	k, ctx := keepertest.LicensesKeeper(t)

	require.NoError(t, k.SetLicenseType(ctx, types.LicenseType{Id: "a", MaxSupply: math.ZeroInt(), IssuedCount: math.ZeroInt()}))
	require.NoError(t, k.SetLicenseType(ctx, types.LicenseType{Id: "b", MaxSupply: math.ZeroInt(), IssuedCount: math.ZeroInt()}))

	resp, err := k.LicenseTypes(ctx, &types.QueryLicenseTypesRequest{})
	require.NoError(t, err)
	require.Len(t, resp.LicenseTypes, 2)
}

func TestLicenseQuery(t *testing.T) {
	k, ctx := keepertest.LicensesKeeper(t)
	holder := sample.AccAddress()

	license := types.License{
		Id: 1, Type: "node", Holder: holder, StartDate: "2026-01-01", Status: "active",
	}
	require.NoError(t, k.SetLicense(ctx, license))

	// Found
	resp, err := k.License(ctx, &types.QueryLicenseRequest{TypeId: "node", Id: 1})
	require.NoError(t, err)
	require.Equal(t, holder, resp.License.Holder)

	// Not found
	_, err = k.License(ctx, &types.QueryLicenseRequest{TypeId: "node", Id: 999})
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestLicensesByTypeQuery(t *testing.T) {
	k, ctx := keepertest.LicensesKeeper(t)

	require.NoError(t, k.SetLicense(ctx, types.License{Id: 1, Type: "node", Holder: sample.AccAddress(), StartDate: "2026-01-01", Status: "active"}))
	require.NoError(t, k.SetLicense(ctx, types.License{Id: 2, Type: "node", Holder: sample.AccAddress(), StartDate: "2026-01-01", Status: "active"}))
	require.NoError(t, k.SetLicense(ctx, types.License{Id: 1, Type: "other", Holder: sample.AccAddress(), StartDate: "2026-01-01", Status: "active"}))

	resp, err := k.LicensesByType(ctx, &types.QueryLicensesByTypeRequest{TypeId: "node"})
	require.NoError(t, err)
	require.Len(t, resp.Licenses, 2)
}

func TestLicensesByHolderQuery(t *testing.T) {
	k, ctx := keepertest.LicensesKeeper(t)
	holder := sample.AccAddress()
	other := sample.AccAddress()

	require.NoError(t, k.SetLicense(ctx, types.License{Id: 1, Type: "a", Holder: holder, StartDate: "2026-01-01", Status: "active"}))
	require.NoError(t, k.SetLicense(ctx, types.License{Id: 1, Type: "b", Holder: holder, StartDate: "2026-01-01", Status: "active"}))
	require.NoError(t, k.SetLicense(ctx, types.License{Id: 2, Type: "a", Holder: other, StartDate: "2026-01-01", Status: "active"}))

	resp, err := k.LicensesByHolder(ctx, &types.QueryLicensesByHolderRequest{Holder: holder})
	require.NoError(t, err)
	require.Len(t, resp.Licenses, 2)
}

func TestLicensesByHolderAndTypeQuery(t *testing.T) {
	k, ctx := keepertest.LicensesKeeper(t)
	holder := sample.AccAddress()

	require.NoError(t, k.SetLicense(ctx, types.License{Id: 1, Type: "a", Holder: holder, StartDate: "2026-01-01", Status: "active"}))
	require.NoError(t, k.SetLicense(ctx, types.License{Id: 2, Type: "a", Holder: holder, StartDate: "2026-01-01", Status: "active"}))
	require.NoError(t, k.SetLicense(ctx, types.License{Id: 1, Type: "b", Holder: holder, StartDate: "2026-01-01", Status: "active"}))

	resp, err := k.LicensesByHolderAndType(ctx, &types.QueryLicensesByHolderAndTypeRequest{Holder: holder, TypeId: "a"})
	require.NoError(t, err)
	require.Len(t, resp.Licenses, 2)
}

func TestAdminKeyQuery(t *testing.T) {
	k, ctx := keepertest.LicensesKeeper(t)
	addr := sample.AccAddress()

	ak := types.AdminKey{
		Address: addr,
		Grants:  []types.AdminKeyGrant{{Permission: "issue", LicenseTypes: []string{"t1"}}},
	}
	require.NoError(t, k.SetAdminKey(ctx, ak))

	// Found
	resp, err := k.AdminKey(ctx, &types.QueryAdminKeyRequest{Address: addr})
	require.NoError(t, err)
	require.Equal(t, addr, resp.AdminKey.Address)
	require.Len(t, resp.AdminKey.Grants, 1)

	// Not found
	_, err = k.AdminKey(ctx, &types.QueryAdminKeyRequest{Address: sample.AccAddress()})
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestAdminKeysQuery(t *testing.T) {
	k, ctx := keepertest.LicensesKeeper(t)

	require.NoError(t, k.SetAdminKey(ctx, types.AdminKey{
		Address: sample.AccAddress(),
		Grants:  []types.AdminKeyGrant{{Permission: "issue", LicenseTypes: []string{"t1"}}},
	}))
	require.NoError(t, k.SetAdminKey(ctx, types.AdminKey{
		Address: sample.AccAddress(),
		Grants:  []types.AdminKeyGrant{{Permission: "revoke", LicenseTypes: []string{"t1"}}},
	}))

	resp, err := k.AdminKeys(ctx, &types.QueryAdminKeysRequest{})
	require.NoError(t, err)
	require.Len(t, resp.AdminKeys, 2)
}

func TestAdminKeysByLicenseTypeQuery(t *testing.T) {
	k, ctx := keepertest.LicensesKeeper(t)

	addr1 := sample.AccAddress()
	addr2 := sample.AccAddress()
	addr3 := sample.AccAddress()

	// addr1: issue for t1
	require.NoError(t, k.SetAdminKey(ctx, types.AdminKey{
		Address: addr1,
		Grants:  []types.AdminKeyGrant{{Permission: "issue", LicenseTypes: []string{"t1"}}},
	}))
	// addr2: issue for t1, revoke for t2
	require.NoError(t, k.SetAdminKey(ctx, types.AdminKey{
		Address: addr2,
		Grants: []types.AdminKeyGrant{
			{Permission: "issue", LicenseTypes: []string{"t1"}},
			{Permission: "revoke", LicenseTypes: []string{"t2"}},
		},
	}))
	// addr3: revoke for t2 only
	require.NoError(t, k.SetAdminKey(ctx, types.AdminKey{
		Address: addr3,
		Grants:  []types.AdminKeyGrant{{Permission: "revoke", LicenseTypes: []string{"t2"}}},
	}))

	// Query by t1 (no permission filter) → addr1, addr2
	resp, err := k.AdminKeysByLicenseType(ctx, &types.QueryAdminKeysByLicenseTypeRequest{
		LicenseTypeId: "t1",
	})
	require.NoError(t, err)
	require.Len(t, resp.AdminKeys, 2)

	// Query by t2 (no permission filter) → addr2, addr3
	resp, err = k.AdminKeysByLicenseType(ctx, &types.QueryAdminKeysByLicenseTypeRequest{
		LicenseTypeId: "t2",
	})
	require.NoError(t, err)
	require.Len(t, resp.AdminKeys, 2)

	// Query by t1 with permission=issue → addr1, addr2
	resp, err = k.AdminKeysByLicenseType(ctx, &types.QueryAdminKeysByLicenseTypeRequest{
		LicenseTypeId: "t1",
		Permission:    "issue",
	})
	require.NoError(t, err)
	require.Len(t, resp.AdminKeys, 2)

	// Query by t2 with permission=issue → none
	resp, err = k.AdminKeysByLicenseType(ctx, &types.QueryAdminKeysByLicenseTypeRequest{
		LicenseTypeId: "t2",
		Permission:    "issue",
	})
	require.NoError(t, err)
	require.Len(t, resp.AdminKeys, 0)

	// Empty license_type_id → error
	_, err = k.AdminKeysByLicenseType(ctx, &types.QueryAdminKeysByLicenseTypeRequest{
		LicenseTypeId: "",
	})
	require.Error(t, err)
}
