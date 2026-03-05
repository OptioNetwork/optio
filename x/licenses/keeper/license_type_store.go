package keeper

import (
	"context"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

func (k Keeper) SetLicenseType(ctx context.Context, lt types.LicenseType) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&lt)
	if err != nil {
		return err
	}
	return store.Set(types.GetLicenseTypeKey(lt.Id), bz)
}

func (k Keeper) GetLicenseType(ctx context.Context, id string) (types.LicenseType, bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.GetLicenseTypeKey(id))
	if err != nil {
		return types.LicenseType{}, false, err
	}
	if bz == nil {
		return types.LicenseType{}, false, nil
	}
	var lt types.LicenseType
	if err := k.cdc.Unmarshal(bz, &lt); err != nil {
		return types.LicenseType{}, false, err
	}
	return lt, true, nil
}

func (k Keeper) GetAllLicenseTypes(ctx context.Context) ([]types.LicenseType, error) {
	store := k.storeService.OpenKVStore(ctx)
	iterator, err := store.Iterator(types.LicenseTypeKey, prefixEndBytes(types.LicenseTypeKey))
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var lts []types.LicenseType
	for ; iterator.Valid(); iterator.Next() {
		var lt types.LicenseType
		if err := k.cdc.Unmarshal(iterator.Value(), &lt); err != nil {
			return nil, err
		}
		lts = append(lts, lt)
	}
	return lts, nil
}
