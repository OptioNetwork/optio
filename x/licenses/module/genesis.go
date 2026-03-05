package licenses

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OptioNetwork/optio/x/licenses/keeper"
	"github.com/OptioNetwork/optio/x/licenses/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	for _, lt := range genState.LicenseTypes {
		if err := k.SetLicenseType(ctx, lt); err != nil {
			panic(err)
		}
	}

	// Track max license ID per type to rebuild counters
	maxIDs := make(map[string]uint64)
	for _, l := range genState.Licenses {
		if err := k.SetLicense(ctx, l); err != nil {
			panic(err)
		}
		if l.Id > maxIDs[l.Type] {
			maxIDs[l.Type] = l.Id
		}
	}

	// Rebuild per-type license counters from imported licenses
	for typeID, maxID := range maxIDs {
		if err := k.SetLicenseCounter(ctx, typeID, maxID); err != nil {
			panic(err)
		}
	}

	for _, ak := range genState.AdminKeys {
		if err := k.SetAdminKey(ctx, ak); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the module's exported genesis state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	lts, err := k.GetAllLicenseTypes(ctx)
	if err != nil {
		panic(err)
	}

	licenses, err := k.GetAllLicenses(ctx)
	if err != nil {
		panic(err)
	}

	aks, err := k.GetAllAdminKeys(ctx)
	if err != nil {
		panic(err)
	}

	if lts == nil {
		lts = []types.LicenseType{}
	}
	if licenses == nil {
		licenses = []types.License{}
	}
	if aks == nil {
		aks = []types.AdminKey{}
	}

	return &types.GenesisState{
		Params:       k.GetParams(ctx),
		LicenseTypes: lts,
		Licenses:     licenses,
		AdminKeys:    aks,
	}
}
