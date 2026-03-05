package keeper

import (
	"context"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

func (k Keeper) SetAdminKey(ctx context.Context, ak types.AdminKey) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&ak)
	if err != nil {
		return err
	}
	return store.Set(types.GetAdminKeyKey(ak.Address), bz)
}

func (k Keeper) GetAdminKey(ctx context.Context, address string) (types.AdminKey, bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.GetAdminKeyKey(address))
	if err != nil {
		return types.AdminKey{}, false, err
	}
	if bz == nil {
		return types.AdminKey{}, false, nil
	}
	var ak types.AdminKey
	if err := k.cdc.Unmarshal(bz, &ak); err != nil {
		return types.AdminKey{}, false, err
	}
	return ak, true, nil
}

func (k Keeper) DeleteAdminKey(ctx context.Context, address string) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Delete(types.GetAdminKeyKey(address))
}

func (k Keeper) GetAllAdminKeys(ctx context.Context) ([]types.AdminKey, error) {
	store := k.storeService.OpenKVStore(ctx)
	iterator, err := store.Iterator(types.AdminKeyStoreKey, prefixEndBytes(types.AdminKeyStoreKey))
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var aks []types.AdminKey
	for ; iterator.Valid(); iterator.Next() {
		var ak types.AdminKey
		if err := k.cdc.Unmarshal(iterator.Value(), &ak); err != nil {
			return nil, err
		}
		aks = append(aks, ak)
	}
	return aks, nil
}

// HasPermission checks whether address has the given permission for licenseTypeID.
func (k Keeper) HasPermission(ctx context.Context, address, permission, licenseTypeID string) bool {
	ak, found, err := k.GetAdminKey(ctx, address)
	if err != nil || !found {
		return false
	}
	for _, grant := range ak.Grants {
		if grant.Permission != permission {
			continue
		}
		for _, lt := range grant.LicenseTypes {
			if lt == licenseTypeID {
				return true
			}
		}
	}
	return false
}
